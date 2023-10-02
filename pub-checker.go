package main

// This is where we connect to an API to check our IP
// and we then publish is to MQTT
// No other action performed herepa

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	Mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ipAddr struct {
	IP string `json:"ip"`
}

var messagePubHandler Mqtt.MessageHandler = func(client Mqtt.Client, msg Mqtt.Message) {
	fmt.Printf("%s - %s : %s\n", curtime, msg.Topic(), msg.Payload())
}

var connectHandler Mqtt.OnConnectHandler = func(client Mqtt.Client) {
	fmt.Printf("%s - Connected to MQTT\n", curtime)
}

var connectLostHandler Mqtt.ConnectionLostHandler = func(client Mqtt.Client, err error) {
	fmt.Printf("%s - Connection to MQTT lost: %v\n", curtime, err)
}

func mqttpub(client Mqtt.Client, message string) {
	text := fmt.Sprint(message)
	token := client.Publish(Mqtttopic, 1, true, text)
	token.Wait()
	time.Sleep(time.Second)
}

var curtime = time.Now().Format(time.RFC3339)

func main() {

	//
	// Connect to IP API & grab the IP
	//
	req, err := http.NewRequest(http.MethodGet, ApiSite, nil)
	if err != nil {
		log.Fatalf("Error %v - IP API connection failed", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error %v - IP API JSON request", err)
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error %v - Error reading IP API HTTP body", err)
	}

	var ip ipAddr

	err = json.Unmarshal(responseData, &ip)
	if err != nil {
		log.Fatalf("Error %v - Error unmarshaling JSON", err)
	}

	if len(ip.IP) == 0 {
		log.Fatalf("%s: Unable to collect IP from API, exiting...", curtime)
	} else {
		fmt.Printf("%s: Pulled IP from API is: %s\n", curtime, ip.IP)
	}

	//Init MQTT broker and connect
	opts := Mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", Mqtthost, Mqtthostport))
	opts.SetClientID(Mqttclid)
	opts.SetUsername(Mqttuser)
	opts.SetPassword(Mqttpwd)
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	myclient := Mqtt.NewClient(opts)
	if token := myclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Publish retrieved IP to MQTT
	mqttpub(myclient, ip.IP)

	myclient.Disconnect(15)
}
