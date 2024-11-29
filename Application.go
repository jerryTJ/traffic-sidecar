package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jerryTJ/sidecar/cmd"
	"github.com/jerryTJ/sidecar/init/db"
	"github.com/jerryTJ/sidecar/init/logger"
	"github.com/jerryTJ/sidecar/internal/app"
	"github.com/jerryTJ/sidecar/tools"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func main() {
	cmd.Execute()
	db.Init(cmd.DB_USER, cmd.DB_PWD, cmd.DB_URL, cmd.DB_NAME)
	logger.Init(cmd.LoggerFile, zerolog.DebugLevel)

	//cron job
	c := cron.New()
	duration := fmt.Sprintf("@every %dm", cmd.Duration)
	coloringServer := app.ServerInfo{}

	serverInfos := make(map[string]app.ServerInfo)
	c.AddFunc(duration, func() {
		serverInfos = coloringServer.QueryServerInfos(serverInfos)
	})
	c.Start()

	handler := tools.ProxyHandler{TargetUrl: "http://gateway.devops.com", ServerInfos: serverInfos}
	http.HandleFunc("/", handler.ServeHTTP)
	// http server
	// Read ports from environment variable
	ports := cmd.Ports
	if ports == "" {
		fmt.Println("No ports specified in SERVER_PORTS environment variable")
		return
	}

	// Split the ports by comma and create a server for each port
	portList := strings.Split(ports, ",")
	for _, port := range portList {
		port = strings.TrimSpace(port) // Clean up any surrounding whitespace
		if port == "" {
			continue
		}
		// Run each server in a separate goroutine
		go func(p string) {
			fmt.Printf("Starting server on port %s...\n", p)
			if err := http.ListenAndServe(":"+p, nil); err != nil {
				fmt.Printf("Server on port %s stopped: %v\n", p, err)
			}
		}(port)
	}
	// Block the main goroutine to keep the servers running
	select {}
}
