package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	c "casinogame/casino"
	g "casinogame/games"
	"casinogame/ui"
	u "casinogame/user"
)

func main() {
	users, err := u.ReadUsersJSON()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println("Ошибка чтения файла:", err.Error())
			return
		}

		newUser, err := u.CreateUser()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		data, err := u.UserMarshal([]u.User{newUser})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := u.WriteToJSON(data); err != nil {
			fmt.Println(err.Error())
			return
		}

		users = []u.User{newUser}
	}

	casino := c.InitCasino()
	user := u.ChoiceUser(users)

	for {
		drawWelcome(user)
		ans := ui.ReadLine("-> ")

		if strings.ToLower(ans) == "выход" {
			fmt.Println("До встречи.")
			return
		}

		if strings.ToLower(ans) == "сп" {
			user = u.ChoiceUser(users)
			continue
		}

		game := g.FindGame(ans)
		if game == nil {
			ui.ReadLine("Игра не найдена, попробуйте еще раз. (Enter)")
			continue
		}

		bid, err := ui.ReadInt("Ваша ставка -> ")
		if err != nil {
			fmt.Println("Некорректная ставка")
			ui.ReadLine("Нажми Enter чтобы продолжить...")
			continue
		}
		factor, err := ui.ReadFloat("Множитель -> ")
		if err != nil {
			fmt.Println("Некорректный множитель")
			ui.ReadLine("Нажми Enter чтобы продолжить...")
			continue
		}

		if err := game.Play(&casino, &user, bid, factor); err != nil {
			fmt.Println(err.Error())
		}
		ui.ReadLine("Нажмите, что бы продолжить...")
	}
}

func drawWelcome(user u.User) {
	ui.ClearConsole()

	ui.PrintArtAt(ui.SlotMachine, 2, 40)

	fmt.Printf("\n\n\n\nДобро пожаловать в казино, %s!\n", user.Name)
	fmt.Printf("Баланс: %d\n", user.Balance)

	fmt.Println("Выбери: ")
	for _, v := range g.Games {
		fmt.Printf("%s(%s) | ", v.Name, v.Command)
	}
	fmt.Println("Выход")
	fmt.Println("Смена пользователя(сп)")
}
