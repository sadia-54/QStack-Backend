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