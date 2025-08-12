package model

type SettingKey string

type Setting struct {
	ID    uint       `gorm:"primaryKey" json:"id"`
	Key   SettingKey `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value string     `gorm:"type:text;not null" json:"value"`
}
