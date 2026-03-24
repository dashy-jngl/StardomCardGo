package funcs

var TestLinks = []string{"https://wwr-stardom.com/event/20260325/"}

func PrintLinks() {
	links, _ := GetCardLinksTwoMonths()
	for _, link := range links {
		println(link)
	}
}

func TestOverrides() {
	matchCard, _ := ParseCard(TestLinks[0])
	overrides, _ := FetchNameOverrides()
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

func TestTranslation() {
	original := "スターダム女子プロレス - 6人タッグマッチ"
	translated, err := TranslateGoogle(original)
	if err != nil {
		println("Translation error:", err.Error())
	} else {
		println("Original:", original)
		println("Translated:", translated)
	}
}