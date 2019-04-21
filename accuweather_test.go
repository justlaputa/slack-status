package main

import (
	"fmt"
	"testing"
)

func TestCurrentCondition(t *testing.T) {
	t.Skip("skip")
	ac := NewAccuWeatherProvider("")

	w, err := ac.getWeather("Tokyo", "JP")
	if err != nil {
		t.Fail()
	}

	fmt.Println(w)
}

func TestForcast(t *testing.T) {
	t.Skip("skip")
	ac := NewAccuWeatherProvider("")

	w, err := ac.get3DaysForcast("Tokyo", "JP")
	if err != nil {
		t.Fail()
	}

	fmt.Println(w)
}
