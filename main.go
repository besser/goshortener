package main

import (
	"runtime"

	"github.com/besser/goshortener/config"
	"github.com/besser/goshortener/server"
)

func init() {
	config.LoadConfig()
}

func main() {
	// Setting the limits the number of operating system threads that can execute user-level Go code simultaneously.
	if numprocs := config.Cfg.GetInt("cpu.numprocs"); numprocs != 0 {
		runtime.GOMAXPROCS(numprocs)
	}

	// Run server
	server.Run()
}