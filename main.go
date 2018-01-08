package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nicholaswang/cinfluxdb/cluster"
)

var (
	configFile = flag.String("config", "", "Configuration file to use")
)

func main() {
	flag.Parse()

	if *configFile == "" {
		fmt.Fprintln(os.Stderr, "Missing configuration file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg, err := cluster.LoadConfigFile(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Problem loading config file:", err)
	}

	p, err := cluster.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		p.Stop()
	}()

	log.Println("starting influxdb proxy...")
	p.Run()
}
