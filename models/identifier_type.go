package models

type IdentifierType struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;size:100" json:"name"`
}
