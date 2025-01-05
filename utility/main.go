package main

import (
	"log/slog"
	"os"

	"github.com/kofuk/dpb/fixed-response/internal/fixedresponse"
)

var commands = map[string]func(){
	"fixed-response": fixedresponse.Run,
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	mode := os.Getenv("MODE")
	fn, ok := commands[mode]
	if !ok {
		slog.Error("Invalid mode")
		os.Exit(1)
	}

	fn()
}
