package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ConditionResponseType struct {
	Response           ResponseType           `json:"response"`
	CurrentObservation CurrentObservationType `json:"current_observation"`
}

type ResponseType struct {
	Version        string       `json:"version"`
	TermsofService string       `json:"termsofService"`
	Features       FeaturesType `json:"features"`
	Error          ErrorType    `json:"error"`
}

type FeaturesType struct {
	Conditions float64 `json:"conditions,omitempty"`
	Forecast   float64 `json:"forecast,omitempty"`
}

type ErrorType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type EstimatedType interface{}

type CurrentObservationType struct {
	Image                 ImageType     `json:"image"`
	DisplayLocation       LocationType  `json:"display_location"`
	ObservationLocation   LocationType  `json:"observation_location"`
	Estimated             EstimatedType `json:"estimated"`
	StationId             string        `json:"station_id"`
	ObservationTime       string        `json:"observation_time"`
	ObservationTimeRfc822 string        `json:"observation_time_rfc822"`
	ObservationEpoch      string        `json:"observation_epoch"`
	LocalTimeRfc822       string        `json:"local_time_rfc822"`
	LocalEpoch            string        `json:"local_epoch"`
	LocalTzShort          string        `json:"local_tz_short"`
	LocalTzLong           string        `json:"local_tz_long"`
	LocalTzOffset         string        `json:"local_tz_offset"`
	Weather               string        `json:"weather"`
	TemperatureString     string        `json:"temperature_string"`
	TempF                 float64       `json:"temp_f"`
	TempC                 float64       `json:"temp_c"`
	RelativeHumidity      string        `json:"relative_humidity"`
	WindString            string        `json:"wind_string"`
	WindDir               string        `json:"wind_dir"`
	WindDegrees           float64       `json:"wind_degrees"`
	WindMph               float64       `json:"wind_mph"`
	WindGustMph           float64       `json:"wind_gust_mph"`
	WindKph               float64       `json:"wind_kph"`
	WindGustKph           float64       `json:"wind_gust_kph"`
	PressureMb            string        `json:"pressure_mb"`
	PressureIn            string        `json:"pressure_in"`
	PressureTrend         string        `json:"pressure_trend"`
	DewpointString        string        `json:"dewpoint_string"`
	DewpointF             float64       `json:"dewpoint_f"`
	DewpointC             float64       `json:"dewpoint_c"`
	// HeatIndexString       string        `json:"heat_index_string"`
	HeatIndexF        string `json:"heat_index_f"`
	HeatIndexC        string `json:"heat_index_c"`
	WindchillString   string `json:"windchill_string"`
	WindchillF        string `json:"windchill_f"`
	WindchillC        string `json:"windchill_c"`
	FeelslikeString   string `json:"feelslike_string"`
	FeelslikeF        string `json:"feelslike_f"`
	FeelslikeC        string `json:"feelslike_c"`
	VisibilityMi      string `json:"visibility_mi"`
	VisibilityKm      string `json:"visibility_km"`
	Solarradiation    string `json:"solarradiation"`
	UV                string `json:"UV"`
	Precip1hrString   string `json:"precip_1hr_string"`
	Precip1hrIn       string `json:"precip_1hr_in"`
	Precip1hrMetric   string `json:"precip_1hr_metric"`
	PrecipTodayString string `json:"precip_today_string"`
	PrecipTodayIn     string `json:"precip_today_in"`
	PrecipTodayMetric string `json:"precip_today_metric"`
	Icon              string `json:"icon"`
	IconUrl           string `json:"icon_url"`
	ForecastUrl       string `json:"forecast_url"`
	HistoryUrl        string `json:"history_url"`
	ObUrl             string `json:"ob_url"`
	Nowcast           string `json:"nowcast"`
}

type ImageType struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type LocationType struct {
	Full           string `json:"full, omitempty"`
	City           string `json:"city, omitempty"`
	State          string `json:"state, omitempty"`
	StateName      string `json:"state_name, omitempty"`
	Country        string `json:"country, omitempty"`
	CountryIso3166 string `json:"country_iso3166, omitempty"`
	Zip            string `json:"zip, omitempty"`
	Magic          string `json:"magic, omitempty"`
	Wmo            string `json:"wmo, omitempty"`
	Latitude       string `json:"latitude, omitempty"`
	Longitude      string `json:"longitude, omitempty"`
	Elevation      string `json:"elevation, omitempty"`
}

type ForecastResponseType struct {
	Response ResponseType `json:"response"`
	Forecast struct {
		SimpleForecast SimpleForecastType `json:"simpleforecast"`
	} `json:"forecast"`
}

type SimpleForecastType struct {
	ForecastDay []ForecastDayType `json:"forecastday"`
}

type ForecastDayType struct {
	Date       ForecastDateType `json:"date"`
	Conditions string           `json:"conditions"`
}

type ForecastDateType struct {
	Epoch string `json:"epoch"`
}

type WundergroundProvider struct {
	ApiKey            string
	reverseWeatherMap map[string]WeatherCondition
}

var WeatherMap = map[WeatherCondition][]string{
	CLEAR:           {"clear"},
	RAIN:            {"drizzle", "rain", "rain mist", "rain showers", "freezing rain"},
	PARTLY_CLOUDY:   {"overcast", "partly cloudy", "scattered clouds"},
	MOSTLY_CLOUDY:   {"mostly cloudy"},
	THUNDER:         {"thunderstorm"},
	THUNDER_RAIN:    {"thunderstorms and rain"},
	TORNADO:         {"thunderstorms and ice pellets", "thunderstorms with hail", "thunderstorms with small hail"},
	SNOW:            {"snow", "snow grains", "low drifting snow", "blowing snow", "snow showers", "snow blowing snow mist", "thunderstorms and snow", "sleet"},
	WINDY:           {"squalls", "funnel cloud"},
	FOGGY:           {"mist", "fog", "fog patches", "freezing fog", "patches of fog", "shallow fog", "partial fog"},
	HAZE:            {"haze"},
	UNKNOWN_WEATHER: {"ice crystals", "ice pellets", "hail", "smoke", "volcanic ash", "widespread dust", "sand", "spray", "dust whirls", "sandstorm", "low drifting widespread dust", "low drifting sand", "blowing widespread dust", "blowing sand", "ice pellet showers", "hail showers", "small hail showers", "freezing drizzle", "small hail", "unknown precipitation", "unknown"},
}

func NewWundergroundProvider(apiKey string) WeatherProvider {
	wunderground := &WundergroundProvider{ApiKey: apiKey}
	wunderground.reverseWeatherMap = buildReverseWeatherMap()

	return wunderground
}

func buildReverseWeatherMap() map[string]WeatherCondition {
	reverseMap := make(map[string]WeatherCondition)

	for wc, ww := range WeatherMap {
		for _, w := range ww {
			reverseMap[w] = wc
		}
	}

	return reverseMap
}

func (w *WundergroundProvider) getWeather(city, country string) (WeatherCondition, error) {
	log.Printf("query weather condition for %s/%s", city, country)
	requestURL := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%s/%s.json", w.ApiKey, country, city)

	log.Printf("request to %s", requestURL)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Get(requestURL)
	if err != nil {
		log.Printf("failed to send Http request, %v", err)
		return UNKNOWN_WEATHER, err
	}

	defer resp.Body.Close()

	conditionResp := &ConditionResponseType{}

	err = json.NewDecoder(resp.Body).Decode(conditionResp)
	if err != nil {
		log.Printf("could not parse response json, %v", err)
		return UNKNOWN_WEATHER, err
	}

	log.Printf("got response from weather api")
	log.Printf("%v", conditionResp)

	if conditionResp.Response.Error.Type != "" {
		log.Printf("server says we send incorrect request: %s, %s",
			conditionResp.Response.Error.Type, conditionResp.Response.Error.Description)
		return UNKNOWN_WEATHER, errors.New(conditionResp.Response.Error.Description)
	}

	return w.convertWeather(conditionResp.CurrentObservation.Weather), nil
}

func (w *WundergroundProvider) get3DaysForcast(city, country string) ([]WeatherCondition, error) {
	log.Printf("query 3 days weather forecast for %s/%s", city, country)
	requestURL := fmt.Sprintf("http://api.wunderground.com/api/%s/forecast/q/%s/%s.json", w.ApiKey, country, city)

	log.Printf("request to %s", requestURL)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Get(requestURL)
	if err != nil {
		log.Printf("failed to send Http request, %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	forecastResponse := &ForecastResponseType{}

	err = json.NewDecoder(resp.Body).Decode(forecastResponse)
	if err != nil {
		log.Printf("could not parse response json, %v", err)
		return nil, err
	}

	log.Printf("got forecast from weather api")
	log.Printf("%v", forecastResponse)

	if forecastResponse.Response.Error.Type != "" {
		log.Printf("server says we send incorrect request: %s, %s",
			forecastResponse.Response.Error.Type, forecastResponse.Response.Error.Description)
		return nil, errors.New(forecastResponse.Response.Error.Description)
	}

	if len(forecastResponse.Forecast.SimpleForecast.ForecastDay) <= 3 {
		log.Printf("server returns less than 3 days forecast, something might be wrong, but anyway we can continue")
	}

	weatherForecasts := make([]WeatherCondition, 0)

	for i, forecast := range forecastResponse.Forecast.SimpleForecast.ForecastDay {
		if i == 0 {
			continue
		}
		weatherForecasts = append(weatherForecasts, w.convertWeather(forecast.Conditions))
	}

	log.Printf("here is next 3 days forecast: %v", weatherForecasts)
	return weatherForecasts, nil
}

func (w *WundergroundProvider) convertWeather(weatherPhrase string) WeatherCondition {
	log.Printf("trying to convert weather phrase \"%s\"", weatherPhrase)
	weather := strings.ToLower(strings.TrimSpace(weatherPhrase))

	weather = strings.TrimPrefix(weather, "light ")
	weather = strings.TrimPrefix(weather, "heavy ")
	weather = strings.TrimPrefix(weather, "chance of ")
	weather = strings.TrimPrefix(weather, "a ")

	result, ok := w.reverseWeatherMap[weather]
	if !ok {
		log.Printf("could not find \"%s\" in reverse map", weather)
		result = UNKNOWN_WEATHER
	}

	log.Printf("convert to %s", result)

	return result
}
