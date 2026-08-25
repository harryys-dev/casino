package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

type Casino struct {
	balance int
}

func InitCasino() Casino {
	return Casino{
		balance: 1000000,
	}
}

func (c *Casino) Dep(user *User, bid int, factor float64) error {
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

	win := rand.Int64N(2) == 1

	if win {
		c.userWin(user, bid, factor)
	} else {
		if err := c.userLoose(user, bid, factor); err != nil {
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

	ans := readLine("Додеп? (Да/Нет): ")
	if strings.ToLower(ans) == "да" {
		return c.Dep(user, bid, factor)
	}
	return nil
}

func (c *Casino) BlackAndRed(user *User, bid int, factor float64) error {
	var ans string
	rbnumber := rand.IntN(2)
	bones := []string{"красное", "чёрное"}

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

	ans = readLine("Выберите: красное или чёрное ->  ")

	if strings.ToLower(ans) == bones[rbnumber] {
		c.userWin(user, bid, factor)
	} else if strings.ToLower(ans) != "красное" && strings.ToLower(ans) != "чёрное" {
		return errors.New("Нужно ввести 'красное' или 'чёрное'")
	} else {
		fmt.Println("Выпало:", bones[rbnumber])
		if err := c.userLoose(user, bid, factor); err != nil {
			return err
		}
	}
	return nil
}

func (c *Casino) userWin(user *User, bid int, factor float64) {
	amount := int(math.Round(float64(bid) * factor))
	c.balance -= amount
	user.Balance += amount

	usersJson, err := ReadUsersJson()
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := range usersJson {
		if usersJson[i].Name == user.Name {
			usersJson[i].Balance = user.Balance
		}
	}

	musersJson, err := UserMarshal(usersJson)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err = WriteToJson(musersJson); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Вы выйграли! Выйгрышь: %d. Баланс: %d\n", amount, user.Balance)
}

func (c *Casino) userLoose(user *User, bid int, factor float64) error {
	amount := int(math.Round(float64(bid) * factor))

	usersJson, err := ReadUsersJson()
	if err != nil {
		return err
	}

	for i := range usersJson {
		if usersJson[i].Name == user.Name {
			usersJson[i].Balance = user.Balance
		}
	}

	musersJson, err := UserMarshal(usersJson)
	if err != nil {
		return err
	}

	if err = WriteToJson(musersJson); err != nil {
		return err
	} else {
		c.balance += amount
		user.Balance -= amount
		fmt.Printf("Вы проиграли! Проигрышь: %d. Баланс: %d\n", amount, user.Balance)
		return nil
	}
}
