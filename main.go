package main

import (
	"stardomcard/funcs"
)

func main() {
	matchCard, _ := funcs.ParseCard("https://wwr-stardom.com/event/20260325/")
	links, _ := funcs.GetCardLinksTwoMonths()
	for _, link := range links {
		println(link)
	}
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