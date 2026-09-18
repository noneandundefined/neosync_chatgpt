package util

import (
	"math/rand"
	"time"
)

/* Генерация рандомного пароля */
func GeneratePassword() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	password := make([]byte, 10)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range password {
		password[i] = characters[rand.Intn(len(characters))]
	}

	return string(password)
}
