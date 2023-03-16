package main

import (
	"fmt"
	"os"

	"github.com/andreykaipov/goobs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: obs_recorder start|stop")
		os.Exit(1)
	}

	action := os.Args[1]

	client, err := goobs.New("localhost:4455", goobs.WithPassword("qTh2WVY6SSC3JNmy"))
	if err != nil {
		fmt.Println("Error connecting to OBS:", err)
		os.Exit(1)
	}

	switch action {
	case "start":
		err = startRecording(client)
	case "stop":
		err = stopRecording(client)
	default:
		fmt.Println("Invalid action. Usage: obs_recorder start|stop")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error executing action:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully executed action:", action)
}

func startRecording(client *goobs.Client) error {
	_, err := client.Record.StartRecord()
	return err
}

func stopRecording(client *goobs.Client) error {
	_, err := client.Record.StopRecord()
	return err
}
