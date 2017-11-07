package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	AccessLog string `yaml:"access_log"`
	Http      struct {
		Servers []Server `yaml:"servers,flow"`
	}
}

var (
	configPath = flag.String("config", "sample/config.yml", "Configuration Path")
)

func usage() {
	fmt.Println("Goinx 0.0.1")
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
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(conf)

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
