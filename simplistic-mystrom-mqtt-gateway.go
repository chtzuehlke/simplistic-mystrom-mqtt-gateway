package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"gopkg.in/resty.v1"
)

// Config holds all the AQS IoT properties
type Config struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	CaCert     string `json:"caCert"`
	ClientCert string `json:"clientCert"`
	PrivateKey string `json:"privateKey"`
}

func getSettingsFromFile(p string, opts *MQTT.ClientOptions) error {
	var conf, err = readFromConfigFile(p)
	if err != nil {
		return err
	}

	var tlsConfig, err2 = makeTLSConfig(conf.CaCert, conf.ClientCert, conf.PrivateKey)
	if err2 != nil {
		return err2
	}

	opts.SetTLSConfig(tlsConfig)

	var brokerURI = fmt.Sprintf("ssl://%s:%d", conf.Host, conf.Port)
	opts.AddBroker(brokerURI)

	return nil
}

func readFromConfigFile(path string) (Config, error) {
	var ret = Config{}

	var b, err = ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func makeTLSConfig(cafile, cert, key string) (*tls.Config, error) {
	var TLSConfig = &tls.Config{InsecureSkipVerify: false}

	var certPool *x509.CertPool
	var err error
	var tlsCert tls.Certificate

	certPool, err = getCertPool(cafile)
	if err != nil {
		return nil, err
	}

	TLSConfig.RootCAs = certPool

	certPool, err = getCertPool(cert)
	if err != nil {
		return nil, err
	}

	TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	TLSConfig.ClientCAs = certPool

	tlsCert, err = tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	TLSConfig.Certificates = []tls.Certificate{tlsCert}

	return TLSConfig, nil
}

func getCertPool(pemPath string) (*x509.CertPool, error) {
	var pemData, err = ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}

	var certs = x509.NewCertPool()
	certs.AppendCertsFromPEM(pemData)

	return certs, nil
}

// ArgOption holds command line arguments
type ArgOption struct {
	Conf            string
	ClientID        string
	MyStomSwitchIPs string
}

// NewOption creates new AWS IoT options (from a configuration file)
func NewOption(args *ArgOption) (*MQTT.ClientOptions, error) {
	var opts *MQTT.ClientOptions = MQTT.NewClientOptions()

	err := getSettingsFromFile(args.Conf, opts)
	if err != nil {
		return nil, err
	}

	opts.SetClientID(args.ClientID)
	opts.SetAutoReconnect(true)

	return opts, nil
}

// MyStromSwitchResponse holds mystrom switch status response (see https://api.mystrom.ch/?version=latest)
type MyStromSwitchResponse struct {
	Power float32 `json:"power"`
}

func getCurrentMyStromSwitchPower(url string) (float32, error) {
	resp, err := resty.R().Get(url)
	if err != nil {
		return -1, err
	}
	if resp.StatusCode() != 200 {
		return -1, errors.New("unexpected non-200 mystrom switch response code")
	}

	var myResp = MyStromSwitchResponse{}
	var parseErr = json.Unmarshal(resp.Body(), &myResp)
	if parseErr != nil {
		return -1, err
	}

	return myResp.Power, nil
}

var args ArgOption

func main() {
	flag.StringVar(&args.Conf, "conf", "", "Config file JSON path and name for accessing to AWS IoT endpoint")
	flag.StringVar(&args.ClientID, "client-id", "", "client id to connect with")
	flag.StringVar(&args.MyStomSwitchIPs, "mystrom-switch-ips", "", "mystrom switch IP addresses")
	flag.Parse()

	opts, err := NewOption(&args)
	if err != nil {
		panic(err)
	}

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer c.Disconnect(250) //rather pointless because of endless loop below

	topic := fmt.Sprintf("mystrom/power/%s", args.ClientID)

	//send power every minute (for every mystrom switch)
	for {
		var ips = strings.Split(args.MyStomSwitchIPs, ",")
		for _, ip := range ips {
			if ip == "" {
				continue
			}

			url := fmt.Sprintf("http://%s/report", ip)
			var power, err = getCurrentMyStromSwitchPower(url)
			if err != nil {
				panic(err) //FIXME Better way to deal with (temporary) problems ;-)
			}

			text := fmt.Sprintf("{\"ip\":\"%s\",\"power\":\"%f\"}", ip, power)
			token := c.Publish(topic, 0, false, text)
			token.Wait()

			fmt.Println(text)
		}

		time.Sleep(60 * time.Second)
	}
}
