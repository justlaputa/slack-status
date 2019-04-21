package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	previousEmoji   Emoji = EMOJI_QUESTION
	previousWeather       = []WeatherCondition{UNKNOWN_WEATHER, UNKNOWN_WEATHER, UNKNOWN_WEATHER}
)

func updatePrevious(emoji Emoji, forecasts []WeatherCondition) (Emoji, []WeatherCondition) {
	if emoji == EMOJI_QUESTION {
		emoji = previousEmoji
	} else {
		previousEmoji = emoji
	}

	for i, w := range forecasts {
		if w == UNKNOWN_WEATHER {
			forecasts[i] = previousWeather[i]
		} else {
			previousWeather[i] = forecasts[i]
		}
	}

	return emoji, forecasts
}

func updateStatus(w WeatherProvider, slackAPI *SlackApi) error {
	log.Printf("updateing slack status...")

	weather, err := w.getWeather("Tokyo", "JP")
	if err != nil {
		log.Printf("failed to get weather, skip")
		return err
	}

	currentEmoji := WeatherEmojiMap[weather]
	text := "Hello :earth_asia"

	weatherForecasts, err := w.get3DaysForcast("Tokyo", "JP")
	if err != nil {
		log.Printf("failed to get 3 days forecast, use boring text")
	} else {
		text = "|"
		currentEmoji, weatherForecasts = updatePrevious(currentEmoji, weatherForecasts)

		for _, wf := range weatherForecasts {
			text = text + string(WeatherEmojiMap[wf])
		}
	}

	result := slackAPI.setStatus(text, currentEmoji)
	if !result {
		log.Printf("failed to set slack status, skip")
		return err
	}

	return nil
}

func main() {
	apiKey := os.Getenv("ACCUWEATHER_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("did you set api key of accuweather by environment variable ACCUWEATHER_API_KEY? I can not start without it, please set it and restart")
	}
	slackToken := os.Getenv("SLACK_API_TOKEN")
	if len(slackToken) == 0 {
		log.Fatal("did you set slack api token by environment variable SLACK_API_TOKEN? I can not start without it, please set it and restart")
	}

	weatherProvider := NewAccuWeatherProvider(apiKey)
	slackAPI := &SlackApi{Token: slackToken}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := updateStatus(weatherProvider, slackAPI)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "success")
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting slack status on: %s", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
