// Package user include User struct and methods associated with it.
package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"casinogame/internal/ui"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Balance int    `json:"balance"`
}

func CreateUser() (User, error) {
	name := ui.ReadLine("Введите имя: ")
	age, err := ui.ReadInt("Введите возраст: ")
	if err != nil {
		return User{}, errors.New("некорректный возраст")
	}

	if age < 18 {
		return User{}, errors.New("вам меньше 18, вход запрещен")
	}

	users, err := ReadUsersJSON()
	if err != nil {
		return User{}, err
	}

	for _, u := range users {
		if u.Name == name {
			return User{}, errors.New("пользователь с таким именем уже существует")
		}
	}

	return User{
		Name:    name,
		Age:     age,
		Balance: 10000,
	}, nil
}

func UserMarshal(users []User) ([]byte, error) {
	data, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UserUnmarshal(data []byte) (User, error) {
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func WriteToJSON(data []byte) error {
	if err := os.WriteFile("users.json", data, 0o644); err != nil {
		return err
	}
	return nil
}

func ReadUsersJSON() ([]User, error) {
	file, err := os.Open("users.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var user []User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func AddNewUserJSON(user User) error {
	users, err := ReadUsersJSON()
	if err != nil {
		return err
	}

	users = append(users, user)
	musers, err := UserMarshal(users)
	if err != nil {
		return err
	}
	if err = WriteToJSON(musers); err != nil {
		return err
	}

	return nil
}

func ChoiceUser(users []User) User {
	for {
		ui.ClearConsole()
		fmt.Println("Выберите пользователя:")

		for _, user := range users {
			fmt.Printf("Имя: %s | Баланс: %d \n", user.Name, user.Balance)
		}

		fmt.Println("Новый пользователь (нп)")
		ans := ui.ReadLine("-> ")

		for _, user := range users {
			if ans == user.Name {
				return user
			}
		}

		if strings.ToLower(ans) == "нп" {
			newUser, err := CreateUser()
			if err != nil {
				fmt.Println(err.Error())
				ui.ReadLine("Enter чтобы продолжить...")
				continue
			}
			if err := AddNewUserJSON(newUser); err != nil {
				fmt.Println(err.Error())
				ui.ReadLine("Enter чтобы продолжить...")
				continue
			}
			updatedUsers, err := ReadUsersJSON()
			if err != nil {
				fmt.Println(err.Error())
				ui.ReadLine("Enter чтобы продолжить...")
				continue
			}
			users = updatedUsers
			continue
		}

		ui.ReadLine("Пользователь не найден..(Enter)")
	}
}
