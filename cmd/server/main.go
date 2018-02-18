package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/andreipimenov/kvstore/config"
	"github.com/andreipimenov/kvstore/store"
)

func main() {
	configFile := flag.String("config", "", "configuration file")
	port := flag.Int("port", -1, "server port")
	flag.Parse()

	c, err := NewConfig(config.New(*configFile))
	if *configFile != "" && err != nil {
		log.Println(err)
	}
	if *port >= 0 {
		c.Port = *port
	}

	s := NewStore(store.New(c.DumpFile, c.DumpInterval))

	r := NewRouter(c, s)

	log.Printf("Start listening on port %d", c.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
