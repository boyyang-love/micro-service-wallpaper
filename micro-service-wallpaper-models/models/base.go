package models

import "time"

type Base struct {
	Id        string     `json:"id" form:"id" gorm:"primaryKey"`
	Created   int64      `gorm:"autoCreateTime:milli"`
	Updated   int64      `gorm:"autoUpdateTime:milli"`
	DeletedAt *time.Time `json:"deleted_at"`
}
