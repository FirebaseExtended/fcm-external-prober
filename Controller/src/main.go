package main

import (
	"controller"
	"flag"
	"io/ioutil"
	"log"
)

func main() {
	cf := flag.String("config", "config.txt", "text file in which a protobuf with config information is located")
	c, err := ioutil.ReadFile(*cf)
	if err != nil {
		log.Fatalf("Main: could not read from specified config file: %s", err.Error())
	}
	ctrl := controller.NewController(string(c))
	ctrl.StartVMs()
	ctrl.StartProbes()
}