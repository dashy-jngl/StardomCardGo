package main

import (
	"flag"
	"fmt"
	"os"
	"stardomcard/funcs"
)

func main() {

	date := flag.String("d", "", "Date in (YYYYMMDD)")
	compact := flag.Bool("c", false, "Compact display")
	vsStyle := flag.Int("vs", 1, "VS Icon Style (1-3)")
	style := flag.Int("s", 3, "Table Style (1-4)")
	english := flag.Bool("e", false, "Translate to English")
	list := flag.Bool("l", false, "List cards for the next 2 months")
	n := flag.Int("n", 1, "nth card to display (default 1, ignored if -l is set)")
	flag.Parse()

	if *style < 1 || *style > 4 {
		fmt.Println("Invalid style. Please choose between 1 and 4.")
		os.Exit(1)
	}
	if *vsStyle < 1 || *vsStyle > 3 {
		fmt.Println("Invalid VS style. Please choose between 1 and 3.")
		os.Exit(1)
	}
	selectedStyle := *style - 1
	selectedVs := *vsStyle - 1

	if *list {
		links, err := funcs.GetCardLinksTwoMonths()
		if err != nil {
			fmt.Println("Error fetching card links:", err)
			os.Exit(1)
		}
		if len(links) == 0 {
			fmt.Println("No cards found.")
			os.Exit(0)
		}
		fmt.Printf("Upcoming Cards:\n")
		for i, link := range links {
			fmt.Printf("%d. %s\n", i+1, link)
		}
		return
	}

	var url string
	if *date != "" {
		url = fmt.Sprintf("https://wwr-stardom.com/schedule/%s/", *date)
	} else {
		links, err := funcs.GetCardLinksTwoMonths()
		if err != nil {
			fmt.Println("Error fetching card links:", err)
			os.Exit(1)
		}
		if len(links) == 0 {
			fmt.Println("No cards found.")
			os.Exit(0)
		}
		if *n < 1 || *n > len(links) {
			fmt.Printf("Invalid card number. Please choose between 1 and %d.\n", len(links))
			os.Exit(1)
		}
		url = links[*n-1]
	}

	card, err := funcs.ParseCard(url)
	if err != nil {
		fmt.Println("Error parsing card:", err)
		os.Exit(1)
	}

	if *english {
		overrides, err := funcs.FetchNameOverrides()
		if err != nil {
			fmt.Println("Error fetching name overrides:", err)
			os.Exit(1)
		}
		funcs.TranslateCard(&card, overrides)
	}

	if *compact {
		funcs.PrintCompact(&card, *&selectedVs)
	} else {
		funcs.PrintVerboseCard(&card, *&selectedStyle, *&selectedVs)
	}
	// matchCard, _ := funcs.ParseCard(funcs.TestLinks[0])
	// funcs.PrintLinks()
	// overrides, _ := funcs.FetchNameOverrides()
	// engMatchCard := matchCard
	// funcs.TranslateCard(&engMatchCard, overrides)
	// funcs.PrintVerboseCard(&matchCard, 1, 1)
	// funcs.PrintVerboseCard(&engMatchCard, 1,1)
	// funcs.PrintCompact(&matchCard, 1)
	// funcs.PrintCompact(&engMatchCard, 1)
}
