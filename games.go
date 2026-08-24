package main

import "strings"

type Game struct {
	Name    string
	Command string
	Play    func(c *Casino, user *User, bid int, factor float64) error
}

var games = []Game{
	{
		Name:    "Депать",
		Command: "деп",
		Play:    (*Casino).Dep,
	},
	{
		Name:    "Красное & Чёрное",
		Command: "кч",
		Play:    (*Casino).BlackAndRed,
	},
}

func findGame(command string) *Game {
	for _, g := range games {
		if strings.ToLower(command) == g.Command {
			return &g
		}
	}
	return nil
}
