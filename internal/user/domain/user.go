package domain

type User struct {
	UserId     uint   `gorm:"primaryKey"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Name       string `gorm:"not null"`
	ProfileURL string
}

type Viewer struct {
	User
}

type Administrator struct {
	User
}
