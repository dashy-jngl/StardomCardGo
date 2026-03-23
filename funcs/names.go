package funcs

import (
	"encoding/json"
	"net/http"
)

const OverridesURL = "https://raw.githubusercontent.com/dashy-jngl/joshi-data/refs/heads/main/name_overrides.json"

func FetchNameOverrides() (map[string]string, error) {
	resp, err := http.Get(OverridesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var overrides map[string]string
	err = json.NewDecoder(resp.Body).Decode(&overrides)
	return overrides, err
}