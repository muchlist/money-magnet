package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Name        string
	Password    string
	Roles       []string
	PocketRoles []string
	Fcm         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

// DTO adalah refresentasi dari balikan database atau API sehingga apabila terjadi perubahan
// aplikasi CORE tidak harus dirubah namun DTO akan digunakan apabila ada perbedaan antara
// domain dan dto kedepannya
// type UserDTO struct {
// 	ID          uuid.UUID
// 	Email       string
// 	Name        string
// 	Password    string
// 	Roles       []string
// 	PocketRoles []string
// 	Fcm         string
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

type UserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResp struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Roles       []string  `json:"roles"`
	PocketRoles []string  `json:"pocket_roles"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"`
}
