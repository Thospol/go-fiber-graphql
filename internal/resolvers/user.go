package resolvers

import (
	"context"
	"fiber-graphql/internal/core/sql"
	"fiber-graphql/internal/models"
	"fiber-graphql/internal/repositories"

	"fiber-graphql/internal/core/config"

	"github.com/sirupsen/logrus"
)

type user struct {
	repository repositories.UserRepository
	config     *config.Configs
	result     *config.ReturnResult
}

// NewService new service
func NewService() *user {
	return &user{
		repository: repositories.NewUserRepository(),
		config:     config.CF,
		result:     config.RR,
	}
}

// Health health check
func (resolver *user) User(ctx context.Context, args struct{ ID int32 }) (*models.User, error) {
	object := &models.User{}
	err := resolver.repository.FindOneObjectByID(sql.Database, uint64(args.ID), object)
	if err != nil {
		logrus.Errorf("find user error: %v", err)
		return nil, resolver.result.Internal.DatabaseNotFound
	}

	return object, nil
}
