package domains

import "time"

type User struct {
    ID        int64  `gorm:"primaryKey;autoIncrement"`
    Email     string `gorm:"type:varchar(255);not null;uniqueIndex"`
    PasswordHash string `gorm:"type:text;not null"`
    Username  string `gorm:"type:varchar(50);not null"`
    Bio       *string `gorm:"type:text"`
    EmailNotificationsEnabled bool `gorm:"type:boolean;not null;default:false"`
    EmailVerified bool `gorm:"type:boolean;not null;default:false"`
    CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
    UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (User) TableName() string {
    return "users"
}

func NewUser(email, username, passwordHash string) *User {
	now := time.Now()

	return &User{
		Email:                     email,
		Username:                  username,
		PasswordHash:              passwordHash,
		EmailVerified:             false,
		EmailNotificationsEnabled: false,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
}