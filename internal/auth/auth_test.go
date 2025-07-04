package auth

import (
	"duna/internal/database"
	"duna/internal/hash"
	"duna/internal/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testingUserUUID     = "4348b9dd-a8e9-448e-8b22-b985d54003b0"
	testingUserUsername = "user-123"
	testingUserEmail    = "testuser@gmail.com"
	testingUserPassword = "testing123"
)

func getTestUser(hash hash.HashStrategy) (models.User, error) {
	email, err := models.NewEmail(testingUserEmail)
	if err != nil {
		return models.User{}, fmt.Errorf("error creating email: %s", err.Error())
	}

	password, err := models.NewPassword(testingUserPassword, true, hash)
	if err != nil {
		return models.User{}, fmt.Errorf("error creating testing password: %s",
			err.Error())
	}

	return models.NewUser(
		testingUserUUID,
		testingUserUsername,
		email,
		password,
	), nil
}

func TestAuthenticate_Success(t *testing.T) {
	mockHash := models.HashStrategyMock{
		FuncCompare: func(enconded, str string) bool {
			assert.Equal(t, enconded, testingUserPassword)
			return true
		},
	}

	testUser, err := getTestUser(mockHash)
	if err != nil {
		t.Errorf("error creating testing password: %s", err.Error())
		return
	}

	mockDB := &database.MockDatabase{
		FuncGetUserByUsername: func(username string, hash hash.HashStrategy) (models.User, error) {
			assert.Equal(t, testingUserUsername, username)
			return testUser, nil
		},
	}

	mockStore := &MockSessionStore{
		FuncSet: func(key string, value SessionUserUUIDAndCsrftTokenPair) error {
			return nil
		},
	}

	auth := New(mockDB, mockStore, mockHash)
	sessionToken, csrfToken, err := auth.Authenticate(
		testingUserUsername, testingUserPassword)

	assert.NoError(t, err)
	assert.NotEmpty(t, sessionToken)
	assert.NotEmpty(t, csrfToken)
	assert.NotEqual(t, sessionToken, csrfToken)
}

func TestAuthenticate_PasswordsDontMatch(t *testing.T) {
	mockHash := models.HashStrategyMock{
		FuncCompare: func(enconded, str string) bool {
			return false
		},
	}

	testUser, err := getTestUser(mockHash)
	if err != nil {
		t.Errorf("error creating testing password: %s", err.Error())
		return
	}

	mockDB := &database.MockDatabase{
		FuncGetUserByUsername: func(username string, hash hash.HashStrategy) (models.User, error) {
			assert.Equal(t, testingUserUsername, username)
			return testUser, nil
		},
	}

	auth := New(mockDB, nil, mockHash)
	sessionToken, csrfToken, err := auth.Authenticate(
		testingUserUsername, testingUserPassword)

	assert.Equal(t, errors.New("passwords don´t match"), err)
	assert.Equal(t, sessionToken, "")
	assert.Equal(t, csrfToken, "")
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	mockDB := &database.MockDatabase{
		FuncGetUserByUsername: func(username string, hash hash.HashStrategy) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}

	mockStore := &MockSessionStore{}
	mockHash := &models.HashStrategyMock{}

	auth := New(mockDB, mockStore, mockHash)

	_, _, err := auth.Authenticate("nonexistent", "password123")

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
}

func TestAuthenticate_generateTokenShouldNotGenerateSameTokenMoreThanOnce(t *testing.T) {
	auth := sessionAuthenticator{}

	t1 := auth.generateToken()
	t2 := auth.generateToken()

	assert.NotEqual(t, t1, t2)
}

func TestGetUserUUID_Success(t *testing.T) {
	testUser, err := getTestUser(nil)
	if err != nil {
		t.Errorf("error creating testing password: %s", err.Error())
		return
	}

	testSessionToken := "random session token"
	testCsrftToken := "random csrft token"

	storeGetWasCalledTimes := 0
	storeGetCalledWith := ""

	mockStore := &MockSessionStore{
		FuncGet: func(key string) (SessionUserUUIDAndCsrftTokenPair, error) {
			storeGetCalledWith = key
			storeGetWasCalledTimes++

			return SessionUserUUIDAndCsrftTokenPair{
				CsrftToken: testCsrftToken,
				UserUUID:   testUser.UUID,
			}, nil
		},
	}

	auth := New(nil, mockStore, nil)
	uuid, err := auth.GetUserUUID(testSessionToken, testCsrftToken)

	assert.NoError(t, err)
	assert.Equal(t, uuid, testUser.UUID)
	assert.Equal(t, storeGetCalledWith, testSessionToken)
	assert.Equal(t, storeGetWasCalledTimes, 1)
}

func TestGetUserUUID_TokenNotFound(t *testing.T) {
	expectedError := errors.New("not found")

	storeGetWasCalledTimes := 0
	mockStore := &MockSessionStore{
		FuncGet: func(key string) (SessionUserUUIDAndCsrftTokenPair, error) {
			storeGetWasCalledTimes++
			return SessionUserUUIDAndCsrftTokenPair{}, expectedError
		},
	}

	auth := New(nil, mockStore, nil)
	uuid, err := auth.GetUserUUID("test token", "test token")

	assert.Equal(t, err, expectedError)
	assert.Equal(t, uuid, "")
	assert.Equal(t, storeGetWasCalledTimes, 1)
}

func TestGetUserUUID_CSFRTNotMatch(t *testing.T) {
	testUser, err := getTestUser(nil)
	if err != nil {
		t.Errorf("error creating testing password: %s", err.Error())
		return
	}

	storeGetWasCalledTimes := 0
	mockStore := &MockSessionStore{
		FuncGet: func(key string) (SessionUserUUIDAndCsrftTokenPair, error) {
			storeGetWasCalledTimes++
			return SessionUserUUIDAndCsrftTokenPair{
				CsrftToken: "invalid token",
				UserUUID:   testUser.UUID,
			}, nil
		},
	}

	auth := New(nil, mockStore, nil)
	uuid, err := auth.GetUserUUID("test token", "test token")

	assert.Equal(t, errors.New("wrong creadentials"), err)
	assert.Equal(t, uuid, "")
	assert.Equal(t, storeGetWasCalledTimes, 1)
}

func TestLogout_Success(t *testing.T) {
	testToken := "Test token"

	storeRemoveWasCalledTimes := 0
	storeRemoveCalledWith := ""
	mockStore := &MockSessionStore{
		FuncRemove: func(key string) error {
			storeRemoveCalledWith = key
			storeRemoveWasCalledTimes++

			return nil
		},
	}

	auth := New(nil, mockStore, nil)
	err := auth.Logout(testToken)

	assert.NoError(t, err)
	assert.Equal(t, storeRemoveCalledWith, testToken)
	assert.Equal(t, storeRemoveWasCalledTimes, 1)
}

func TestLogout_Fail(t *testing.T) {
	testToken := "Test token"

	storeRemoveWasCalledTimes := 0
	storeRemoveCalledWith := ""

	expectedError := errors.New("error")
	mockStore := &MockSessionStore{
		FuncRemove: func(key string) error {
			storeRemoveCalledWith = key
			storeRemoveWasCalledTimes++

			return expectedError
		},
	}

	auth := New(nil, mockStore, nil)
	err := auth.Logout(testToken)

	assert.Equal(t, err, expectedError)
	assert.Equal(t, storeRemoveCalledWith, testToken)
	assert.Equal(t, storeRemoveWasCalledTimes, 1)
}
