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
    Matches []Match
}

func ParseCard(url string) (MatchCard, error) {
    resp, err := http.Get(url)
    if err != nil {
        return MatchCard{}, err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return MatchCard{}, err
    }

    title := doc.Find("h1.match_head_title").First().Text()
    title = strings.TrimSpace(title)

    fmt.Println("Title:", title)
    return MatchCard{Title: title}, nil
}