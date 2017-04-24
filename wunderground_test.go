package main

import (
	"fmt"
	"testing"
)

func TestGetWeather(t *testing.T) {
	t.Skip("skip")
	wd := NewWundergroundProvider("")

	we, err := wd.getWeather("Tokyo", "JP")
	if err != nil {
		t.Fail()
	}

	fmt.Println(we)
}

func TestGet3DaysForecast(t *testing.T) {
	t.Skip("skip")
	wd := NewWundergroundProvider("")

	we, err := wd.get3DaysForcast("Tokyo", "JP")
	if err != nil {
		t.Fail()
	}

	fmt.Println(we)
}
