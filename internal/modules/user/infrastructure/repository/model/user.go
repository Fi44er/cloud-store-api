package user_model

type User struct {
	ID       string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string `gorm:"not null;unique;type:varchar(255)"`
	Email    string `gorm:"not null;unique;type:varchar(255)"`

	QuotaMax   int64  `gorm:"not null;default:0"`
	AvatarPath string `gorm:"type:varchar(255)"`
	BannerPath string `gorm:"type:varchar(255)"`
}
