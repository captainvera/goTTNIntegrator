package main

import (
	"fmt"
	"github.com/captainvera/go-utils/mqtt"
	"os"
	"strings"
)

func main() {

	args := os.Args[1:]

	inputs := handleArgs(args)

	fmt.Println("Redirecting:" + strings.Join(inputs, " ") + " from ttn device.")

	//Config should also be an input argument.
	mqtt.ConfigMQTT("config")
	defer mqtt.CloseMQTT()

	messages := make(chan map[string]interface{})

	//Maybe also input argument? if wanted to withdraw data from several different apps
	go mqtt.ReadMQTT(messages)

	for {
		msg := <-messages

		fmt.Println()
		for _, in := range inputs {

			fmt.Printf("%s: ", in)
			fmt.Println(msg[in])
		}
	}
}

func handleArgs(args []string) []string {

	receivingInputs := false

	var inputs []string

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
			fmt.Println("-h     : -h        : print help")
			fmt.Println("-f     : -f [args] : specify arguments to redirect to OpenSensors.io")
			fmt.Println("---------------")
			os.Exit(0)

		case "-f":
			receivingInputs = true

		default:
			if receivingInputs == true && arg[0] != '-' {
				inputs = append(inputs, arg)
			} else {
				panic("Input:" + arg + " invalid. Try '-h'")
			}

		}

	}

	if len(inputs) < 1 {
		panic("Something needs to be redirected! Try '-h'")
	}

	return inputs
}
