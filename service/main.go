package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
)

var pzPath = flag.String("p", "./", "dedicated project-zombie path")
var pzParma = strings.Split(`//c ./java.exe -Djava.awt.headless=true -Dzomboid.steam=1 -Dzomboid.znetlog=1 -XX:+UseZGC -XX:-CreateCoredumpOnCrash -XX:-OmitStackTraceInFastThrow -Xms8g -Xmx8g -Djava.library.path=natives/;natives/win64/;. -cp %PZ_CLASSPATH% zombie.network.GameServer -statistic 0 %1 %2`, " ")

func init() {
	flag.Parse()

	if _, err := os.Stat(*pzPath); err != nil {
		log.Fatal(err)
	}
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
	log.Println("Start Server ...")
	cmd := exec.CommandContext(ctx, "cmd.exe", pzParma...)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func startProxy(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	frpcPath, err := filepath.Abs("./frpc.exe")
	if err != nil {
		log.Fatal(err)
	}
	frpcConfig, err := filepath.Abs("./frpc.toml")
	cmd := exec.CommandContext(ctx, "cmd.exe", "/c", frpcPath, "-c", frpcConfig)
	log.Println("Start Proxy ...")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	//go func() {
	//	startServer(ctx, &wg)
	//}()
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
