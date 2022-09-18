package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var envErr error = godotenv.Load()
	if envErr != nil {
		log.Fatalf("Some error occured. Err: %s", envErr)
	}

	// var sensiboProxy = SensiboProxy{
	// 	apiKey: os.Getenv("SENSIBO_API_KEY"),
	// }
	// var pods = sensiboProxy.FetchPods()
	// fmt.Println(pods[0].Id)

	// var smartMode = sensiboProxy.FetchSmartModeForPod(pods[0])
	// fmt.Println(smartMode.Enabled)

	var tibberProxy = TibberProxy{
		apiKey: os.Getenv("TIBBER_API_KEY"),
	}

	var prices = tibberProxy.FetchPricesToday()
	log.Println(prices)
	for {
		// check price for current hour
		// if cheap -> enable smartMode (bool in ram to see if it's enabled already or not)
		// if expensive -> disable smartMode
		time.Sleep(time.Hour) // TODO: calculate how much to sleep
	}
}
