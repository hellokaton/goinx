package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	PidFile     = "goinx.pid"
)

var (
	configPath = flag.String("config", "config.yml", "Configuration Path")
	cmds       = []string{"start", "stop", "restart"}
)

func usage() {
	fmt.Printf("💖 %s %s\n", ProjectName, Version)
	fmt.Println("Author: biezhi")
	fmt.Println("Github: https://github.com/biezhi/goinx")
	fmt.Println("\nUsage: goinx [start|stop|restart]\n")
	fmt.Println("Options:")
	fmt.Println("    --config\tConfiguration path")
	fmt.Println("    --help\tHelp info")
}

func startArgs() *Config {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	if !Contains(cmds, cmd) {
		usage()
		os.Exit(0)
	}

	// start goinx
	if cmd == cmds[0] {
		return start()
	}
	// stop goinx
	if cmd == cmds[1] {
		stop()
	}
	if cmd == cmds[2] {
		stop()
		return start()
	}

	return nil
}

func start() *Config {

	if Exist(PidFile) {
		log.Warning("Goinx has bean started.")
		os.Exit(0)
	}

	conf := Config{}
	if pid := syscall.Getpid(); pid != 1 {
		err := ioutil.WriteFile(PidFile, []byte(strconv.Itoa(pid)), 0777)
		if err != nil {
			fmt.Println(err)
		}
	}

	flag.Usage = usage
	flag.Parse()
	bytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Error("%v", err)
		os.Exit(0)
	}
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		log.Error("%v", err)
		os.Exit(0)
	}
	return &conf
}

func stop() {
	bytes, err := ioutil.ReadFile(PidFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	pid, err := strconv.Atoi(string(bytes))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	syscall.Kill(pid, 9)
	os.Remove("goinx.pid")
}

func shutdownHook() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				os.Remove("goinx.pid")
				log.Info("Shutown Goinx.")
				os.Exit(0)
			default:
				log.Info("other", s)
			}
		}
	}()
}

func main() {

	shutdownHook()

	conf := startArgs()
	log.Info("Start Goinx.")

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
		go func(s Server) {
			s.Start()
			exit_chan <- 1
		}(server)
		count++
	}

	for i := 0; i < count; i++ {
		<-exit_chan
	}

}
