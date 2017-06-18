package main

import (
	"fmt"
	"github.com/captainvera/go-utils/mqtt"
	"os"
)

func main() {

	args := os.Args[1:]

	fmt.Println(args)

	receivingInputs := false
	inputs := []string{}

	for _, arg := range args {
		switch arg {
		case "-h":
			if receivingInputs {
				receivingInputs = false
			}
			fmt.Println("Welcome to goTTNIntegrator!")
			fmt.Println("You are expected to have a 'config' file in this directory.")
			fmt.Println("Refer to README")
			fmt.Println("---------------")
			fmt.Println("-h 	: -h 		: print help")
			fmt.Println("-f     : -f [args] : specify arguments to redirect to OpenSensors.io")
			fmt.Println("---------------")
			os.Exit(0)

		case "-f":
			receivingInputs = true

		default:
			fmt.Println("Received: " + arg)
			if receivingInputs == true {
				inputs = append(inputs, arg)
			} else {
				fmt.Println("Input:" + arg + " invalid. Try '-h'")
			}

		}

	}

	fmt.Println(inputs)

	if len(inputs) < 1 {
		panic("Something needs to be redirected!")
	}

	mqtt.ConnectMQTT()
	defer mqtt.CloseMQTT()

	messages := make(chan map[string]interface{})

	go mqtt.ReadMQTT(messages)

	for {
		msg := <-messages

		fmt.Println(msg)
	}

}
