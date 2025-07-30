package configs

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	SMTP   SMTPConfig   `mapstructure:"smtp"`
	JWT    JWTConfig    `mapstructure:"jwt"`
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

func Init() {
	// 设置工作目录，方便读取
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 设置配置文件，不带后缀
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 添加配置文件的搜索路径
	// 这里假设 `configs` 目录与可执行文件或项目根目录在同一层级
	viper.AddConfigPath(workDir)

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper.ReadInConfig failed, err: %v", err)
	}

	// 解析配置
	err = viper.Unmarshal(Conf)
	if err != nil {
		log.Fatalf("viper.Unmarshal failed, err: %v", err)
	}

	// (可选) 开启监视功能，当配置文件发生变化时自动重新加载
	// 这在开发环境中非常有用，但在生产环境中需要谨慎使用
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			log.Printf("viper.Unmarshal on config change failed, err: %v", err)
		}
	})

	log.Println("Configuration loaded and initialized successfully!")

}
