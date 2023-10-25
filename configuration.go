package main

import "os"

type Configuration struct {
	Home string
}

func getConfiguration() Configuration {
	return Configuration{Home: os.Getenv("ZET_HOME")}
}
