package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jerryTJ/sidecar/cmd"
	"github.com/jerryTJ/sidecar/init/logger"
	"github.com/jerryTJ/sidecar/tools"
	"github.com/rs/zerolog"
)

func main() {
	cmd.Execute()
	logger.Init(cmd.LoggerFile, zerolog.DebugLevel)

	handler := tools.ProxyHandler{TargetUrl: "http://gateway.devops.com", Addr: "traffic.devops.com:10080"}
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
