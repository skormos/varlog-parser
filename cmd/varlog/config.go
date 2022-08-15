package main

import (
	"flag"
)

type config struct {
	dirPath string
	http    httpConfig
}

type httpConfig struct {
	port int
}

func parseFlags() config {
	logPath := flag.String("logPath", "/var/log", "Tells the service where to look for requested files. Default is `/var/log`.")
	httpPort := flag.Int("httpPort", 8080, "The port on which the http server will listen. Default is 8080.")

	flag.Parse()

	return config{
		dirPath: *logPath,
		http: httpConfig{
			port: *httpPort,
		},
	}
}
