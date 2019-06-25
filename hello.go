package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	Conf     string
	ClientID string
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

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var args ArgOption

///

//curl -H "Content-Type: text/json" -d '{"power":120,"relay":"FIXME","temperature":"FIXME"}' https://httpecho-167219.appspot.com/register/get/200/mystrom/switch/sampleresponse1.json
//curl -v https://httpecho-167219.appspot.com/mystrom/switch/sampleresponse1.json | jq

// MyStromgSwitchResponse holds mystrom switch status response (see https://api.mystrom.ch/?version=latest)
type MyStromgSwitchResponse struct {
	Power       int    `json:"power"`
	Relay       string `json:"relay"`
	Temperature string `json:"temperature"`
}

///

func main() {

	///

	resp, err := resty.R().Get("https://httpecho-167219.appspot.com/mystrom/switch/sampleresponse1.json")
	if err != nil || resp.StatusCode() != 200 {
		panic(err)
	}

	// explore response object
	//fmt.Printf("\nError: %v", err)
	//fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	//fmt.Printf("\nResponse Status: %v", resp.Status())
	//fmt.Printf("\nResponse Time: %v", resp.Time())
	//fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	//fmt.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())

	var myResp = MyStromgSwitchResponse{}
	var parseErr = json.Unmarshal(resp.Body(), &myResp)
	if parseErr != nil {
		panic(err)
	}
	fmt.Printf("\nResponse Body power: %v\n", myResp.Power)

	///

	flag.StringVar(&args.Conf, "conf", "", "Config file JSON path and name for accessing to AWS IoT endpoint")
	flag.StringVar(&args.ClientID, "client-id", "", "client id to connect with")
	flag.Parse()

	opts, err := NewOption(&args)
	if err != nil {
		panic(err)
	}

	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Disconnect(250)
}
