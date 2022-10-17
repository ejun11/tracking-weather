package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WeatherData struct {
	Name  string `json:"name`
	Id    int    `json:"id`
	Coord struct {
		Lont float64 `json:"lon"`
		Latt float64 `json:"lat"`
	}

	Main struct {
		Tempp    float64 `json:"temp"`
		TemppMin float64 `json:"temp_min"`
		TemppMax float64 `json:"temp_max"`
	}

	Weather []struct {
		Weather     string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`

	Datetime int64 `json:"dt"`
	Timezone int64 `json:"timezone"`
}

type Timelocal struct {
	Name string `json:"name`
}

func query(city string) (WeatherData, error) {

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=8bd90625c19d08d8d1a1c20bd04b9ce0" + "&q=" + city + "&units=metric")
	if err != nil {
		return WeatherData{}, err
	}

	defer resp.Body.Close()

	var d WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}
	return d, nil
}

func timelocal(city string) (Timelocal, error) {

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=8bd90625c19d08d8d1a1c20bd04b9ce0" + "&q=" + city + "&units=metric")
	if err != nil {
		return Timelocal{}, err
	}

	defer resp.Body.Close()

	var e Timelocal
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return Timelocal{}, err
	}
	return e, nil
}

func main() {

	http.HandleFunc("/waktu/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := timelocal(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			now := time.Now().Local()
			tf := now.Format("2006-1-2 15:04:05")
			fmt.Println("GMT:", tf)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
			json.NewEncoder(w).Encode("GMT: " + tf)
		})

	http.HandleFunc("/tracking-weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	fmt.Println("Server jalan di: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
