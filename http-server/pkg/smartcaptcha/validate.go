package smartcaptcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"neomatica/neosync/config"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/pkg/clientip"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var validationClient = &http.Client{
	Timeout: 15 * time.Second,
	// A validation token is single-use; do not replay it to a redirect target.
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

var errVerificationFailed = errors.New("captcha verification failed")

func verificationFailed(tr locale.Translator) error {
	return errors.New(tr.TErr("smartcaptcha-verification-failed"))
}

func Verify(tr locale.Translator, token string, remoteIP string) error {
	return verify(context.Background(), tr, token, remoteIP)
}

func verify(ctx context.Context, tr locale.Translator, token string, remoteIP string) error {
	if os.Getenv("GO_ENV") == "DEV" {
		return nil
	}

	err := validate(ctx, validationClient, token, remoteIP)
	if err != nil && !errors.Is(err, errVerificationFailed) {
		logger.Error("smartcaptcha.Verify: %s", err.Error())
	}

	return validationResult(tr, err)
}

func validationResult(tr locale.Translator, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, errVerificationFailed) {
		return verificationFailed(tr)
	}

	return errors.New(tr.TErr("smartcaptcha-unavailable"))
}

// validate makes one attempt: a timed-out POST may have already consumed the token.
func validate(ctx context.Context, client *http.Client, token string, remoteIP string) error {
	if token == "" {
		return errVerificationFailed
	}

	secret := os.Getenv("YANDEX_SMARTCAPTCHA_SECRET_KEY")
	if secret == "" {
		return errors.New("YANDEX_SMARTCAPTCHA_SECRET_KEY is empty")
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("token", token)

	if ip := normalizeIP(remoteIP); ip != "" {
		form.Set("ip", ip)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.YANDEX_SMARTCAPTCHA_VALIDATE_URL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("validate request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("validate returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read validation response: %w", err)
	}

	var result ValidateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse validation response: %w", err)
	}

	switch result.Status {
	case "ok":
		return nil
	case "failed":
		return errVerificationFailed
	default:
		return errors.New("unexpected validation response status")
	}
}

// VerifyRequest sends the verified client address only when the request came
// through a configured trusted proxy.
func VerifyRequest(tr locale.Translator, token string, r *http.Request) error {
	return verify(r.Context(), tr, token, clientip.IP(r))
}

func normalizeIP(remoteIP string) string {
	remoteIP = strings.TrimSpace(remoteIP)
	if remoteIP == "" {
		return ""
	}

	if host, _, err := net.SplitHostPort(remoteIP); err == nil {
		return host
	}

	return remoteIP
}
