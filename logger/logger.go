// Package logger for logging games.
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	User   string
	Game   string
	Bid    int
	Profit int
	Time   time.Time
}

var logchan = make(chan Logger, 100)
var done = make(chan struct{})

func init() {
	go worker()
}

func Log(user, game string, bid, profit int) {
	logchan <- Logger{
		User:   user,
		Game:   game,
		Bid:    bid,
		Profit: profit,
		Time:   time.Now(),
	}
}

func worker() {
	file, err := os.OpenFile("gamelog.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Println("logger: не удалось открыть файл:", err)
		close(done)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for l := range logchan {
		if err := encoder.Encode(l); err != nil {
			fmt.Println("logger: записи в файл:", err)
		}
	}

	close(done)
}

func Close() {
	close(logchan)
	<-done
}
