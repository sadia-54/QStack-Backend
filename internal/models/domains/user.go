package domains

import "time"

type User struct {
    ID                        int64
    Email                     string
    PasswordHash              string
    Username                  string
    Bio                       *string
    EmailNotificationsEnabled bool
    EmailVerified             bool
    CreatedAt                 time.Time
    UpdatedAt                 time.Time
}

func NewUser(email, username, passwordHash string) *User {
    return &User{
        Email: email,
        Username: username,
        PasswordHash: passwordHash,
        EmailVerified: false,
        EmailNotificationsEnabled: false,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}