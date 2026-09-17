package pkg

import (
	"crypto/tls"
	"fmt"
	"io"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

const smtpDeadline = 45 * time.Second

/* Отправка email письма на почту */
func SendEmail(to, title, content string, tr locale.Translator) error {
	from := os.Getenv("SMTP_EMAIL")
	host := os.Getenv("SMTP_ADDR")
	portStr := os.Getenv("SMTP_PORT")
	password := os.Getenv("SMTP_PASSWORD")
	goEnv := os.Getenv("GO_ENV")
	started := time.Now()

	port, _ := strconv.Atoi(portStr)
	ssl := goEnv != "DEV"

	logger.Mail("SendEmail START to={%s} from={%s} subject={%s} host={%s} port={%d} ssl={%t} env={%s} deadline=%s", to, from, title, host, port, ssl, goEnv, smtpDeadline)

	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", title)
	mail.SetBody("text/html", BuildEmailTemplate(content, tr))

	d := gomail.NewDialer(host, port, from, password)
	d.SSL = ssl

	err := dialAndSend(d, mail)
	elapsed := time.Since(started)
	if err != nil {
		logger.Mail("SendEmail FAIL to={%s} from={%s} subject={%s} host={%s} port={%d} ssl={%t} elapsed=%s err={%s}", to, from, title, host, port, ssl, elapsed, err.Error())
		logger.Error("SendEmail to={%s}: %s", to, err.Error())
		return fmt.Errorf("%s: %v", tr.TErr("email-send-failed"), err)
	}

	logger.Mail("SendEmail OK to={%s} from={%s} subject={%s} host={%s} port={%d} ssl={%t} elapsed=%s", to, from, title, host, port, ssl, elapsed)
	logger.Info("SendEmail email={%s}: Email has been sent!", to)
	return nil
}

/* gomail.DialAndSend hardcodes a 10s dial; nc to mail.nic.ru needs ~20s. */
func dialAndSend(d *gomail.Dialer, mail *gomail.Message) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", d.Host, d.Port), smtpDeadline)
	if err != nil {
		return err
	}

	tlsConfig := d.TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{ServerName: d.Host}
	}

	if d.SSL {
		conn = tls.Client(conn, tlsConfig)
	}

	c, err := smtp.NewClient(conn, d.Host)
	if err != nil {
		conn.Close()
		return err
	}
	defer c.Close()

	if !d.SSL {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}

	if d.Username != "" {
		auth := smtpAuth(c, d)
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return err
			}
		}
	}

	return gomail.Send(&smtpGoMailSender{c}, mail)
}

func smtpAuth(c *smtp.Client, d *gomail.Dialer) smtp.Auth {
	ok, auths := c.Extension("AUTH")
	if !ok {
		return nil
	}
	if strings.Contains(auths, "CRAM-MD5") {
		return smtp.CRAMMD5Auth(d.Username, d.Password)
	}
	return smtp.PlainAuth("", d.Username, d.Password, d.Host)
}

type smtpGoMailSender struct {
	*smtp.Client
}

func (s *smtpGoMailSender) Send(from string, to []string, msg io.WriterTo) error {
	if err := s.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err := s.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := s.Data()
	if err != nil {
		return err
	}
	if _, err = msg.WriteTo(w); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func BuildEmailTemplate(html string, tr locale.Translator) string {
	return fmt.Sprintf(`
		<body style="margin:0; padding:0; font-family:Arial,Helvetica,sans-serif; font-size:16px;">
			<div style="max-width:625px; margin:40px auto; background-color:#ffffff; overflow:hidden;">
				<div style="text-align:center;margin-bottom:40px;">
					<img src="https://neosync.neomatica.ru/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="Neomatica" width="211" />
				</div>

				<div style="color:#222222; line-height:200%%;">
					%s

					<p>%s</p>
				</div>

				<hr style="border:none;border-top:1px solid #e9edf2;margin:30px 0;">

				<div style="text-align:center; font-size:13px; color:#222222;">
					<img src="https://neosync.neomatica.ru/local/templates/neomatica/images/neomatica-logo.png" alt="Neomatica" width="55" />

					<p>%s</p>

					<table role="presentation" align="center" cellpadding="0" cellspacing="0">
						<tr>
							<td style="padding:0 4px;">
								<a href="https://www.youtube.com/channel/UCLFMfvdU1PjLMNHaMeA1vBw/"><img src="%s/local/templates/social-media/youtube-ico.png" width="30" alt="YouTube"></a>
							</td>
							<td style="padding:0 4px;">
								<a href="https://t.me/+QmORq_HZlwoYqiNj"><img src="%s/local/templates/social-media/tg-ico.png" width="30" alt="Telegram"></a>
							</td>
							<td style="padding:0 4px;">
								<a href="https://vk.com/neomatica/"><img src="%s/local/templates/social-media/vk-ico.png" width="30" alt="VK"></a>
							</td>
						</tr>
					</table>

					<p style="color:#adb1b8;">%s</p>
					<p style="color:#adb1b8;">%s</p>
				</div>
			</div>
		</body>
	`, html,
		tr.T("email-regards"),
		tr.T("email-stay-informed"),
		os.Getenv("CLIENT_URL"),
		os.Getenv("CLIENT_URL"),
		os.Getenv("CLIENT_URL"),
		tr.T("email-created-automatic"),
		fmt.Sprintf(tr.T("copyright"), time.Now().Year()),
	)
}
