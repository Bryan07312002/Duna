package auth

import (
	"crypto/rand"
	"duna/internal/database"
	"duna/internal/hash"
	"encoding/base64"
	"errors"
)

type SessionAuthenticator interface {
	Authenticate(username, password string) (string, string, error)
	GetUserUUID(sessionToken, csrftToken string) (string, error)
	Logout(sessionToken string) error
}

type SessionUserUUIDAndCsrftTokenPair struct {
	CsrftToken string
	UserUUID   string
}

type SessionStore interface {
	Set(key string, value SessionUserUUIDAndCsrftTokenPair) error
	Get(key string) (SessionUserUUIDAndCsrftTokenPair, error)
	Remove(key string) error
}

type sessionAuthenticator struct {
	db    database.Database
	store SessionStore
	hash  hash.HashStrategy
}

func New(db database.Database, store SessionStore, hash hash.HashStrategy) SessionAuthenticator {
	return &sessionAuthenticator{
		db:    db,
		store: store,
		hash:  hash,
	}
}

func (s *sessionAuthenticator) Authenticate(
	username, password string) (string, string, error) {
	user, err := s.db.GetUserByUsername(username, s.hash)
	if err != nil {
		return "", "", err
	}

	if !user.Password().Compare(password) {
		return "", "", errors.New("passwords don´t match")
	}

	session, csrft := s.generateToken(), s.generateToken()
	s.store.Set(session, SessionUserUUIDAndCsrftTokenPair{
		CsrftToken: csrft,
		UserUUID:   user.UUID,
	})

	return session, csrft, nil
}

func (s *sessionAuthenticator) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *sessionAuthenticator) GetUserUUID(
	sessionToken, csrftToken string) (string, error) {
	pair, err := s.store.Get(sessionToken)
	if err != nil {
		return "", err
	}

	if pair.CsrftToken != csrftToken {
		return "", errors.New("wrong creadentials")
	}

	return pair.UserUUID, nil
}

func (s *sessionAuthenticator) Logout(sessionToken string) error {
	return s.store.Remove(sessionToken)
}
