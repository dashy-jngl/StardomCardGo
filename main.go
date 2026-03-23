package main

import (
	"stardomcard/funcs"
)

func main() {
	funcs.ParseCard("https://wwr-stardom.com/event/20260325/")
	links, _ := funcs.GetCardLinksTwoMonths()
	for _, link := range links {
		println(link)
	}
}