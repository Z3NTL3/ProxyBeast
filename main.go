package main

import (
	"Z3NTL3/ProxyBeast/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}