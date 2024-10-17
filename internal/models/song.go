package models

type Song struct {
	Id          uint   `gorm:"primarykey;autoIncrement"`
	Band        string `gorm:"-"`
	Song        string `gorm:"index:song_name_index;notnull"`
	ReleaseDate string `gorm:"column:release_date"`
	Text        string `gorm:"column:text"`
	Link        string `gorm:"index:link_index;column:link"`
	GroupId     uint   `gorm:"foreignKey:group_id"`
	Group       Group  `json:"-"`
}

func (Song) TableName() string {
	return "songs"
}
