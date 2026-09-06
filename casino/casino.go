// Package casino include User struct and methods associated with it.
package casino

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"

	"casinogame/logger"
	"casinogame/ui"
	u "casinogame/user"
)

type Casino struct {
	balance int
}

func InitCasino() Casino {
	return Casino{
		balance: 1000000,
	}
}

func (c *Casino) Dep(user *u.User, bid int, factor float64, gameName string) error {
	if bid <= 0 {
		return errors.New("ставка должна быть положительной")
	}
	if factor <= 0 {
		return errors.New("множитель должен быть положительным")
	}
	if user.Balance < 0 {
		return errors.New("ваш баланс отрицательный")
	}
	if user.Balance < bid {
		return errors.New("баланс меньше чем ставка")
	}

	win := rand.Int64N(2) == 1

	if win {
		if err := c.userWin(user, bid, factor, gameName); err != nil {
			return err
		}
	} else {
		if err := c.userLoose(user, bid, gameName); err != nil {
			return err
		}
	}

	if c.balance < 0 {
		fmt.Println("У казино закончились деньги..попробуйте в другой раз.")
		return nil
	}

	if user.Balance <= 0 {
		fmt.Println("Баланс закончился, доигрались.")
		return nil
	}

	ans := ui.ReadLine("Додеп? (Да/Нет): ")
	if strings.ToLower(ans) == "да" {
		return c.Dep(user, bid, factor, gameName)
	}
	return nil
}

func (c *Casino) BlackAndRed(user *u.User, bid int, factor float64, gameName string) error {
	if bid <= 0 {
		return errors.New("ставка должна быть положительной")
	}
	if factor <= 0 {
		return errors.New("множитель должен быть положительным")
	}
	if user.Balance < 0 {
		return errors.New("ваш баланс отрицательный")
	}
	if user.Balance < int(math.Round(float64(bid)*factor)) {
		return errors.New("баланс меньше чем возможный проигрыш")
	}
	if c.balance < 0 {
		fmt.Println("У казино закончились деньги..попробуйте в другой раз.")
		return nil
	}

	var ans string
	rbnumber := rand.IntN(2)
	bones := []string{"красное", "чёрное"}
	ans = ui.ReadLine("Выберите: красное или чёрное ->  ")

	if strings.ToLower(ans) == bones[rbnumber] {
		if err := c.userWin(user, bid, factor, gameName); err != nil {
			return err
		}
	} else if strings.ToLower(ans) != "красное" && strings.ToLower(ans) != "чёрное" {
		return errors.New("нужно ввести 'красное' или 'чёрное'")
	} else {
		fmt.Println("Выпало:", bones[rbnumber])
		if err := c.userLoose(user, bid, gameName); err != nil {
			return err
		}
	}
	return nil
}

func (c *Casino) userWin(user *u.User, bid int, factor float64, gameName string) error {
	amount := int(math.Round(float64(bid) * factor))
	c.balance -= amount
	user.Balance += amount

	usersJSON, err := u.ReadUsersJSON()
	if err != nil {
		return err
	}

	for i := range usersJSON {
		if usersJSON[i].Name == user.Name {
			usersJSON[i].Balance = user.Balance
		}
	}

	musersJSON, err := u.UserMarshal(usersJSON)
	if err != nil {
		return err
	}

	if err = u.WriteToJSON(musersJSON); err != nil {
		return err
	}

	logger.Log(user.Name, gameName, bid, amount)

	fmt.Printf("Вы выйграли! Выйгрышь: %d. Баланс: %d\n", amount, user.Balance)
	return nil
}

func (c *Casino) userLoose(user *u.User, bid int, gameName string) error {
	c.balance += bid
	user.Balance -= bid

	usersJSON, err := u.ReadUsersJSON()
	if err != nil {
		return err
	}

	for i := range usersJSON {
		if usersJSON[i].Name == user.Name {
			usersJSON[i].Balance = user.Balance
		}
	}

	musersJSON, err := u.UserMarshal(usersJSON)
	if err != nil {
		return err
	}

	if err = u.WriteToJSON(musersJSON); err != nil {
		return err
	}

	logger.Log(user.Name, gameName, bid, bid)
	fmt.Printf("Вы проиграли! Проигрышь: %d. Баланс: %d\n", bid, user.Balance)
	return nil
}
