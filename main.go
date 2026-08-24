package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	esc = "\033["
	dol = "$ "
)

var (
	scanner     = bufio.NewScanner(os.Stdin)
	slotMachine = strings.Split(`⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣀⣤⣤⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣶⣤⣤⣀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢙⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⣛⡋⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⣿⡇⠀⠀⣤⣄⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⠀⢸⣿⣿⡇⢸⣿⣿⡇⢸⣿⣿⡇⠀⣿⢸⠀⠀⣿⠛⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⠀⢸⣿⣿⡇⢸⣿⣿⡇⢸⣿⣿⡇⠀⣿⢸⠀⠀⣿⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⠀⢸⣿⣿⡇⢸⣿⣿⡇⢸⣿⣿⡇⠀⣿⢸⠀⣾⡇⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣿⡆⠀⣿⡿⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠁⠀⠙⠃⠀⠀⠀
⠀⠀⠀⠀⢀⣴⣿⠟⠛⠛⢻⡿⠛⠛⠛⢻⣿⣿⡟⠋⠉⠉⠛⢿⣦⡀⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣤⣤⣤⣤⣾⣧⣤⣤⣤⣿⣿⣿⣷⣦⣤⣤⣶⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠀⠀⠀⠀`, "\n")
)

func readLine(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func readInt(prompt string) (int, error) {
	s := readLine(prompt)
	return strconv.Atoi(s)
}

func readFloat(prompt string) (float64, error) {
	s := readLine(prompt)
	return strconv.ParseFloat(s, 64)
}

func main() {
	users, err := ReadUsersJson()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println("Ошибка чтения файла:", err.Error())
			return
		}

		newUser, err := NewUser()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		data, err := UserMarshal([]User{newUser})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := WriteToJson(data); err != nil {
			fmt.Println(err.Error())
			return
		}

		users = []User{newUser}
	}

	casino := InitCasino()
	user := ChoiceUser(users)

	for {
		drawWelcome(user)
		ans := readLine("-> ")

		if strings.ToLower(ans) == "выход" {
			fmt.Println("До встречи.")
			return
		}

		if strings.ToLower(ans) == "сп" {
			user = ChoiceUser(users)
			continue
		}

		game := findGame(ans)
		if game == nil {
			readLine("Игра не найдена, попробуйте еще раз. (Enter)")
			continue
		}

		bid, err := readInt("Ваша ставка -> ")
		if err != nil {
			fmt.Println("Некорректная ставка")
			readLine("Нажми Enter чтобы продолжить...")
			continue
		}
		factor, err := readFloat("Множитель -> ")
		if err != nil {
			fmt.Println("Некорректный множитель")
			readLine("Нажми Enter чтобы продолжить...")
			continue
		}

		if err := game.Play(&casino, &user, bid, factor); err != nil {
			fmt.Println(err.Error())
		}
		readLine("Нажмите, что бы продолжить...")
	}
}

func clearConsole() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func moveCursor(row, col int) {
	fmt.Printf("%s%d;%dH", esc, row, col)
}

func saveCursor() {
	fmt.Print(esc + "s")
}

func restoreCursor() {
	fmt.Print(esc + "u")
}

func printArtAt(art []string, row, col int) {
	saveCursor()
	for i, line := range art {
		moveCursor(row+i, col)
		fmt.Print(line)
	}
	restoreCursor()
}

func drawWelcome(user User) {
	clearConsole()

	// печатаем арт справа, начиная с 1-й строки, с 50-й колонки
	printArtAt(slotMachine, 2, 40)

	// обычный текст печатается как ни в чём не бывало, слева, построчно
	fmt.Printf("\n\n\n\nДобро пожаловать в казино, %s!\n", user.Name)
	fmt.Printf("Баланс: %d\n", user.Balance)

	fmt.Println("Выбери: ")
	for _, v := range games {
		fmt.Printf("%s(%s) | ", v.Name, v.Command)
	}
	fmt.Println("Выход")
	fmt.Println("Смена пользователя(сп)")
}
