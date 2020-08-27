package main

import (
	"encoding/json"
	"github.com/hinoshiba/go-onamaeddns/src/onamaeddns"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type ResCh struct {
	IP     string
	Domain string
	Status bool
}

type Domain struct {
	Host   string
	Domain string
}

type JsonDate struct {
	Ip string
}

var Domains []Domain

func main() {
	envLoad()
	setDomains()

	ip, err := getRemoteIp()

	if err != nil {
		os.Exit(1)
	}

	limitCh := make(chan struct{}, len(Domains))
	defer close(limitCh)
	resCh := make(chan ResCh)
	defer close(resCh)

	for _, d := range Domains {
		go resolveIp(d.Host+d.Domain, limitCh, resCh)
	}
	domainEqual := false
	for i := 0; i < len(Domains); i++ {
		res := <-resCh
		if res.IP == ip && !domainEqual {
			domainEqual = true
		}
	}

	if !domainEqual {
		updateDns(ip)
	}
	for _, d := range Domains {
		go resolveIp(d.Host+d.Domain, limitCh, resCh)
	}
	for i := 0; i < len(Domains); i++ {
		res := <-resCh
		log.Println(res.IP)
	}

}

func resolveIp(domain string, limitCh chan struct{}, resCh chan<- ResCh) {
	limitCh <- struct{}{}
	addr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		resCh <- ResCh{"", domain, false}
		//os.Exit(1)
	} else {
		resCh <- ResCh{addr.String(), domain, true}
	}
	<-limitCh
}

func getRemoteIp() (string, error) {
	d := os.Getenv("get_ip_host")
	res, err := http.Get(d)
	if err != nil {
		log.Printf("URL%sが見つからなかった\n", d)
		return "", err
	}

	var jd JsonDate
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&jd)
	if err != nil {
		log.Printf("jsonが解析できなかった")
		return "", err
	}
	log.Println(jd.Ip)
	return jd.Ip, nil
}

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func updateDns(remoteIp string) {

	cl, err := onamaeddns.Dial(os.Getenv("onamae_server"), os.Getenv("username"), os.Getenv("password"), time.Minute)
	if err != nil {
		log.Println(err)
		return
	}
	defer cl.Close()

	for _, d := range Domains {
		if err := cl.UpdateIPv4(d.Host, d.Domain, remoteIp); err != nil {
			log.Println("tls error")
			log.Println(err)
			return
		}
	}
	log.Println("updated")
}

func setDomains() {
	hosts := strings.Split(os.Getenv("hosts"), ",")
	domain := os.Getenv("domain")
	for _, h := range hosts {
		log.Println(h)
		Domains = append(Domains, Domain{Host: h, Domain: domain})
	}
}
