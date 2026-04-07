package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	cfg       configs.MailConfig
	dialer    *gomail.Dialer
	templates *template.Template
}

func NewSMTPMailer(cfg configs.MailConfig) (*SMTPMailer, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	pattern := filepath.Join(dir, "templates", "*.html")

	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &SMTPMailer{
		cfg:       cfg,
		dialer:    dialer,
		templates: tmpl,
	}, nil
}

func (m *SMTPMailer) send(to, subject, tmplName string, data any) error {
	var body bytes.Buffer
	if err := m.templates.ExecuteTemplate(&body, tmplName, data); err != nil {
		return fmt.Errorf("failed to render template %s: %w", tmplName, err)
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.cfg.From, m.cfg.FromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body.String())

	return m.dialer.DialAndSend(msg)
}

func (m *SMTPMailer) SendVerificationEmail(_ context.Context, to, name, otp string) error {
	return m.send(to, "Verify your email", "verification.html", map[string]string{
		"Name": name,
		"OTP":  otp,
	})
}

func (m *SMTPMailer) SendPasswordResetEmail(_ context.Context, to, name, resetLink string) error {
	return m.send(to, "Reset your password", "password_reset.html", map[string]string{
		"Name":      name,
		"ResetLink": resetLink,
	})
}

func (m *SMTPMailer) SendWelcomeEmail(_ context.Context, to, name string) error {
	return m.send(to, "Welcome!", "welcome.html", map[string]string{
		"Name": name,
	})
}

func (m *SMTPMailer) SendNotificationEmail(_ context.Context, to, name, subject, message string) error {
	return m.send(to, subject, "notification.html", map[string]string{
		"Name":    name,
		"Message": message,
	})
}
