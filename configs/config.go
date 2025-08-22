package configs

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// var Conf = new(Config)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	SMTP   SMTPConfig   `mapstructure:"smtp"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Qiniu  QiniuConfig  `mapstructure:"qiniu"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// SMTPConfig 定义了发送邮件所需的 SMTP 服务器配置
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	FromName string `mapstructure:"fromName"`
}

// JWTConfig 定义了生成和验证 JWT 所需的配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expireHours"`
}

type QiniuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
	Zone      string `mapstructure:"zone"`
}


// LoadConfig 用于Wire依赖注入
func LoadConfig() (*Config, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	// 可选：监视配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&config); err != nil {
			log.Printf("viper.Unmarshal on config change failed, err: %v", err)
		}
	})

	log.Println("Configuration loaded successfully!")
	return &config, nil
}
