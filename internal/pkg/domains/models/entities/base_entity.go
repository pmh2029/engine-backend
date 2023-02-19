package entity

import (
	"gorm.io/gorm"
)

type BaseEntity struct {
	gorm.Model
}
