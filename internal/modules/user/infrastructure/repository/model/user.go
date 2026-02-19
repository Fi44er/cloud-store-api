package user_model

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username     string `gorm:"not null;unique;type:varchar(255)"`
	Email        string `gorm:"not null;unique;type:varchar(255)"`
	PasswordHash string `gorm:"not null;type:varchar(255)"`
	AvatarPath   string `gorm:"type:varchar(255)"`
	BannerPath   string `gorm:"type:varchar(255)"`
}
