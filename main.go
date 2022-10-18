package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"time"
)

type apiConfigData struct {
	AbstractApiKey string `json:"AbstractApiKey"`
}

type JSONTime time.Time

type timeData struct {
	Country string `json:"requested_location"`
	JSONTime string `json:"datetime"`
	GMT int `json:"gmt_offset"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func query(city string)(timeData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return timeData{}, err
	}
	
	resp, err := http.Get("https://timezone.abstractapi.com/v1/current_time/?api_key=" + apiConfig.AbstractApiKey + "&location=" + city)

	if err != nil {
		return timeData{}, err
	}

	defer resp.Body.Close()

	var d timeData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return timeData{}, err
	}
	return d, nil
}

func main() {
	http.HandleFunc("/time/",
	func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset = utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}