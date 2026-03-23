package main

import (
	"fmt"
	"stardomcard/funcs"
)

var testLinks = []string{"https://wwr-stardom.com/event/20260325/"}

func printLinks() {
	links, _ := funcs.GetCardLinksTwoMonths()
	for _, link := range links {
		println(link)
	}
}

func testOverrides() {
	matchCard, _ := funcs.ParseCard(testLinks[0])
	overrides, _ := funcs.FetchNameOverrides()
	for _, match := range matchCard.Matches {
		for _, team := range match.Teams {
			for _, wrestler := range team {
				println("Original name:", wrestler)
				if override, exists := overrides[wrestler]; exists {
					println("Overridden name:", override)
				}
			}
		}
	}
}

func testTranslation() {
	original := "スターダム女子プロレス - 6人タッグマッチ"
	translated, err := funcs.TranslateGoogle(original)
	if err != nil {
		println("Translation error:", err.Error())
	} else {
		println("Original:", original)
		println("Translated:", translated)
	}
}

func main() {
	matchCard, _ := funcs.ParseCard(testLinks[0])
	printLinks()
	// testOverrides()
	// testTranslation()
	overrides, _ := funcs.FetchNameOverrides()
	funcs.TranslateCard(&matchCard, overrides)
	fmt.Println(matchCard)
}
