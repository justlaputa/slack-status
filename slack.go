package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Emoji string

const (
	EMOJI_SUNNY              Emoji = ":sunny:"
	EMOJI_MOSTLY_SUNNY             = ":mostly_sunny:"
	EMOJI_PARTLY_SUNNY             = ":partly_sunny:"
	EMOJI_BARELY_SUNNY             = ":barely_sunny:"
	EMOJI_PARTLY_SUNNY_RAIN        = ":partly_sunny_rain:"
	EMOJI_CLOUD                    = ":cloud:"
	EMOJI_RAIN_CLOUD               = ":rain_cloud:"
	EMOJI_THUNDER_CLOUD_RAIN       = ":thunder_cloud_and_rain:"
	EMOJI_LIGHTING                 = ":lightning:"
	EMOJI_SNOW_FLAKE               = ":snowflake:"
	EMOJI_SNOW                     = ":snow_cloud:"
	EMOJI_SNOWMAN                  = ":snowman:"
	EMOJI_TORNADO                  = ":tornado:"
	EMOJI_FOG                      = ":fog:"
	EMOJI_UMBRELLA_RAIN            = ":umbrella_with_rain_drops:"
	EMOJI_WIND                     = ":wind_blowing_face:"
	EMOJI_MASK                     = ":mask:"
	EMOJI_QUESTION                 = ":question:"
)

var (
	WeatherEmojiMap = map[WeatherCondition]Emoji{
		CLEAR:           EMOJI_SUNNY,
		RAIN:            EMOJI_RAIN_CLOUD,
		PARTLY_CLOUDY:   EMOJI_MOSTLY_SUNNY,
		MOSTLY_CLOUDY:   EMOJI_BARELY_SUNNY,
		THUNDER:         EMOJI_LIGHTING,
		THUNDER_RAIN:    EMOJI_THUNDER_CLOUD_RAIN,
		TORNADO:         EMOJI_TORNADO,
		SNOW:            EMOJI_SNOW_FLAKE,
		WINDY:           EMOJI_WIND,
		FOGGY:           EMOJI_FOG,
		HAZE:            EMOJI_MASK,
		UNKNOWN_WEATHER: EMOJI_QUESTION,
	}
)

type SimpleResponseType struct {
	OK      bool   `json:"ok"`
	Warning string `json:"warning"`
	Error   string `json:"error"`
}

type SlackApi struct {
	Token string
}

func (s *SlackApi) setStatus(text string, emoji Emoji) bool {
	statusRequest, err := http.NewRequest("GET", "https://slack.com/api/users.profile.set", nil)
	if err != nil {
		log.Printf("failed to create request object, %v", err)
		return false
	}

	status := struct {
		StatusText  string `json:"status_text"`
		StatusEmoji Emoji  `json:"status_emoji"`
	}{text, emoji}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		log.Printf("unable to convert status object to json string, %v", err)
		return false
	}

	log.Printf("trying to set user status profile to: %s", statusJSON)

	q := statusRequest.URL.Query()
	q.Add("token", s.Token)
	q.Add("profile", string(statusJSON))
	statusRequest.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(statusRequest)
	if err != nil {
		log.Printf("failed to send request to slack api, %v", err)
		return false
	}

	defer resp.Body.Close()
	statusResponse := &SimpleResponseType{}
	err = json.NewDecoder(resp.Body).Decode(statusResponse)
	if err != nil {
		log.Printf("failed to parse response as json, %v", err)
		return false
	}

	log.Printf("got response from server: %v", statusResponse)

	if !statusResponse.OK {
		log.Printf("slack server says our api is not ok, error: %s", statusResponse.Error)
		return false
	} else if statusResponse.Warning != "" {
		log.Printf("slack server accepted our request, but send an warning that something might be wrong: %s", statusResponse.Warning)
	}

	return true
}
