package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func TranslateBatch(texts []string) (map[string]string, error) {
	joined := strings.Join(texts, "\n")
	u := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=ja&tl=en&dt=t&q=%s", url.QueryEscape(joined))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	outer, ok := result[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected Response Format")
	}

	var translated string
	for _, seg := range outer {
		inner, ok := seg.([]interface{})
		if !ok || len(inner) == 0 {
			continue
		}
		part, ok := inner[0].(string)
		if ok {
			translated += part
		}
	}
	parts := strings.Split(translated, "\n")
	mapping := make(map[string]string)
	for i, orig := range texts {
		if i < len(parts) {
			mapping[orig] = strings.TrimSpace(parts[i])
		}
	}
	return mapping, nil
}

func TranslateCard(card *MatchCard, overrides map[string]string) {

	//collect unique strings
	seen := make(map[string]bool)
	var unique []string

	for _, s := range []string{card.Title} {
		if s != "" && !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}
	for _, m := range card.Matches {
		if !seen[m.MatchType] {
			seen[m.MatchType] = true
			unique = append(unique, m.MatchType)
		}
		for _, team := range m.Teams {
			for _, name := range team {
				if !seen[name] {
					seen[name] = true
					unique = append(unique, name)
				}
			}
		}
	}
	// build mapping - overrides first, then translate the rest
	mapping := make(map[string]string)
	var toTranslate []string
	for _, s := range unique {
		if en, ok := overrides[s]; ok {
			mapping[s] = en
		} else {
			toTranslate = append(toTranslate, s)
		}
	}
	if len(toTranslate) > 0 {
		translated, err := TranslateBatch(toTranslate)
		if err == nil {
			for k, v := range translated {
				mapping[k] = v
			}
		}
	}

	// apply mapping back to card
	if t, ok := mapping[card.Title]; ok {
		card.Title = t
	}
	for i := range card.Matches {
		if t, ok := mapping[card.Matches[i].MatchType]; ok {
			card.Matches[i].MatchType = t
		}

		for j := range card.Matches[i].Teams {
			for k := range card.Matches[i].Teams[j] {
				if t, ok := mapping[card.Matches[i].Teams[j][k]]; ok {
					card.Matches[i].Teams[j][k] = t
				}
			}
		}
	}
}
