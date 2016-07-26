package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"text/tabwriter"

	"github.com/justincampbell/timeago"

	readline "gopkg.in/readline.v1"
)

type CommandHandler interface {
	Run(string)
	Register() []readline.PrefixCompleterInterface
	Match(string) bool
}

type Context struct {
	C        interface{}
	Hint     string
	Commands *Commander
}

type Commander struct {
	Completer *readline.PrefixCompleter
	Handlers  []CommandHandler
}

type Cd struct{}

func (c *Cd) Run(in string) {
	args := strings.Split(in, " ")
	if len(args) != 2 || args[1] == "" {
		println("missing project name, cd <name>")
		return
	}

	if args[1] == ".." || args[1] == "../" {
		context.C = nil
		context.Hint = "main"
		context.Commands = NewCommander(context.Commands.Completer)
		rl.SetPrompt(fmt.Sprintf("%s ~> ", green("main")))
		return
	}

	project := Project{}
	db.Where("name = ?", args[1]).First(&project)
	rl.SetPrompt(fmt.Sprintf("%s ~> ", green(project.Name)))
	context.C = project
	context.Hint = "project"
	context.Commands = NewProjectCommander(context.Commands.Completer)
}

func (c *Cd) Match(in string) bool {
	return strings.Contains(in, "cd")
}

func (c *Cd) Register() []readline.PrefixCompleterInterface {
	completers := []readline.PrefixCompleterInterface{}
	completers = append(completers, readline.PcItem("cd"))
	projects := []Project{}
	db.Find(&projects)
	for _, p := range projects {
		completers = append(completers, readline.PcItem(fmt.Sprintf("cd %s", p.Name)))
	}
	return completers
}

type CreateProject struct{}

type ListProjects struct{}

func (l *ListProjects) Run(in string) {
	//rl.SetPrompt(fmt.Sprintf("%s ~> ", green("ls")))
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	projects := []Project{}
	db.Find(&projects)
	fmt.Fprintln(w, green("Projects"), ": ", len(projects))
	for _, p := range projects {
		duration := time.Since(p.CreatedAt)
		fmt.Fprintln(w, "\t", green("-"), "\t", p.Name, "\t", yellow(timeago.FromDuration(duration)), "ago")
	}
	w.Flush()
}

func (l *ListProjects) Match(in string) bool {
	return in == "ls"
}

func (l *ListProjects) Register() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{readline.PcItem("ls")}
}

type ListTasks struct{}

func (l *ListTasks) Run(in string) {
	println("listing tasks...")
}

func (l *ListTasks) Register() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{readline.PcItem("ls")}
}

func (l *ListTasks) Match(in string) bool {
	return in == "ls"
}

func (c *CreateProject) Run(in string) {
	// check args, only 2
	args := strings.Split(in, " ")
	if len(args) != 2 || args[1] == "" {
		println("missing project name, create-project <name>")
		return
	}

	project := Project{Name: args[1]}
	db.Create(&project)
	println("project created")
}

func (c *CreateProject) Register() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{readline.PcItem("create-project")}
}

func (c *CreateProject) Match(in string) bool {
	return strings.Contains(in, "create-project")
}

func NewProjectCommander(completer *readline.PrefixCompleter) *Commander {
	println("new project commander")
	handlers := []CommandHandler{
		&ListTasks{},
		&Cd{},
	}

	for _, h := range handlers {
		completer.Children = append(completer.Children, h.Register()...)
	}
	return &Commander{Completer: completer, Handlers: handlers}
}

func NewCommander(completer *readline.PrefixCompleter) *Commander {
	handlers := []CommandHandler{
		&CreateProject{},
		&ListProjects{},
		&Cd{},
	}

	for _, h := range handlers {
		completer.Children = append(completer.Children, h.Register()...)
	}
	return &Commander{Completer: completer, Handlers: handlers}
}

func (c *Commander) Run(line string) {
	for _, h := range c.Handlers {
		if h.Match(line) {
			h.Run(line)
		}
	}
}
