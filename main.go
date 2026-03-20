package main

import (
	server "dcrcs-go/pkg"
)

func main() {
	server := server.Server()
	server.Run()
}
