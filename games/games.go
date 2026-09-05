// Package games include User struct and methods associated with it.
package games

import (
	"strings"

	"casinogame/casino"
	"casinogame/user"
)

type Game struct {
	Name    string
	Command string
	Play    func(c *casino.Casino, user *user.User, bid int, factor float64) error
}

var Games = []Game{
	{
		Name:    "Депать",
		Command: "деп",
		Play:    (*casino.Casino).Dep,
	},
	{
		Name:    "Красное & Чёрное",
		Command: "кч",
		Play:    (*casino.Casino).BlackAndRed,
	},
}

func FindGame(command string) *Game {
	for i, g := range Games {
		if strings.ToLower(command) == g.Command {
			return &Games[i]
		}
	}
	return nil
}
