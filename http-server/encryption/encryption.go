package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"neomatica/neosync/infra/locale"
	"os"
)

var key []byte
var iv []byte

func Encrypt(tr locale.Translator, plaintext string) (string, error) {
	key, _ = base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	iv, _ = base64.StdEncoding.DecodeString(os.Getenv("IV"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-block-create-error"))
	}

	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-gcm-create-error"))
	}

	ciphertext := gsm.Seal(nil, iv, []byte(plaintext), nil)

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

func EncryptSeed(plaintext string) (string, error) {
	key, _ = base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	iv, _ = base64.StdEncoding.DecodeString(os.Getenv("IV"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := gsm.Seal(nil, iv, []byte(plaintext), nil)

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

func Decrypt(tr locale.Translator, encodedCiphertext string) (string, error) {
	key, _ = base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	iv, _ = base64.StdEncoding.DecodeString(os.Getenv("IV"))

	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-decode-error"))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-block-create-error"))
	}

	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-gcm-create-error"))
	}

	plaintext, err := gsm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", errors.New(tr.TErr("encryption-decrypt-error"))
	}

	return string(plaintext), nil
}
