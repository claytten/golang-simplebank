package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

/** start testing normally **/
func CreateRandomUser(t *testing.T) db.Users {
	HashedPassword, err := util.HashingPassword(util.RandomString(10))
	require.NoError(t, err)
	require.NotEmpty(t, HashedPassword)

	arg := db.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: HashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestGetUserUsingEmail(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserUsingEmail(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUserAllFields(t *testing.T) {
	user1 := CreateRandomUser(t)
	password := util.RandomString(10)

	newFullname := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword, err := util.HashingPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, newPassword)

	user2, err := testQueries.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: user1.Username,
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, newFullname, user2.FullName)
	require.Equal(t, newEmail, user2.Email)
	require.Equal(t, newPassword, user2.HashedPassword)
	require.NotZero(t, user2.PasswordChangedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUserOnlyFullname(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newFullname := util.RandomOwner()

	user2, err := testQueries.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, newFullname, user2.FullName)
	require.WithinDuration(t, oldUser.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newEmail := util.RandomEmail()

	user2, err := testQueries.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, newEmail, user2.Email)
	require.WithinDuration(t, oldUser.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := CreateRandomUser(t)
	password := util.RandomString(10)
	newPassword, err := util.HashingPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, newPassword)

	user2, err := testQueries.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, newPassword, user2.HashedPassword)
	require.NoError(t, util.ComparePassword(user2.HashedPassword, password))
	require.NotZero(t, user2.PasswordChangedAt)
	require.WithinDuration(t, oldUser.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.UpdatedAt, user2.UpdatedAt, time.Second)
}

/** end testing normally **/
