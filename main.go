package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	var cfgfile *string = flag.String("config", "", "configuration file")

	hosts := make(map[string]string)
	sites := make(map[string]*Site)

	flag.Parse()

	if *cfgfile == "" {
		usage()
	}

	cfg, err := ReadConfig(*cfgfile)
	if err != nil {
		log.Printf("opening %s failed: %v", *cfgfile, err)
		os.Exit(1)
	}

	var access_f io.WriteCloser
	if cfg.Global.Accesslog != "" {
		access_f, err = os.OpenFile(cfg.Global.Accesslog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err == nil {
			defer access_f.Close()
		} else {
			log.Printf("Opening access log %s failed: %v", cfg.Global.Accesslog, err)
		}
	}

	for _, item := range cfg.Sites {
		for _, host := range item.Domains {
			hosts[*host] = item.ProxyPass
		}
	}

	for _, item := range cfg.Sites {
		sites[item.Name] = &Site{
			Name:         item.Name,
			ListenAddr:   item.ListenAddr,
			Domains:      item.Domains,
			EnableSSL:    item.EnableSSL,
			AddForwarded: item.AddForwarded,
			KeyFile:      item.KeyFile,
			CertFile:     item.CertFile,
			ProxyPass:    item.ProxyPass,
		}
	}

	count := 0
	exit_chan := make(chan int)
	for name, frontend := range sites {
		log.Printf("Bind site [ %s ]", name)
		go func(site *Site, name string) {
			var accesslogger *log.Logger
			if access_f != nil {
				accesslogger = log.New(access_f, "frontend:"+name+" ", log.Ldate|log.Ltime|log.Lmicroseconds)
			}
			site.Start(hosts, accesslogger)
			exit_chan <- 1
		}(frontend, name)
		count++
	}

	for i := 0; i < count; i++ {
		<-exit_chan
	}
}

func usage() {
	fmt.Fprintf(os.Stdout, "usage: %s -config=<configfile>\n", os.Args[0])
	os.Exit(1)
}
