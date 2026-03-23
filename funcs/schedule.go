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

    //
    date := doc.Find("p.date").First().Text()
    date = strings.TrimSpace(date)

    fmt.Println("Title:", title)
    fmt.Println("Date:", date)
    return MatchCard{Title: title, Date: date}, nil
}