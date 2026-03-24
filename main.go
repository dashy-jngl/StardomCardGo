package main

import (
	"stardomcard/funcs"
)

func main() {
	matchCard, _ := funcs.ParseCard(funcs.TestLinks[0])
	funcs.PrintLinks()
	overrides, _ := funcs.FetchNameOverrides()
	funcs.TranslateCard(&matchCard, overrides)
	funcs.PrintVerboseCard(&matchCard, 1, 1)
	funcs.PrintCompact(&matchCard, 1)
}
