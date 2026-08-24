package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Balance int    `json:"balance"`
}

func NewUser() (User, error) {
	name := readLine("Введите имя: ")
	age, err := readInt("Введите возраст: ")
	if err != nil {
		fmt.Println("Некорректный возраст")
		return User{}, errors.New("некорректный возраст")
	}

	if age < 18 {
		return User{}, errors.New("вам меньше 18, вход запрещен")
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

func WriteToJson(data []byte) error {
	if err := os.WriteFile("users.json", data, 0o644); err != nil {
		return err
	}
	return nil
}

func ReadUsersJson() ([]User, error) {
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

func ChoiceUser(users []User) User {
	for {
		clearConsole()
		fmt.Println("Выберите пользователя:")

		for _, user := range users {
			fmt.Printf("Имя: %s | Баланс: %d \n", user.Name, user.Balance)
		}

		fmt.Println("Новый пользователь (нп)")
		ans := readLine("-> ")

		for _, user := range users {
			if strings.EqualFold(ans, user.Name) {
				return user
			}
		}

		if strings.ToLower(ans) == "нп" {
			newUser, err := NewUser()
			if err != nil {
				fmt.Println(err.Error())
				readLine("Enter чтобы продолжить...")
				continue
			}
			if err := AddNewUserJson(newUser); err != nil {
				fmt.Println(err.Error())
				readLine("Enter чтобы продолжить...")
				continue
			}
			updatedUsers, err := ReadUsersJson()
			if err != nil {
				fmt.Println(err.Error())
				readLine("Enter чтобы продолжить...")
				continue
			}
			users = updatedUsers
			continue
		}

		readLine("Пользователь не найден..(Enter)")
	}
}

func AddNewUserJson(user User) error {
	users, err := ReadUsersJson()
	if err != nil {
		return err
	}

	users = append(users, user)
	musers, err := UserMarshal(users)
	if err != nil {
		return err
	}
	if err = WriteToJson(musers); err != nil {
		return err
	}

	return nil
}
