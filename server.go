package main

import (
	"log"
	"time"
)

const (
	UPDATEPERIOD = time.Minute * 1
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

	weatherProvider := NewWundergroundProvider("")
	slackApi := &SlackApi{Token: ""}

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
