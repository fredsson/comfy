package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const getUserPodIdsPath = "https://home.sensibo.com/api/v2/users/me/pods?fields=id&apiKey="
const getPodSmartModePath = "https://home.sensibo.com/api/v2/pods/{device_id}/smartmode?apiKey="

type SensiboProxy struct {
	apiKey string
}

type Pod struct {
	Id string
}
type PodsResponse struct {
	Status string
	Result []Pod
}

type TemperatureState struct {
	On                bool
	TargetTemperature int
	TemperatureUnit   string
	Mode              string
	FanLevel          string
	Swing             string
}

type SmartModeResult struct {
	Enabled                  bool
	Type                     string
	DevideUid                string
	LowTemperatureThreshold  float32
	HighTemperatureThreshold float32
	LowTemperatureState      TemperatureState
	HighTemperatureState     TemperatureState
}

type SmartModeResponse struct {
	Status string
	Result SmartModeResult
}

type mapper func(io.ReadCloser) interface{}

func (p SensiboProxy) FetchPods() []Pod {
	var podsMapper = func(body io.ReadCloser) interface{} {
		return mapToPodsResponse(body)
	}

	var response = Get(getUserPodIdsPath+p.apiKey, podsMapper).(*PodsResponse)
	return response.Result
}

func (p SensiboProxy) FetchSmartModeForPod(pod Pod) *SmartModeResult {
	var smartModeMapper = func(body io.ReadCloser) interface{} {
		return mapToSmartModeResponse(body)
	}
	var smartModePathWithDeviceId = strings.ReplaceAll(getPodSmartModePath, "{device_id}", pod.Id)
	var response = Get(smartModePathWithDeviceId+p.apiKey, smartModeMapper).(*SmartModeResponse)
	return &response.Result
}

func mapToPodsResponse(body io.ReadCloser) *PodsResponse {
	response := new(PodsResponse)
	err := decodeFromJson(body, response)
	if err != nil {
		log.Fatal("Could not decode Pods response")
	}

	return response
}

func mapToSmartModeResponse(body io.ReadCloser) *SmartModeResponse {
	response := new(SmartModeResponse)
	err := decodeFromJson(body, response)
	if err != nil {
		log.Fatal("Could not decode Smart mode response", err)
	}
	return response
}

func Get(url string, m mapper) interface{} {
	log.Println("Sending Get request to " + url)
	var resp, reqErr = http.Get(url)
	if reqErr != nil {
		log.Fatal("woopsi ", reqErr)
	}
	defer resp.Body.Close()

	return m(resp.Body)
}

func decodeFromJson(body io.ReadCloser, target interface{}) error {
	return json.NewDecoder(body).Decode(target)
}
