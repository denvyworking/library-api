package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$PwpqIpesHQr7S66hrZGkzumVjFuwZwU2l9Eq084Ctk0ZNZeHla3ni"

	// Проверяем разные пароли
	passwords := []string{"password", "admin", "123456", "Den"}

	for _, pass := range passwords {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
		if err == nil {
			fmt.Printf("Хэш совпадает с паролем: %s\n", pass)
			return
		}
	}

	fmt.Println("Не найден подходящий пароль")
}
