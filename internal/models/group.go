package models

type Group struct {
	Id    uint   `gorm:"primarykey;autoIncrement"`
	Name  string `gorm:"notnull"`
	Songs []Song `json:"-"`
}

func (Group) TableName() string {
	return "groups"
}
