package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const TIBBER_URL = "api.tibber.com"

const getUserPodIdsPath = "https://home.sensibo.com/api/v2/users/me/pods?fields=id&apiKey="
const getPodSmartModePath = "https://home.sensibo.com/api/v2/pods/{device_id}/smartmode?apiKey="
const setAcStatePath = "/pods/{podUid}/acStates"

type Pod struct {
	Id string
}
type PodResult struct {
	Status string
	Result []Pod
}

func decodeFromJson(body io.ReadCloser, target interface{}) error {
	return json.NewDecoder(body).Decode(target)
}

func mapToPodResult(body io.ReadCloser) *PodResult {
	result := new(PodResult)
	err := decodeFromJson(body, result)
	if err != nil {
		log.Fatal("Could not decode pod response")
	}

	return result
}

func main() {
	var envErr error = godotenv.Load()
	if envErr != nil {
		log.Fatalf("Some error occured. Err: %s", envErr)
	}

	var sensiboApiKey string = os.Getenv("SENSIBO_API_KEY")

	var resp, reqErr = http.Get(getUserPodIdsPath + sensiboApiKey)
	if reqErr != nil {
		log.Fatal("woopsi ", reqErr)
	}
	defer resp.Body.Close()

	var result = mapToPodResult(resp.Body)

	fmt.Println(result.Status)
}

// get /pods/{device_id}/smartmode

// put /pods/{device_id}/smartmode
/*
{
  "enabled": false
}
*/

// post /pods/{device_id}/smartmode
/*
{
  "enabled": false,
  "lowTemperatureThreshold": 0,
  "lowTemperatureState": {},
  "highTemperatureThreshold": 0,
  "highTemperatureState": {}
}
*/
