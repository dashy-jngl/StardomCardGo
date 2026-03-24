package funcs

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

var vsIcons = []string{" vs "," ⚔️ ", " ⚔ ", " 🆚 ", }
var borderChars = [][]string{{"┌","┬", "┐", "└","┴", "┘", "─", "│"}, {"╔","╦", "╗", "╚","╩", "╝", "═", "║"}, {"┏", "┳", "┓", "┗","┻", "┛", "━", "┃"}, {"╭","─", "╮", "╰","─","╯", "─", "│"}, {"╒", "╤", "╕", "╘","╧", "╛", "═", "│"}}

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Default width if unable to get terminal size
	}
	return w
}

func centerLine(line string) string {
	w:= termWidth()
	lineW := runewidth.StringWidth(line)
	if lineW >= w {
		return line // No centering if line is wider than terminal
	}
	pad := (w - lineW) / 2
	return strings.Repeat(" ", pad) + line
}

func padCenter(text string, width int) string {
	len := runewidth.StringWidth(text)
	if len >= width {
		return text // No padding if text is wider than specified width
	}
	space := (width - len)
	left := space / 2
	right := space - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func PrintHeader(card *MatchCard) {
	fmt.Println()
	var parts []string
	if card.Time != "" {
		parts = append(parts, card.Time)
	}
	if card.Date != "" {
		parts = append(parts, card.Date)
	}
	if len(parts) > 0 {
		fmt.Println(centerLine(strings.Join(parts, " ~ ")))
	}
	fmt.Println(centerLine(fmt.Sprintf("『%s』", card.Title)))
	fmt.Println()
}

func maxRows(teams [][]string) int {
	rows := 0
	for _, t := range teams {
		if len(t) > rows {
			rows = len(t)
		}
	}
	return rows
}
func colWidth(teams [][]string,) map[int]int {
	widths := make(map[int]int)
	for i, team := range teams {
		for _, name := range team {
			w := runewidth.StringWidth(name)
			if w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}

func topParts(lenTeams int, widths map[int]int, style int, vs, matchType string) {
	borderstyle := borderChars[style]
	var topParts []string
	for i := 0; i < lenTeams; i++ {
		topParts = append(topParts, strings.Repeat(borderstyle[6], widths[i]))
		if i < lenTeams-1 {
			topParts = append(topParts, strings.Repeat(borderstyle[6], runewidth.StringWidth(vs)))
		}
	}
	fmt.Println(centerLine(matchType))
	fmt.Println(centerLine(borderstyle[0] + strings.Join(topParts, borderstyle[1]) + borderstyle[2]))
}

func midRow(teams [][]string, rows int, widths map[int]int, vsW int, borderstyle []string, vsStyle int){
	midRow := rows / 2
	for r := 0; r < rows; r++ {
		var rowParts []string
		for i, team := range teams {
			cell := ""
			if r < len(team) {
				cell = team[r]
			}
			rowParts = append(rowParts, padCenter(cell, widths[i]))
			if i < len(teams)-1 {
				sep := ""
				if r == midRow {
					sep = vsIcons[vsStyle]
				}
				rowParts = append(rowParts,padCenter(sep, vsW))
			}
		}
		fmt.Println(centerLine(borderstyle[7] + strings.Join(rowParts, borderstyle[7]) + borderstyle[7]))
	}
}

func bottomParts(lenTeams int, widths map[int]int, style int, vs string) {
	borderstyle := borderChars[style]
	var bottomParts []string
	for i := 0; i < lenTeams; i++ {
		bottomParts = append(bottomParts, strings.Repeat(borderstyle[6], widths[i]))
		if i < lenTeams-1 {
			bottomParts = append(bottomParts, strings.Repeat(borderstyle[6], runewidth.StringWidth(vs)))
		}
	}
	fmt.Println(centerLine(borderstyle[3] + strings.Join(bottomParts, borderstyle[4]) + borderstyle[5]))
}

func PrintMatchTable(matchType string, teams [][]string, style,vsStyle int,) {

	for i, team := range teams {
		for j, name := range team {
			teams[i][j] = " " + name + " "
		}
	}
	n:= len(teams)
	if n == 0 {
		return
	}

	//max rows
	rows := maxRows(teams)
	if rows == 0 {
		return
	}

	//col widths
	widths := colWidth(teams)
	//vs width
	vsW := runewidth.StringWidth(vsIcons[vsStyle])

	//top border
	topParts(n, widths, style, vsIcons[vsStyle], matchType)
	//rows
	midRow(teams, rows, widths, vsW, borderChars[style], vsStyle)
	//bottom border
	bottomParts(n, widths, style, vsIcons[vsStyle])
}

func PrintVerboseCard(card *MatchCard, style, vsStyle int) {
	PrintHeader(card)
	for i, m:= range card.Matches {
		header:= fmt.Sprintf("Match %d: %s", i+1, m.MatchType)
		PrintMatchTable(header, m.Teams, style, vsStyle)
	}
}

func PrintCompact(card *MatchCard, vsStyle int) {
	var title []string
	if card.Time != "" {
		title = append(title, card.Time)
	}
	if card.Date != "" {
		title = append(title, card.Date)
	}
	if card.Title != "" {
		title = append(title, fmt.Sprintf("『%s』", card.Title))
	}
	if len(title) > 0 {
		fmt.Println(strings.Join(title, " ~ "))
		fmt.Println()
	}
	
	for i, m:= range card.Matches {
		var teamStrs []string
		for _, t := range m.Teams {
			teamStrs = append(teamStrs, strings.Join(t, ", "))
		}
		fmt.Printf("%d. %s: %s\n", i+1, m.MatchType, strings.Join(teamStrs, fmt.Sprintf(" %s ", vsIcons[vsStyle])))
	}

}
