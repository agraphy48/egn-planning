package mailer

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// SMTPConfig holds SMTP server connection details
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	Secure   bool // true for direct SSL/TLS (port 465)
}

var currentConfig *SMTPConfig

// LoadEnv loads .env files from current or parent directories if present
func LoadEnv() {
	searchPaths := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("backend", ".env"),
		filepath.Join(".", "backend", ".env"),
	}

	for _, p := range searchPaths {
		if file, err := os.Open(p); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					// Strip wrapping quotes
					val = strings.Trim(val, `"'`)
					// Always update env from .env file
					_ = os.Setenv(key, val)
				}
			}
			log.Printf("[Mailer] Loaded environment variables from %s\n", p)
			break
		}
	}
}

// GetConfig reads the active SMTP configuration from environment
func GetConfig() *SMTPConfig {
	LoadEnv()

	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		host = "smtp.gmail.com"
	}

	portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	port := 587
	if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
		port = p
	}

	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	pass := strings.TrimSpace(os.Getenv("SMTP_PASS"))
	// Clean whitespace in app passwords (e.g., "abcd efgh ijkl mnop" -> "abcdefghijklmnop")
	if strings.Contains(host, "gmail.com") {
		pass = strings.ReplaceAll(pass, " ", "")
	}

	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if from == "" {
		if user != "" {
			from = fmt.Sprintf("Pulse Workspace <%s>", user)
		} else {
			from = "Pulse Workspace <noreply@pulse.workspace>"
		}
	}

	secure := port == 465 || strings.ToLower(os.Getenv("SMTP_SECURE")) == "true"

	currentConfig = &SMTPConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		From:     from,
		Secure:   secure,
	}

	return currentConfig
}

// IsConfigured checks whether SMTP has credentials configured
func IsConfigured() bool {
	cfg := GetConfig()
	return cfg.User != "" && cfg.Password != "" && cfg.Host != ""
}

// GetStatus returns the current state of SMTP service for admin monitoring
func GetStatus() map[string]interface{} {
	cfg := GetConfig()
	configured := IsConfigured()

	maskedPass := ""
	if len(cfg.Password) > 4 {
		maskedPass = cfg.Password[:2] + strings.Repeat("*", len(cfg.Password)-4) + cfg.Password[len(cfg.Password)-2:]
	} else if len(cfg.Password) > 0 {
		maskedPass = "****"
	}

	return map[string]interface{}{
		"configured":  configured,
		"smtp_host":   cfg.Host,
		"smtp_port":   cfg.Port,
		"smtp_user":   cfg.User,
		"smtp_from":   cfg.From,
		"smtp_secure": cfg.Secure,
		"has_pass":    cfg.Password != "",
		"masked_pass": maskedPass,
	}
}

// sendRawMail sends an email message via STARTTLS (port 587) or Direct SSL/TLS (port 465)
func sendRawMail(cfg *SMTPConfig, to string, subject string, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	host := cfg.Host

	// Build MIME Header
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", cfg.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	auth := smtp.PlainAuth("", cfg.User, cfg.Password, host)

	// Mode 1: Direct SSL/TLS (typically port 465)
	if cfg.Secure || cfg.Port == 465 {
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("gagal menghubungkan TLS ke %s: %w", addr, err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("gagal membuat SMTP client: %w", err)
		}
		defer client.Quit()

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("otentikasi SMTP ditolak: %w", err)
		}

		if err := client.Mail(cfg.User); err != nil {
			return fmt.Errorf("kesalahan MAIL FROM: %w", err)
		}

		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("kesalahan RCPT TO: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("kesalahan DATA: %w", err)
		}
		defer w.Close()

		if _, err := w.Write(msg.Bytes()); err != nil {
			return fmt.Errorf("gagal menulis isi email: %w", err)
		}

		return nil
	}

	// Mode 2: Standard STARTTLS (typically port 587 or 25)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke %s: %w", addr, err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("gagal membuat SMTP client: %w", err)
	}
	defer client.Quit()

	// Check if server supports STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("gagal memulai STARTTLS: %w", err)
		}
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("otentikasi SMTP ditolak: %w", err)
	}

	if err := client.Mail(cfg.User); err != nil {
		return fmt.Errorf("kesalahan MAIL FROM: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("kesalahan RCPT TO: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("kesalahan DATA: %w", err)
	}
	defer w.Close()

	if _, err := w.Write(msg.Bytes()); err != nil {
		return fmt.Errorf("gagal menulis isi email: %w", err)
	}

	return nil
}

// InvitationEmailData defines payload for email template
type InvitationEmailData struct {
	Name        string
	Email       string
	Role        string
	InviterName string
	Projects    []string
	InviteURL   string
	ExpiresIn   string
}

const invitationHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Undangan Bergabung ke Pulse Workspace</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f8fafc; margin: 0; padding: 24px; color: #1e293b; }
    .card { max-width: 580px; margin: 0 auto; background: #ffffff; border-radius: 16px; border: 1px solid #e2e8f0; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); }
    .header { background: #0f172a; padding: 32px 28px; text-align: center; }
    .header h1 { color: #ffffff; margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -0.02em; }
    .header p { color: #94a3b8; margin: 6px 0 0; font-size: 13px; }
    .content { padding: 32px 28px; line-height: 1.6; }
    .badge { display: inline-block; background: #eff6ff; color: #2563eb; font-weight: 600; padding: 4px 12px; border-radius: 9999px; font-size: 12px; text-transform: uppercase; margin-bottom: 12px; }
    .projects-box { background: #f8fafc; border-radius: 12px; padding: 16px; margin: 20px 0; border: 1px solid #e2e8f0; }
    .projects-title { font-size: 12px; font-weight: 700; color: #475569; text-transform: uppercase; margin-bottom: 8px; }
    .project-item { display: inline-block; background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 5px 12px; font-size: 13px; color: #334155; margin: 4px; font-weight: 500; }
    .btn-container { text-align: center; margin: 32px 0 24px; }
    .btn { display: inline-block; background: #2563eb; color: #ffffff !important; text-decoration: none; padding: 14px 32px; border-radius: 10px; font-weight: 600; font-size: 14px; box-shadow: 0 2px 4px rgba(37, 99, 235, 0.2); }
    .footer { border-top: 1px solid #e2e8f0; padding: 20px 28px; font-size: 12px; color: #64748b; background: #f8fafc; text-align: center; }
    .link-alt { word-break: break-all; font-family: monospace; font-size: 11px; color: #2563eb; }
  </style>
</head>
<body>
  <div class="card">
    <div class="header">
      <h1>PULSE WORKSPACE</h1>
      <p>Kedukan Aksara Enterprise Project Management</p>
    </div>
    <div class="content">
      <span class="badge">Undangan Akses Proyek</span>
      <h2 style="margin: 0 0 16px; font-size: 18px; color: #0f172a;">Halo {{if .Name}}{{.Name}}{{else}}Bapak/Ibu{{end}},</h2>
      <p style="margin: 0 0 16px;">
        Anda telah diundang oleh <strong>{{.InviterName}}</strong> untuk berkolaborasi di <strong>Pulse Workspace</strong> dengan hak akses sebagai <strong>{{.Role}}</strong>.
      </p>
      {{if .Projects}}
      <div class="projects-box">
        <div class="projects-title">Proyek yang Diberikan Hak Akses:</div>
        {{range .Projects}}
        <span class="project-item">📁 {{.}}</span>
        {{end}}
      </div>
      {{end}}
      <p style="margin: 0 0 8px; color: #475569; font-size: 14px;">
        Silakan klik tombol di bawah untuk mengonfirmasi dan mengaktifkan akun Anda. Sistem akan secara otomatis mengaktifkan akses Anda dan mengirimkan kredensial login (username & password) ke email Anda:
      </p>
      <div class="btn-container">
        <a href="{{.InviteURL}}" class="btn" target="_blank">Konfirmasi & Aktifkan Akun</a>
      </div>
      <p style="font-size: 12px; color: #94a3b8; margin: 24px 0 0; text-align: center;">
        * Tautan aktivasi ini aman dan berlaku selama 7 hari.
      </p>
      <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0 16px;" />
      <p style="font-size: 11px; color: #64748b; margin: 0;">
        Jika tombol di atas tidak dapat diklik, salin dan buka tautan berikut di browser Anda:<br>
        <a href="{{.InviteURL}}" class="link-alt">{{.InviteURL}}</a>
      </p>
    </div>
    <div class="footer">
      &copy; 2026 Kedukan Aksara. Seluruh hak cipta dilindungi undang-undang.<br>
      Pesan otomatis ini dikirimkan ke <strong>{{.Email}}</strong>.
    </div>
  </div>
</body>
</html>`

// SendInvitationEmail dispatches an invitation email using the configured SMTP server
func SendInvitationEmail(toEmail, toName, inviterName, roleName string, projectNames []string, inviteURL string) error {
	cfg := GetConfig()
	if !IsConfigured() {
		log.Printf("[Mailer] SMTP belum terkonfigurasi (SMTP_PASS / SMTP_USER kosong). Email tidak terkirim otomatis ke %s. Invite URL: %s\n", toEmail, inviteURL)
		return fmt.Errorf("SMTP belum dikonfigurasi pada server (silakan atur SMTP_USER dan SMTP_PASS)")
	}

	roleLabel := roleName
	switch roleName {
	case "client":
		roleLabel = "Klien (Client Portal)"
	case "member":
		roleLabel = "Anggota Tim (Team Member)"
	case "pm":
		roleLabel = "Manajer Proyek (Project Manager)"
	case "super_admin":
		roleLabel = "Administrator Sistem (Super Admin)"
	}

	tmpl, err := template.New("invitation").Parse(invitationHTMLTemplate)
	if err != nil {
		return fmt.Errorf("gagal mem-parsing template email: %w", err)
	}

	data := InvitationEmailData{
		Name:        toName,
		Email:       toEmail,
		Role:        roleLabel,
		InviterName: inviterName,
		Projects:    projectNames,
		InviteURL:   inviteURL,
		ExpiresIn:   "7 hari",
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("gagal merender template email: %w", err)
	}

	subject := fmt.Sprintf("Undangan Akses Pulse Workspace: Kolaborasi Proyek [%s]", roleLabel)

	log.Printf("[Mailer] Mengirimkan email undangan ke %s via %s:%d...\n", toEmail, cfg.Host, cfg.Port)
	err = sendRawMail(cfg, toEmail, subject, body.String())
	if err != nil {
		log.Printf("[Mailer] Gagal mengirimkan email ke %s: %v\n", toEmail, err)
		return err
	}

	log.Printf("[Mailer] Email undangan sukses terkirim ke %s!\n", toEmail)
	return nil
}

// WelcomeCredentialsEmailData defines payload for credentials email template
type WelcomeCredentialsEmailData struct {
	Name     string
	Email    string
	Password string
	Role     string
	Projects []string
	LoginURL string
}

const welcomeCredentialsHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Konfirmasi Akun & Kredensial Login Pulse Workspace</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f8fafc; margin: 0; padding: 24px; color: #1e293b; }
    .card { max-width: 580px; margin: 0 auto; background: #ffffff; border-radius: 16px; border: 1px solid #e2e8f0; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); }
    .header { background: #0f172a; padding: 32px 28px; text-align: center; }
    .header h1 { color: #ffffff; margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -0.02em; }
    .header p { color: #94a3b8; margin: 6px 0 0; font-size: 13px; }
    .content { padding: 32px 28px; line-height: 1.6; }
    .badge { display: inline-block; background: #ecfdf5; color: #059669; font-weight: 600; padding: 4px 12px; border-radius: 9999px; font-size: 12px; text-transform: uppercase; margin-bottom: 12px; border: 1px solid #a7f3d0; }
    .cred-box { background: #f8fafc; border-radius: 12px; padding: 20px; margin: 20px 0; border: 1px solid #cbd5e1; }
    .cred-title { font-size: 11px; font-weight: 700; color: #475569; text-transform: uppercase; margin-bottom: 12px; letter-spacing: 0.05em; border-bottom: 1px solid #e2e8f0; padding-bottom: 6px; }
    .pass-highlight { background: #fef08a; padding: 3px 8px; border-radius: 6px; color: #854d0e; font-weight: 800; border: 1px dashed #ca8a04; font-family: monospace; font-size: 14px; }
    .projects-box { background: #f1f5f9; border-radius: 12px; padding: 14px 16px; margin: 16px 0; border: 1px solid #e2e8f0; }
    .projects-title { font-size: 11px; font-weight: 700; color: #475569; text-transform: uppercase; margin-bottom: 8px; }
    .project-item { display: inline-block; background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #334155; margin: 3px; font-weight: 500; }
    .btn-container { text-align: center; margin: 30px 0 20px; }
    .btn { display: inline-block; background: #2563eb; color: #ffffff !important; text-decoration: none; padding: 14px 34px; border-radius: 10px; font-weight: 600; font-size: 14px; box-shadow: 0 2px 4px rgba(37, 99, 235, 0.2); }
    .security-note { background: #eff6ff; border-left: 4px solid #3b82f6; padding: 12px 14px; border-radius: 0 8px 8px 0; font-size: 12px; color: #1e40af; margin-top: 20px; }
    .footer { border-top: 1px solid #e2e8f0; padding: 20px 28px; font-size: 12px; color: #64748b; background: #f8fafc; text-align: center; }
  </style>
</head>
<body>
  <div class="card">
    <div class="header">
      <h1>PULSE WORKSPACE</h1>
      <p>Kedukan Aksara Enterprise Project Management</p>
    </div>
    <div class="content">
      <span class="badge">✓ Akun Berhasil Diaktifkan</span>
      <h2 style="margin: 0 0 16px; font-size: 18px; color: #0f172a;">Halo {{if .Name}}{{.Name}}{{else}}Bapak/Ibu{{end}},</h2>
      <p style="margin: 0 0 16px; font-size: 14px; color: #334155;">
        Selamat! Pendaftaran dan konfirmasi akun Anda di <strong>Pulse Workspace</strong> telah berhasil diselesaikan. Berikut adalah rincian kredensial akun Anda untuk login ke sistem:
      </p>

      <div class="cred-box">
        <div class="cred-title">Kredensial Akses Pengguna</div>
        <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
          <tr>
            <td style="padding: 7px 0; color: #64748b; width: 38%;">Email / Username:</td>
            <td style="padding: 7px 0; font-weight: 700; color: #0f172a; font-family: monospace;">{{.Email}}</td>
          </tr>
          <tr>
            <td style="padding: 7px 0; color: #64748b;">Kata Sandi (Password):</td>
            <td style="padding: 7px 0;"><span class="pass-highlight">{{.Password}}</span></td>
          </tr>
          <tr>
            <td style="padding: 7px 0; color: #64748b;">Hak Akses / Peran:</td>
            <td style="padding: 7px 0; font-weight: 600; color: #2563eb;">{{.Role}}</td>
          </tr>
        </table>
      </div>

      {{if .Projects}}
      <div class="projects-box">
        <div class="projects-title">Proyek yang Diberikan Hak Akses:</div>
        {{range .Projects}}
        <span class="project-item">📁 {{.}}</span>
        {{end}}
      </div>
      {{end}}

      <div class="btn-container">
        <a href="{{.LoginURL}}" class="btn" target="_blank">Masuk ke Pulse Workspace</a>
      </div>

      <div class="security-note">
        <strong>Tips Keamanan:</strong> Simpan email ini sebagai arsip kredensial Anda. Anda dapat mengubah kata sandi kapan saja setelah masuk melalui menu Profil Pengguna di pojok kanan atas.
      </div>
    </div>
    <div class="footer">
      &copy; 2026 Kedukan Aksara. Seluruh hak cipta dilindungi undang-undang.<br>
      Kredensial dikirimkan secara otomatis ke <strong>{{.Email}}</strong>.
    </div>
  </div>
</body>
</html>`

// SendWelcomeCredentialsEmail dispatches username and password to the user upon confirmation
func SendWelcomeCredentialsEmail(toEmail, toName, password, roleName string, projectNames []string, loginURL string) error {
	cfg := GetConfig()
	if !IsConfigured() {
		log.Printf("[Mailer] SMTP belum terkonfigurasi. Kredensial tidak dapat dikirimkan ke %s\n", toEmail)
		return fmt.Errorf("SMTP belum dikonfigurasi pada server")
	}

	roleLabel := roleName
	switch roleName {
	case "client":
		roleLabel = "Klien (Client Portal)"
	case "member":
		roleLabel = "Anggota Tim (Team Member)"
	case "pm":
		roleLabel = "Manajer Proyek (Project Manager)"
	case "super_admin":
		roleLabel = "Administrator Sistem (Super Admin)"
	}

	tmpl, err := template.New("credentials").Parse(welcomeCredentialsHTMLTemplate)
	if err != nil {
		return fmt.Errorf("gagal mem-parsing template email kredensial: %w", err)
	}

	data := WelcomeCredentialsEmailData{
		Name:     toName,
		Email:    toEmail,
		Password: password,
		Role:     roleLabel,
		Projects: projectNames,
		LoginURL: loginURL,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("gagal merender template kredensial: %w", err)
	}

	subject := fmt.Sprintf("Akun Pulse Workspace Aktif - Kredensial Login Anda [%s]", toEmail)

	log.Printf("[Mailer] Mengirimkan email kredensial login ke %s...\n", toEmail)
	err = sendRawMail(cfg, toEmail, subject, body.String())
	if err != nil {
		log.Printf("[Mailer] Gagal mengirimkan kredensial ke %s: %v\n", toEmail, err)
		return err
	}

	log.Printf("[Mailer] Email kredensial login sukses terkirim ke %s!\n", toEmail)
	return nil
}

// SendTestEmail sends a quick test email to verify credentials
func SendTestEmail(toEmail string) error {
	cfg := GetConfig()
	if !IsConfigured() {
		return fmt.Errorf("SMTP belum dikonfigurasi (SMTP_USER atau SMTP_PASS kosong)")
	}

	subject := "Test Koneksi SMTP: Pulse Workspace"
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; padding: 24px; color: #1e293b;">
  <h2>Koneksi SMTP Berhasil!</h2>
  <p>Pesan ini menandakan bahwa konfigurasi SMTP pada <strong>Pulse Workspace</strong> berjalan dengan normal.</p>
  <ul>
    <li>Host: %s</li>
    <li>Port: %d</li>
    <li>Pengirim: %s</li>
    <li>Waktu: %s</li>
  </ul>
</body>
</html>`, cfg.Host, cfg.Port, cfg.From, time.Now().Format(time.RFC1123))

	return sendRawMail(cfg, toEmail, subject, html)
}
