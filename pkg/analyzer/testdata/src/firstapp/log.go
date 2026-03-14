package myapp

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

func main() {
	// Стандартный log
	log.Print("hello world")        // должно найти
	log.Printf("hello %s", "world") // должно найти
	log.Println("hello world")      // должно найти

	log.Fatal("fatal error") // должно найти
	log.Panic("panic error") // должно найти

	// slog верхнего уровня
	slog.Info("info message")   // должно найти
	slog.Debug("debug message") // должно найти
	slog.Warn("warn message")   // должно найти
	slog.Error("error message") // должно найти

	// slog с дополнительными атрибутами
	slog.Info("message with attrs", "key", "value") // должно найти, сообщение "message with attrs"

	// Вызовы, которые НЕ должны быть найдены
	fmt.Println("not a log call") // не должно найти (другой пакет)
	strings.ToUpper("test")       // не должно найти

	// Вызов метода на логгере (пока не должно находить - отложим на потом)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("method call") // пока не будет найдено (сложный случай)
}
