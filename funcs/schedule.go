package funcs

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

type Match struct {
	MatchType string
	Teams     [][]string
}
type MatchCard struct {
	Time, Date, Title string
	Matches           []Match
}

func ParseCard(cardUrl string) (MatchCard, error) {

	//url for time:
	timeUrl := strings.Replace(cardUrl, "event", "schedule", 1)

	var cardDoc, timeDoc *goquery.Document
	g := new(errgroup.Group)
	//fetch event page
	g.Go(func() error {
		resp, err := http.Get(cardUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		cardDoc, err = goquery.NewDocumentFromReader(resp.Body)
		return err
	})
	//fetch schedule page
	g.Go(func() error {
		resp, err := http.Get(timeUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		timeDoc, err = goquery.NewDocumentFromReader(resp.Body)
		return err
	})

	if err := g.Wait(); err != nil {
		return MatchCard{}, err
	}

	// Extract the title of the match card
	title := cardDoc.Find("h1.match_head_title").First().Text()
	title = strings.TrimSpace(title)

	// Extract the date of the match card
	date := cardDoc.Find("p.date").First().Text()
	date = strings.TrimSpace(date)

	// Extract the Time
	time := timeDoc.Find("span.time").Eq(2).Text()
	time = strings.TrimSpace(time)

	// Extract the matches
	var matches []Match

	cardDoc.Find("div.match_wrap").Each(func(i int, wrap *goquery.Selection) {
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
	return MatchCard{Title: title, Date: date, Time: time, Matches: matches}, nil
}

func GetCardLinksTwoMonths() ([]string, error) {
	baseURL := "https://wwr-stardom.com/schedule/"

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	links := extractCardLinks(doc, baseURL)

	doc.Find("a.calendar_btn_next").First().Each(func(_ int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			ref, _ := url.Parse(href)
			base, _ := url.Parse(baseURL)
			nextURL := base.ResolveReference(ref)
			resp2, err := http.Get(nextURL.String())
			if err != nil {
				return
			}
			defer resp2.Body.Close()
			doc2, err := goquery.NewDocumentFromReader(resp2.Body)
			if err != nil {
				return
			}
			moreLinks := extractCardLinks(doc2, nextURL.String())
			links = append(links, moreLinks...)
		}
	})
	return links, nil

}

func extractCardLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	base, _ := url.Parse(baseURL)

	doc.Find("a.btn").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "対戦カード" {
			if href, exists := s.Attr("href"); exists {
				ref, err := url.Parse(href)
				if err == nil {
					fullURL := base.ResolveReference(ref)
					links = append(links, fullURL.String())
				}
			}
		}
	})
	return links
}
