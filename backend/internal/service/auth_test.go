package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"creditanalysis/internal/model"
	"creditanalysis/internal/service"
)

// fakeUserRepo is an in-memory repository.UserRepository for testing the
// AuthService in isolation from the database.
type fakeUserRepo struct {
	user *model.User
	err  error
}

func (f *fakeUserRepo) FindByEmail(_ context.Context, _ string) (*model.User, error) {
	return f.user, f.err
}

func TestAuthService_Login(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	svc := service.NewAuthService(&fakeUserRepo{
		user: &model.User{ID: 1, Email: "admin@creditanalysis.com", PasswordHash: string(hash)},
	}, "test-secret")

	t.Run("valid credentials return a token", func(t *testing.T) {
		token, err := svc.Login(context.Background(), "admin@creditanalysis.com", "senha123")
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "admin@creditanalysis.com", "wrong")
		require.ErrorIs(t, err, service.ErrInvalidCredentials)
	})
}
