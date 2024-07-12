package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
)

var pzPath = flag.String("p", "./", "project-zombie path")

func init() {
	flag.Parse()
}

func startServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
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
	log.Println("Start Server ...")
	cmd := exec.CommandContext(ctx, "cmd.exe", "/c", "start "+javaPath)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func startProxy(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.CommandContext(ctx, "frpc.exe", "-c", "./proxy/frpc.toml")
	log.Println("Start Proxy ...")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		startServer(ctx, &wg)
	}()
	go func() {
		startProxy(ctx, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	cancel()

	wg.Wait()

	log.Println("Bye!")
}
