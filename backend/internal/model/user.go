package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
