package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	log "github.com/biezhi/goinx/log"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	AccessLog string `yaml:"access_log"`
	LogLevel  string `yaml:"log_level"`
	Http      struct {
		Servers []Server `yaml:"servers,flow"`
	}
}

const (
	ProjectName = "Goinx"
	Version     = "0.0.1"
)

var (
	configPath = flag.String("config", "sample/config.yml", "Configuration Path")
)

func usage() {
	fmt.Printf("%s %s\n", ProjectName, Version)
	fmt.Println("Usage: goinx --config=<configfile>\n")
	fmt.Println("Options:")
	fmt.Println("\t--config\tConfiguration Path")
}

func main() {

	flag.Usage = usage
	flag.Parse()

	conf := Config{}
	bytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Error("%v", err)
		return
	}
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		log.Error("%v", err)
		return
	}

	if conf.LogLevel == "debug" {
		log.LogLevelNum = 1
	}
	if conf.LogLevel == "info" {
		log.LogLevelNum = 2
	}
	if conf.LogLevel == "warn" {
		log.LogLevelNum = 3
	}
	if conf.LogLevel == "error" {
		log.LogLevelNum = 4
	}

	log.Debug("Config Content: %v", conf)

	count := 0
	exit_chan := make(chan int)
	for _, server := range conf.Http.Servers {
		go func(server Server) {
			server.Start()
			exit_chan <- 1
		}(server)
		count++
	}

	for i := 0; i < count; i++ {
		<-exit_chan
	}

}
