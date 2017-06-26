package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ttnConfiguration struct {
	AppId  string
	ApiKey string
}

type osConfiguration struct {
	ApiKey   string
	ClientId int
	Password string
	Topic    string
}

type configuration struct {
	TTN         ttnConfiguration
	OpenSensors osConfiguration
}

var log ttnlog.Interface

func main() {
	log = apex.Stdout()
	ttnlog.Set(log)

	args := os.Args[1:]

	inputs, filepath := handleArgs(args)

	if len(inputs) > 0 {
		log.Info("goTTNIntegrator: Redirecting [" + strings.Join(inputs, ", ") + "] from ttn device.")
	}

	config := ttnsdk.NewCommunityConfig("goTTNIntegrator")
	config.ClientVersion = "0.1"

	fileConfig := readConfig(filepath)

	client := config.NewClient(fileConfig.TTN.AppId, fileConfig.TTN.ApiKey)
	defer client.Close()

	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatal("goTTNIntegrator: Could not get application pub/sub")
	}

	allDevices := pubsub.AllDevices()
	uplink, err := allDevices.SubscribeUplink()
	if err != nil {
		log.WithError(err).Fatal("goTTNIntegrator: Could not subscribe to uplink messages")
	}

	for {
		for message := range uplink {
			payload := message.PayloadFields
			log.WithFields(payload).Info("goTTNIntegrator: Received uplink")
			if len(inputs) != 0 {
				payload = filterPayload(payload, inputs)
			}
			go forwardPayload(payload, fileConfig.OpenSensors)
		}
	}
}

func filterPayload(payload map[string]interface{}, fields []string) map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.WithField("panic", err).Error("goTTNIntegrator: Error when filtering payload")
		}
	}()

	res := make(map[string]interface{})

	for _, arg := range fields {
		res[arg] = payload[arg]
	}

	log.WithFields(res).Info("goTTNIntegrator: Filtered payload contents before redirecting")

	return res
}

func forwardPayload(payload map[string]interface{}, config osConfiguration) {
	url := "https://realtime.opensensors.io/v1/topics/" + config.Topic + "?client-id=" + strconv.Itoa(config.ClientId) + "&password=" + string(config.Password)

	data, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("goTTNIntegrator: Couldn't format payload data")
	}
	datamap := map[string]string{"data": string(data)}

	data2, err := json.Marshal(datamap)
	if err != nil {
		log.WithError(err).Error("goTTNIntegrator: Couldn't format payload data")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data2))
	if err != nil {
		log.WithError(err).Error("goTTNIntegrator: Couldn't create http request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "api-key "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("goTTNIntegrator: Error sending POST request to OpenSensors.io")
	} else {
		defer resp.Body.Close()
	}

	log.WithField("Response", resp.Status).Info("goTTNIntegrator: Received response from OpenSensors.io")
}

func handleArgs(args []string) ([]string, string) {
	var inputs []string
	var filepath string

	for num := 0; num < len(args); num++ {

		switch args[num] {
		case "-h":
			fmt.Println("Welcome to goTTNIntegrator!")
			fmt.Println("Refer to README")
			fmt.Println("---------------")
			fmt.Println("-h     : -h        : print help")
			fmt.Println("-f     : -f [args] : specify specific arguments to redirect to OpenSensors.io, if not used, entire payload will be redirected")
			fmt.Println("-c     : -c path   : specify path to config file, default is \"config\"")
			fmt.Println("---------------")
			os.Exit(0)

		case "-f":
			num += 1
			for num < len(args) && args[num][0] != '-' {
				inputs = append(inputs, args[num])
				num += 1
			}

		case "-c":
			num += 1
			filepath = args[num]

		default:
			log.Fatal("goTTNIntegrator: Input:" + args[num] + " invalid. Try '-h'")
		}

	}

	if filepath != "" {
		return inputs, filepath
	}
	return inputs, "config"
}

func readConfig(path string) configuration {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithError(err).Fatal("goTTNIntegrator: Could not get read config file on specified path")
	}

	config := configuration{}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		log.WithError(err).Fatal("goTTNIntegrator: Could not get read config file json")
	}
	data, _ = json.Marshal(config)

	log.WithField("data", string(data)).Debug("goTTNIntegrator: Config successfully read")

	return config
}
