package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		cmd := exec.CommandContext(ctx, "frpc.exe", "-c", "./frpc.toml")
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Bye!")
}
