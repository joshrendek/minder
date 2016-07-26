package main

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
	ProjectID   uint
}

type Project struct {
	gorm.Model
	Name  string
	Tasks []Task
}
