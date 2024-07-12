package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		frpcPath, err := filepath.Abs("./frpc.exe")
		if err != nil {
			log.Fatal(err)
		}
		frpcConfig, err := filepath.Abs("./frpc.toml")
		cmd := exec.CommandContext(ctx, "cmd.exe", "/c", frpcPath, "-c", frpcConfig)
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Bye!")
}
