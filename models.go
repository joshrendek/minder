package main

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
	ProjectID   uint
	Completed   bool
}

type Project struct {
	gorm.Model
	Name  string
	Tasks []Task
}
