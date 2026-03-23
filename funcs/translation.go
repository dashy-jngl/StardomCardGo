package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func TranslateGoogle(text string) (string, error) {
	u := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=ja&tl=en&dt=t&q=%s", url.QueryEscape(text))

	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	
	outer, ok := result[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	inner, ok := outer[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	translated, ok := inner[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return translated, nil
}