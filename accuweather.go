package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	//TokyoLocationKey hard coded location key for Tokyo
	TokyoLocationKey = "226396"
	//ConditionAPIURL api endpoint for current condition
	ConditionAPIURL = "http://dataservice.accuweather.com/currentconditions/v1"
	//Forcast5DaysURL api endpoint for 5 days forcast
	Forcast5DaysURL = "http://dataservice.accuweather.com/forecasts/v1/daily/5day"
)

//icon code map to weather condition
//see: https://developer.accuweather.com/weather-icons
var accuIconMap = map[int]WeatherCondition{
	1:  CLEAR,         //Sunny
	2:  PARTLY_CLOUDY, //Mostly Sunny
	3:  PARTLY_CLOUDY, //Partly Sunny
	4:  PARTLY_CLOUDY, //Intermittent Clouds
	5:  HAZE,          //Hazy Sunshine
	6:  MOSTLY_CLOUDY, //Mostly Cloudy
	7:  MOSTLY_CLOUDY, //Cloudy
	8:  MOSTLY_CLOUDY, //Dreary (Overcast)
	11: FOGGY,         //Fog
	12: RAIN,          //Showers
	13: RAIN,          //Mostly Cloudy w/ Showers
	14: RAIN,          //Partly Sunny w/ Showers
	15: RAIN,          //T-Storms
	16: RAIN,          //Mostly Cloudy w/ T-Storms
	17: RAIN,          //Partly Sunny w/ T-Storms
	18: RAIN,          //Rain
	19: WINDY,         //Flurries
	20: MOSTLY_CLOUDY, //Mostly Cloudy w/ Flurries
	21: MOSTLY_CLOUDY, //Partly Sunny w/ Flurries
	22: SNOW,          //Snow
	23: MOSTLY_CLOUDY, //Mostly Cloudy w/ Snow
	24: SNOW,          //Ice
	25: RAIN,          //Sleet
	26: RAIN,          //Freezing Rain
	29: SNOW,          //Rain and Snow
	30: CLEAR,         //Hot
	31: SNOW,          //Cold
	32: WINDY,         //Windy
	33: CLEAR,         //Clear
	34: CLEAR,         //Mostly Clear
	35: PARTLY_CLOUDY, //Partly Cloudy
	36: PARTLY_CLOUDY, //Intermittent Clouds
	37: CLEAR,         //Hazy Moonlight
	38: MOSTLY_CLOUDY, //Mostly Cloudy
	39: RAIN,          //Partly Cloudy w/ Showers
	40: RAIN,          //Mostly Cloudy w/ Showers
	41: RAIN,          //Partly Cloudy w/ T-Storms
	42: RAIN,          //Mostly Cloudy w/ T-Storms
	43: RAIN,          //Mostly Cloudy w/ Flurries
	44: SNOW,          //Mostly Cloudy w/ Snow
}

//CurrentConditionResponse response struct for current condition api
type CurrentConditionResponse struct {
	LocalObservationDateTime string
	WeatherIcon              int
	WeatherText              string
	HasPrecipitation         bool
	PrecipitationType        string
	Temperature              struct {
		Metric TemperatureValueType
	}
}

//TemperatureValueType temperature struct, only include C
type TemperatureValueType struct {
	Value    float64
	Unit     string
	UnitType int
}

//Condition short condition of a day
type Condition struct {
	Icon       int
	IconPhrase string
}

//DailyForcastResponse api response for daily forcast
type DailyForcastResponse struct {
	DailyForecasts []DailyForcastType
}

//DailyForcastType contains one day's condition
type DailyForcastType struct {
	Date        string
	Temperature struct {
		Minimum TemperatureValueType
		Maximum TemperatureValueType
	}
	Day   Condition
	Night Condition
}

//AccuWeatherProvider accuweather provider
type AccuWeatherProvider struct {
	apiKey string
	client *http.Client
}

func (a *AccuWeatherProvider) getWeather(city, country string) (WeatherCondition, error) {
	location := getLocationKey(city, country)
	url := fmt.Sprintf("%s/%s", ConditionAPIURL, location)

	currentCondition := []CurrentConditionResponse{}

	err := a.callAPI(url, &currentCondition)
	if err != nil {
		return UNKNOWN_WEATHER, err
	}

	if len(currentCondition) < 1 {
		return UNKNOWN_WEATHER, fmt.Errorf("failed to get current condition")
	}

	log.Printf("%v", currentCondition)

	return convertIconToConditoin(currentCondition[0].WeatherIcon), nil
}

func (a *AccuWeatherProvider) get3DaysForcast(city, country string) ([]WeatherCondition, error) {
	location := getLocationKey(city, country)
	url := fmt.Sprintf("%s/%s", Forcast5DaysURL, location)

	forcast := DailyForcastResponse{}

	err := a.callAPI(url, &forcast)

	if err != nil {
		return nil, err
	}

	conditions := []WeatherCondition{}

	for i := 1; i < 4; i++ {
		c := UNKNOWN_WEATHER
		if len(forcast.DailyForecasts) > i {
			dayIcon := forcast.DailyForecasts[i].Day.Icon
			c = convertIconToConditoin(dayIcon)
		}

		conditions = append(conditions, c)
	}

	return conditions, nil
}

func (a *AccuWeatherProvider) callAPI(url string, v interface{}) error {
	url = fmt.Sprintf("%s?apikey=%s", url, a.apiKey)
	resp, err := a.client.Get(url)

	if err != nil {
		log.Printf("failed to send Http request, %v", err)
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		log.Printf("could not parse response, %v", err)
		return err
	}

	return nil
}

func convertIconToConditoin(icon int) WeatherCondition {
	if w, ok := accuIconMap[icon]; ok {
		return w
	}
	return UNKNOWN_WEATHER
}

//TODO: now always return tokyo's location key
func getLocationKey(city, country string) string {
	return TokyoLocationKey
}

//NewAccuWeatherProvider create a accuweather provider
// which use the accuweather api
func NewAccuWeatherProvider(apiKey string) WeatherProvider {
	return &AccuWeatherProvider{apiKey, &http.Client{}}
}
