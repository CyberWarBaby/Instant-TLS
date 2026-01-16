package models

import (
	"time"

	"github.com/google/uuid"
)

type Plan string

const (
	PlanFree Plan = "free"
	PlanPro  Plan = "pro"
	PlanTeam Plan = "team"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Plan         Plan      `db:"plan" json:"plan"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Token struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	Name       string     `db:"name" json:"name"`
	Prefix     string     `db:"prefix" json:"prefix"`
	TokenHash  string     `db:"token_hash" json:"-"`
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

type Machine struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	Hostname   string    `db:"hostname" json:"hostname"`
	OS         string    `db:"os" json:"os"`
	Arch       string    `db:"arch" json:"arch"`
	LastSeenAt time.Time `db:"last_seen_at" json:"last_seen_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// API response types
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Plan      Plan      `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
}

type LicenseResponse struct {
	Plan   Plan           `json:"plan"`
	Limits map[string]int `json:"limits"`
	User   UserResponse   `json:"user"`
}

type TokenResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type TokenCreateResponse struct {
	Token string        `json:"token"`
	Data  TokenResponse `json:"data"`
}
