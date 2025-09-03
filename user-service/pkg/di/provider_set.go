package di

import (
	"github.com/anrisys/quicket/user-service/internal"
	"github.com/anrisys/quicket/user-service/pkg/config"
	"github.com/anrisys/quicket/user-service/pkg/database"
	"github.com/anrisys/quicket/user-service/pkg/security"
	"github.com/anrisys/quicket/user-service/pkg/token"
	"github.com/google/wire"
)

var (
	ConfigSet = wire.NewSet(
		config.Load,
		database.ConnectMySQL,
		config.NewZerolog,
		security.NewAccountSecurity,
		token.NewTokenGenerator,
		wire.Bind(new(security.AccountSecurityInterface), new(*security.AccountSecurity)),
		wire.Bind(new(token.TokenGeneratorInterface), new(*token.TokenGenerator)),
	)
	UserAppProviderSet = wire.NewSet(
		ConfigSet,
		internal.NewUserRepository,
		internal.NewUserService,
		internal.NewUserHandler,
		wire.Bind(new(internal.UserRepositoryInterface), new(*internal.UserRepository)),
		wire.Bind(new(internal.UserServiceInterface), new(*internal.UserService)),
		wire.Struct(new(UserServiceApp), "*"),
	)
)