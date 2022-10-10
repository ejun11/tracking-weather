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

	Datetime int64     `json:"dt"`
	Timezone int64     `json:"timezone"`
	DTime    time.Time `json:"ConvertUnix"`
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
	converttoDate(d)
	return d, nil
}

func converttoDate(d WeatherData) {
	x := d.Datetime
	now := time.Now()
	ConvertUnix := time.Unix(x, 0)

	fmt.Println("\nCity:", d.Name)
	fmt.Println("Convert GMT:", ConvertUnix)
	fmt.Println("Local Time:", now)
}

func main() {

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
