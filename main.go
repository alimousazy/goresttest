package main

import (
	"flag"
	"fmt"
	"github.com/chat/chatserver"
	"os"
)

func main() {
	configPath := flag.String("config", "", "Please specify config path.")
	flag.Parse()
	if *configPath == "" {
		fmt.Fprintf(os.Stderr, "Config file should be specifed -config [path-to-config]\n")
		os.Exit(1)
	}
	config, err := chatserver.ReadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configration (%s)\n", err)
		os.Exit(1)
	}
	server := chatserver.NewChatServer(config)
	server.Start()
}
