package main

import (
	"os"

	"./webserver_sl"
)

func main() {
	if len(os.Args) < 3 {
		panic("Missing command line arguments. Argument 1 should be the listening Address and Argument 2 should be the path to the request-file mapping.")
	}
	listenAddress := os.Args[1]
	requestFileMapPath := os.Args[2]
	webserver_sl.Start(listenAddress, requestFileMapPath)
}
