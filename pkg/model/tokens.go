package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"log"
	"time"

	"github.com/n31t/go-project/pkg/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type (
	Token struct {
		PlainText string    `json:"token"`
		Hash      []byte    `json:"-"`
		UserId    int64     `json:"-"`
		Scope     string    `json:"-"`
		ExpiresAt time.Time `json:"-"`
	}

	TokenModel struct {
		DB       *sql.DB
		InfoLog  *log.Logger
		ErrorLog *log.Logger
	}
)

func (t *TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

//	func (t *TokenModel) Insert(token *Token) error {
//		query := `
//		INSERT INTO tokens (hash, user_id, scope, expires_at)
//		VALUES ($1, $2, $3, $4)
//		RETURNING id, hash, user_id, scope, expires_at`
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//		_, err := t.DB.ExecContext(ctx, query, token.Hash, token.UserId, token.Scope, token.ExpiresAt)
//		return err
//	}
func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
		`

	args := []interface{}{token.Hash, token.UserId, token.ExpiresAt, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId:    userId,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (t *TokenModel) DeleteAllForUser(scope string, userId int64) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, userId, scope)
	return err

}

func ValidateTokenPlainText(v *validator.Validator, tokenText string) {
	v.Check(tokenText != "", "token", "must be provided")
	v.Check(len(tokenText) == 26, "token", "must be 26 bytes long")
}
