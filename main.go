package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/weather", getWeather)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getWeather(w http.ResponseWriter, r *http.Request) {
    // Fetch city from request query parameters
    city := r.URL.Query().Get("city")

    // Get API key from environment variable
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        http.Error(w, "API key not set", http.StatusInternalServerError)
        return
    }

    // Construct OpenWeatherMap API URL
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=imperial", city, apiKey)

    // Make API call
    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Parse JSON response
    var weatherData map[string]interface{}
    err = json.Unmarshal(body, &weatherData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send response back to client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(weatherData)
}
