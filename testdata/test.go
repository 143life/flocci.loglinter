package main

import (
	"log"
)

func main() {
	// ---------- Правило 1: строчная буква в начале ----------
	log.Println("valid message starts with lowercase")   // OK
	log.Println("Invalid message starts with uppercase") // Должно ругаться: "log message should start with lowercase"

	// ---------- Правило 2: только английский язык (не кириллица и т.п.) ----------
	log.Println("english only")  // OK
	log.Println("русский текст") // Должно ругаться: "log message contains non-ASCII characters"

	// ---------- Правило 3: без спецсимволов и эмодзи ----------
	log.Println("clean message")        // OK
	log.Println("message with emoji 😊") // Должно ругаться: "log message contains emoji"
	log.Println("message with !!!")     // Должно ругаться (спецсимволы), но пока этой проверки нет
	log.Println("message with ...")     // Должно ругаться (спецсимволы), тоже нет проверки

	// ---------- Правило 4: без чувствительных данных ----------
	password := "secret123"
	log.Printf("user password: %s", password)       // Должно ругаться: "potential sensitive data: "password""
	log.Println("api_key=abc123")                   // Должно ругаться, если "api_key" в списке sensitive_patterns
	log.Println("token generated")                  // OK, если "token" нет в сообщении как отдельное слово (но здесь есть слово "token", может сработать, если паттерн "token" просто как подстрока)
	log.Println("operation completed successfully") // OK
}
