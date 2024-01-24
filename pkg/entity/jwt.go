package entity

import "github.com/google/uuid"

type JWTClaims struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Subject  uuid.UUID `json:"subject"`
	Issuer   string    `json:"iss,omitempty"`
}
