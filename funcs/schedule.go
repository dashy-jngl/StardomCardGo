package funcs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Match struct {
	MatchType string
	Teams     [][]string
}
type MatchCard struct {
	Time, Date, Title string
	Matches           []Match
}

func ParseCard(url string) (MatchCard, error) {
	// Make an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return MatchCard{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return MatchCard{}, err
	}

	// Extract the title of the match card
	title := doc.Find("h1.match_head_title").First().Text()
	title = strings.TrimSpace(title)

	// Extract the date of the match card
	date := doc.Find("p.date").First().Text()
	date = strings.TrimSpace(date)

	// Extract the matches
	var matches []Match

	doc.Find("div.match_wrap").Each(func(i int, wrap *goquery.Selection) {
		matchType := strings.TrimSpace(wrap.Find("h2.sub_content_title1").Text())

		row := wrap.Find("div.match_block_row")
		if row.Length() > 0 {
			var left []string
			row.Find("div.leftside h3.name").Each(func(_ int, s *goquery.Selection) {
				left = append(left, strings.TrimSpace(s.Text()))
			})
			var right []string
			row.Find("div.rightside h3.name").Each(func(_ int, s *goquery.Selection) {
				right = append(right, strings.TrimSpace(s.Text()))
			})
			matches = append(matches, Match{MatchType: matchType, Teams: [][]string{left, right}})
		} else {
			// 3-col match(e.g. 3v3)
			var teams [][]string

			wrap.Find("div.match_block_column ul.match_block_3col").Each(func(_ int, ul *goquery.Selection) {
				var members []string
				ul.Find("h3.name").Each(func(_ int, s *goquery.Selection) {
					members = append(members, strings.TrimSpace(s.Text()))
				})
				teams = append(teams, members)
			})
			matches = append(matches, Match{MatchType: matchType, Teams: teams})

		}
	})

	// Debug output
	fmt.Println("Title:", title)
	fmt.Println("Date:", date)
	for _, match := range matches {
		fmt.Println("Match Type:", match.MatchType)
		for i, team := range match.Teams {
			fmt.Printf("  Team %d: %s\n", i+1, strings.Join(team, ", "))
		}
	}
	return MatchCard{Title: title, Date: date, Matches: matches}, nil
}
