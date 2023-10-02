package main

import (
	"os"
	"strconv"
)

// From OS environment import:
//   - your telegram bot API key
//   - the target user (or channel)
//   - the website which can return a JSON string with your IP address in format
//     {"ip":"actual_IP_address"}
//   - name or address of the MQTT server
var Tgbotkey = os.Getenv("MYIP2TG_APIKEY")
var Tgtarget, _ = strconv.Atoi(os.Getenv("MYIP2TG_TARGET"))
var ApiSite = os.Getenv("MYIP2TG_IPSITE")
var Mqtthost = os.Getenv("MYIP2TG_MQTT")
var Mqtthostport, _ = strconv.Atoi(os.Getenv("MYIP2TG_MQTTPORT"))
var Mqttuser = os.Getenv("MYIP2TG_MQTT_USER")
var Mqttpwd = os.Getenv("MYIP2TG_MQTT_PWD")
var Mqtttopic = os.Getenv("MYIP2TG_MQTT_TOPIC")
var Mqttclid = os.Getenv("MYIP2TG_MQTT_CLID")
