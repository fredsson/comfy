package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const getUserPodIdsPath = "https://home.sensibo.com/api/v2/users/me/pods?fields=id&apiKey="
const podSmartModePath = "https://home.sensibo.com/api/v2/pods/{device_id}/smartmode?apiKey="

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
	On                bool   `json:"on"`
	TargetTemperature int    `json:"targetTemperature"`
	TemperatureUnit   string `json:"temperatureUnit"`
	Mode              string `json:"mode"`
	FanLevel          string `json:"fanLevel"`
	Swing             string `json:"swing"`
}

type SmartModeResult struct {
	Enabled                  bool   `json:"enabled"`
	Type                     string `json:"type"`
	DevideUid                string
	LowTemperatureThreshold  float32          `json:"lowTemperatureThreshold"`
	HighTemperatureThreshold float32          `json:"highTemperatureThreshold"`
	LowTemperatureState      TemperatureState `json:"lowTemperatureState"`
	HighTemperatureState     TemperatureState `json:"highTemperatureState"`
}

type SmartModeResponse struct {
	Status string
	Result SmartModeResult
}

type SmartModeRequest struct {
	Enabled                  bool             `json:"enabled"`
	Type                     string           `json:"type"`
	LowTemperatureThreshold  float32          `json:"lowTemperatureThreshold"`
	LowTemperatureState      TemperatureState `json:"lowTemperatureState"`
	HighTemperatureThreshold float32          `json:"highTemperatureThreshold"`
	HighTemperatureState     TemperatureState `json:"highTemperatureState"`
}

func getDefaultSmartModeRequest() SmartModeRequest {
	return SmartModeRequest{
		Enabled:                 false,
		Type:                    "temperature",
		LowTemperatureThreshold: 19,
		LowTemperatureState: TemperatureState{
			On:                true,
			TargetTemperature: 19,
			TemperatureUnit:   "C",
			Mode:              "auto",
			FanLevel:          "auto",
			Swing:             "rangeFull",
		},
		HighTemperatureThreshold: 22,
		HighTemperatureState: TemperatureState{
			On:                false,
			TargetTemperature: 22,
			TemperatureUnit:   "C",
			Mode:              "auto",
			FanLevel:          "low",
			Swing:             "stopped",
		},
	}
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
	var smartModePathWithDeviceId = strings.ReplaceAll(podSmartModePath, "{device_id}", pod.Id)
	var response = Get(smartModePathWithDeviceId+p.apiKey, smartModeMapper).(*SmartModeResponse)
	return &response.Result
}

func (p SensiboProxy) EnableSmartMode(pod Pod) {
	log.Println("Enabling smart mode!")
	p.SetSmartMode(pod, true)
}

func (p SensiboProxy) DisableSmartMode(pod Pod) {
	log.Println("Disabling smart mode!")
	p.SetSmartMode(pod, false)
}

func (p SensiboProxy) SetSmartMode(pod Pod, enabled bool) {
	var mapper = func(body io.ReadCloser) interface{} {
		var b, err = io.ReadAll(body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
		return nil
	}

	var request SmartModeRequest = getDefaultSmartModeRequest()
	request.Enabled = enabled

	var smartModePathWithDeviceId = strings.ReplaceAll(podSmartModePath, "{device_id}", pod.Id)
	var body, marshalError = json.Marshal(request)
	if marshalError != nil {
		log.Fatal("Could not marshal smart mode request")
	}
	Post(smartModePathWithDeviceId+p.apiKey, strings.NewReader(string(body)), []HeaderDefinition{}, mapper)
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
