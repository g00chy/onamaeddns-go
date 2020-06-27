package network

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type OnameId struct {
	UserName string
	password string
}

type Url struct {
	getIpHost    string
	onamaeHost   string
	targetDomain string
}

type Ips struct {
	currentIp string
	domainIp  []string
}

type Network struct {
	OnameId *OnameId
	Url     *Url
	Ips     *Ips
}

func (n2 Network) CheckIp() {
	n2.getIp()
	n2.getIpOfDomain()
	log.Print(n2.Ips.domainIp)
	log.Print(n2.Ips.currentIp)
}

var network = Network{}

type jsondata struct {
	Ip string `json:"ip"`
}

func (n2 *Network) getIp() error {
	res, err := http.Get(n2.Url.getIpHost)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("StatusCode=%d", res.StatusCode)
		return errors.New("error")
	}
	// jsonを読み込む
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// jsonを構造体へデコード
	var jsonData jsondata
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return err
	}

	n2.Ips.currentIp = jsonData.Ip
	return nil
}

func (n2 *Network) getIpOfDomain() error {
	addr, err := net.LookupHost(n2.Url.targetDomain)
	if err != nil {
		return err
	}
	n2.Ips.domainIp = addr

	return nil
}

func NewNetwork() *Network {
	var id = OnameId{}
	var url = Url{}
	var ips = Ips{}

	network.OnameId = &id
	network.Url = &url
	network.Ips = &ips

	network.loadEnv()

	return &network
}

func (n2 *Network) loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	network.OnameId.UserName = os.Getenv("USER_NAME")
	network.OnameId.password = os.Getenv("PASSWORD")
	network.Url.getIpHost = os.Getenv("GET_IP_ADDR_HOST")
	network.Url.onamaeHost = os.Getenv("ONAMEDNS_HOST")
	var domains []string
	domains[0] = os.Getenv("TARGET_DOMAIN")
}
