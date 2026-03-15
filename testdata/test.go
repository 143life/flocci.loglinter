package main

import (
	"context"
	"log"
	"log/slog"
)

func main() {
	password := "secret123"

	// ---------- log ----------
	log.Println("valid message starts with lowercase")   // OK
	log.Println("Invalid message starts with uppercase") // должно ругаться (правило 1)

	log.Println("english only")  // OK
	log.Println("русский текст") // должно ругаться (правило 2)

	log.Println("clean message")        // OK
	log.Println("message with emoji 😊") // должно ругаться (правило 3: эмодзи)
	log.Println("message with !!!")     // должно ругаться (правило 3: спецсимволы)
	log.Println("message with ,.")      // должно ругаться (спецсимволы , и .)
	log.Println("message with 123")     // OK

	log.Printf("user password: %s", password)       // должно ругаться (правило 4)
	log.Println("api_key=abc123")                   // должно ругаться (правило 4, если "api_key" в паттернах)
	log.Println("operation completed successfully") // OK

	// ---------- slog ----------
	slog.Info("starting server")       // OK
	slog.Error("failed to connect")    // OK
	slog.Info("Error happened")        // должно ругаться (правило 1)
	slog.Info("message with !!!")      // должно ругаться (правило 3: спецсимволы)
	slog.Info("русский текст")         // должно ругаться (правило 2)
	slog.Info("user password: secret") // должно ругаться (правило 4)

	// slog with context
	ctx := context.Background()
	slog.InfoContext(ctx, "context message")   // OK
	slog.ErrorContext(ctx, "Error in context") // должно ругаться (правило 1)
}
