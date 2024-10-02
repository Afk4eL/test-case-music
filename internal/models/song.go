package models

type Song struct {
	Id          uint   `gorm:"primarykey"`
	Group       string `gorm:"notnull;column:band"`
	Song        string `gorm:"notnull"`
	ReleaseDate string `gorm:"column:release_date"`
	Text        string `gorm:"column:text"`
	Link        string `gorm:"column:link"`
}

func (Song) TableName() string {
	return "songs"
}
