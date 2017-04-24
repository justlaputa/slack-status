package main

import "testing"

func TestSetStatus(t *testing.T) {
	t.Skip("skip test slack")
	api := &SlackApi{Token: ""}

	result := api.setStatus("Hello :earth_asia:", EMOJI_SUNNY)
	if !result {
		t.Fail()
	}
}
