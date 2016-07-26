package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	readline "gopkg.in/readline.v1"
)

var (
	db      *gorm.DB
	err     error
	cyan    = color.New(color.FgCyan).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	rl      *readline.Instance
	context Context
)

func main() {
	db, err = gorm.Open("sqlite3", "minder.db")
	defer db.Close()

	db.AutoMigrate(&Project{}, &Task{})

	completer := readline.NewPrefixCompleter(
		readline.PcItem("say",
			readline.PcItem("hello"),
			readline.PcItem("bye"),
		),
		readline.PcItem("help"),
	)

	rl, err = readline.NewEx(&readline.Config{
		Prompt:       fmt.Sprintf("%s ~> ", green("main")),
		AutoComplete: completer},
	)

	commander := NewCommander(completer)
	context.Commands = commander

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		context.Commands.Run(line)
	}
}
