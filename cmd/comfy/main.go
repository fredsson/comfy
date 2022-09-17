package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var envErr error = godotenv.Load()
	if envErr != nil {
		log.Fatalf("Some error occured. Err: %s", envErr)
	}

	sensiboProxy := SensiboProxy{
		apiKey: os.Getenv("SENSIBO_API_KEY"),
	}
	var pods = sensiboProxy.FetchPods()
	fmt.Println(pods[0].Id)

	var smartMode = sensiboProxy.FetchSmartModeForPod(pods[0])
	fmt.Println(smartMode.Enabled)
}
