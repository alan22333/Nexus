// nexus/internal/service/email_service.go
package service

import (
	"Nuxus/configs"
	"fmt"

	"gopkg.in/gomail.v2"
)

// SendVerificationMail 是一个专门发送验证码邮件的便捷函数
func SendRegisterMail(toEmail, code string) error {
	cfg := configs.Conf.SMTP

	subject := fmt.Sprintf("[%s] 您的邮箱验证码", cfg.FromName)

	// 这里是邮件正文，可以使用 HTML 来美化它
	// 确保内容清晰地告诉用户验证码是什么，以及它的有效期
	body := fmt.Sprintf(`
	<html>
	<body>
		<h3>您好！</h3>
		<p>感谢您注册 <strong>%s</strong>。您的邮箱验证码是：</p>
		<h2 style="font-weight: bold; color: #1E90FF;">%s</h2>
		<p>此验证码将在10分钟内失效，请尽快完成验证。</p>
		<p>如果这不是您本人的操作，请忽略此邮件。</p>
		<br/>
		<p>此致</p>
		<p><strong>%s 团队</strong></p>
	</body>
	</html>
	`, cfg.FromName, code, cfg.FromName)

	return sendMail(toEmail, subject, body)
}

// SendResetPasswordMail 发送重置密码邮件的便捷函数
func SendResetPasswordMail(toEmail, code string) error {
	cfg := configs.Conf.SMTP
	subject := fmt.Sprintf("[%s] 您的密码重置请求", cfg.FromName)
	body := fmt.Sprintf(`
    <html><body>
        <h3>您好！</h3>
        <p>我们收到了您在 <strong>%s</strong> 的密码重置请求。您的验证码是：</p>
        <h2 style="font-weight: bold; color: #FF4500;">%s</h2>
        <p>此验证码将在10分钟内失效。请使用此验证码来设置您的新密码。</p>
        <p>如果这不是您本人的操作，请忽略此邮件。</p>
    </body></html>
    `, cfg.FromName, code)
	return sendMail(toEmail, subject, body)
}

// sendMail 是底层的邮件发送实现
// 它不关心邮件内容，只负责发送
func sendMail(toEmail, subject, body string) error {
	// 从全局配置中获取 SMTP 信息
	cfg := configs.Conf.SMTP

	m := gomail.NewMessage()

	// 设置发件人
	// m.FormatAddress() 可以同时设置邮箱地址和发件人名称，避免乱码
	m.SetHeader("From", m.FormatAddress(cfg.Username, cfg.FromName))

	// 设置收件人
	m.SetHeader("To", toEmail)

	// 设置邮件主题
	m.SetHeader("Subject", subject)

	// 设置邮件正文，指定为 HTML 格式
	m.SetBody("text/html", body)

	// 创建一个拨号器，用于连接 SMTP 服务器
	// 参数：主机、端口、发件人邮箱、授权码
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// 发送邮件
	// d.DialAndSend() 会自动处理连接、认证和发送的全过程
	return d.DialAndSend(m)
}
