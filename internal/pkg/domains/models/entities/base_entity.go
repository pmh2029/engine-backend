package entities

import (
	"time"
)

type BaseEntity struct {
	ID        uint       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}
