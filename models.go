package main

import (
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	Name string
}
