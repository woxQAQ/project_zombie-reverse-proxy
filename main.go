package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

var pzPath = flag.String("p", "./", "project-zombie path")
var zuluJDKPath = ""

func init() {
	flag.Parse()
}

func startServer(ctx context.Context) {
	err := os.Chdir(*pzPath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("PZ_CLASSPATH", "java/istack-commons-runtime.jar;java/jassimp.jar;java/javacord-2.0.17-shaded.jar;java/javax.activation-api.jar;java/jaxb-api.jar;java/jaxb-runtime.jar;java/lwjgl.jar;java/lwjgl-natives-windows.jar;java/lwjgl-glfw.jar;java/lwjgl-glfw-natives-windows.jar;java/lwjgl-jemalloc.jar;java/lwjgl-jemalloc-natives-windows.jar;java/lwjgl-opengl.jar;java/lwjgl-opengl-natives-windows.jar;java/lwjgl_util.jar;java/sqlite-jdbc-3.27.2.1.jar;java/trove-3.0.3.jar;java/uncommons-maths-1.2.3.jar;java/commons-compress-1.18.jar;java/")
	if err != nil {
		log.Fatal(err)
	}
	javaPath, err := filepath.Abs("./zulu21-win_x64/bin/java.exe")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("cmd.exe", "/c", "start "+javaPath)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func startFrpc(ctx context.Context) {
	cmd := exec.Command("frpc.exe", "-c", "./proxy/frpc.toml")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startServer(ctx)
	go startFrpc(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
}
