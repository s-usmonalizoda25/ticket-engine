package domain

import "time"

type User struct {
    ID        int64      `json:"id"`
    Email     string     `json:"email"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}