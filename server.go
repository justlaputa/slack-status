package main

import (
	"log"
	"os"
	"time"
)

const (
	UPDATEPERIOD = time.Hour * 1
)

func updateStatus(w WeatherProvider, slackApi *SlackApi) {
	log.Printf("updateing slack status...")

	weather, err := w.getWeather("Tokyo", "JP")
	if err != nil {
		log.Printf("failed to get weather, skip")
		return
	}

	currentEmoji := WeatherEmojiMap[weather]
	text := "Hello :earth_asia"

	weatherForecasts, err := w.get3DaysForcast("Tokyo", "JP")
	if err != nil {
		log.Printf("failed to get 3 days forecast, use boring text")
	} else {
		text = ""
		for _, wf := range weatherForecasts {
			text = text + string(WeatherEmojiMap[wf])
		}
	}

	result := slackApi.setStatus(text, currentEmoji)
	if !result {
		log.Printf("failed to set slack status, skip")
		return
	}
}

func main() {
	log.Printf("starting slack status...")
	log.Printf("changeing status in %d minutes", UPDATEPERIOD/60000000000)

	apiKey := os.Getenv("WUNDERGROUND_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("did you set api key of wunderground by environment variable WUNDERGROUND_API_KEY? I can not start without it, please set it and restart")
	}
	slackToken := os.Getenv("SLACK_API_TOKEN")
	if len(slackToken) == 0 {
		log.Fatal("did you set slack api token by environment variable SLACK_API_TOKEN? I can not start without it, please set it and restart")
	}

	weatherProvider := NewWundergroundProvider(apiKey)
	slackApi := &SlackApi{Token: slackToken}

	tickerChan := time.NewTicker(UPDATEPERIOD).C
	for {
		select {
		case <-tickerChan:
			go updateStatus(weatherProvider, slackApi)
		default:
			time.Sleep(time.Second)
		}
	}
}
