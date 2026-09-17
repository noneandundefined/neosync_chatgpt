// Package authtoken handles versioned, authenticated access tokens and opaque refresh tokens.
package authtoken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"neomatica/neosync/types"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	AccessTTL    = 24 * time.Hour
	RefreshTTL   = 30 * 24 * time.Hour
	accessPrefix = "v2."
)

func deriveKey(purpose string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	if err != nil || (len(key) != 16 && len(key) != 24 && len(key) != 32) {
		return nil, errors.New("invalid SUPER_SECRET_KEY")
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("neosync/auth/v2/" + purpose))

	return mac.Sum(nil), nil
}

func IPHash(ip string) (string, error) {
	key, err := deriveKey("client-ip")
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(ip))

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func Equal(a, b string) bool {
	return a != "" && b != "" && hmac.Equal([]byte(a), []byte(b))
}

func accessCipher() (cipher.AEAD, error) {
	key, err := deriveKey("access-encryption")
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func SealAccess(token types.RequestAuthToken) (string, error) {
	aead, err := accessCipher()
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	sealed := aead.Seal(nonce, nonce, data, []byte(accessPrefix))
	return accessPrefix + base64.RawURLEncoding.EncodeToString(sealed), nil
}

func OpenAccess(value string) (types.RequestAuthToken, error) {
	var token types.RequestAuthToken

	if !strings.HasPrefix(value, accessPrefix) {
		return token, errors.New("unsupported access token")
	}

	aead, err := accessCipher()
	if err != nil {
		return token, err
	}

	data, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(value, accessPrefix))
	if err != nil || len(data) < aead.NonceSize()+aead.Overhead() {
		return token, errors.New("invalid access token")
	}

	plain, err := aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], []byte(accessPrefix))
	if err != nil {
		return token, err
	}

	if err := json.Unmarshal(plain, &token); err != nil {
		return token, err
	}

	if token.UUID == "" || token.SessionID == "" || token.IPAddress == "" || token.Timestamp.IsZero() || token.Timestamp.After(time.Now().Add(time.Minute)) {
		return token, errors.New("invalid access claims")
	}

	return token, nil
}

func NewID() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}

func NewRefresh(sessionID string) (string, error) {
	secret, err := NewID()
	if err != nil {
		return "", err
	}

	return sessionID + "." + secret, nil
}

func RefreshSessionID(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 2 || len(parts[0]) != 64 || len(parts[1]) != 64 {
		return ""
	}

	for _, part := range parts {
		if _, err := hex.DecodeString(part); err != nil {
			return ""
		}
	}

	return parts[0]
}

func RefreshHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func SetCookie(w http.ResponseWriter, name, value string, expires time.Time) {
	w.Header().Set("Cache-Control", "no-store")

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   os.Getenv("GO_ENV") != "DEV",
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearCookie(w http.ResponseWriter, name string) {
	w.Header().Set("Cache-Control", "no-store")

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   os.Getenv("GO_ENV") != "DEV",
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearCookies(w http.ResponseWriter) {
	ClearCookie(w, "access_token")
	ClearCookie(w, "refresh_token")
}
