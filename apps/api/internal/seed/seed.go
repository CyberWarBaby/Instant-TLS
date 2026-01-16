package seed

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func SeedDemoUser(db *sqlx.DB) error {
	// Check if demo user exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", "demo@instanttls.dev")
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Demo user already exists
	}

	// Create demo user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("demo1234"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, plan)
		VALUES ($1, $2, $3, $4)
	`, userID, "demo@instanttls.dev", string(passwordHash), "pro")

	return err
}
