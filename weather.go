package main

type WeatherCondition uint8

const (
	CLEAR WeatherCondition = iota
	RAIN
	PARTLY_CLOUDY
	MOSTLY_CLOUDY
	THUNDER
	THUNDER_RAIN
	TORNADO
	SNOW
	WINDY
	FOGGY
	HAZE
	UNKNOWN_WEATHER
)

type WeatherProvider interface {
	getWeather(string, string) (WeatherCondition, error)
	get3DaysForcast(string, string) ([]WeatherCondition, error)
}
