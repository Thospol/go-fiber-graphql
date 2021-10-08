package routes

import (
	ctx "context"
	"fiber-graphql/internal/core/config"
	"fiber-graphql/internal/core/context"
	"fiber-graphql/internal/core/utils"
	"fiber-graphql/internal/handlers/middlewares"
	"fiber-graphql/internal/resolvers"
	"fiber-graphql/internal/schemas"
	"fmt"

	"os"
	"os/signal"
	"time"

	fgg "github.com/cckwes/fiber-graphql-go"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/sirupsen/logrus"
)

const (
	// MaximumSize100MB body limit 100 mb.
	MaximumSize100MB = 1024 * 1024 * 100
	// MaximumSize1MB body limit 1 mb.
	MaximumSize1MB = 1024 * 1024 * 1
)

// New new router
func New() {
	app := fiber.New(
		fiber.Config{
			IdleTimeout:    5 * time.Second,
			BodyLimit:      MaximumSize100MB,
			ReadBufferSize: MaximumSize1MB,
		},
	)

	app.Use(
		compress.New(),
		requestid.New(),
		cors.New(),
		middlewares.WrapError(),
		middlewares.TransactionDatabase(func(c *fiber.Ctx) bool {
			return c.Method() == fiber.MethodGet
		}),
	)

	auth := jwtware.New(jwtware.Config{
		Claims:        &context.Claims{},
		SigningMethod: jwt.SigningMethodES256.Name,
		SigningKey:    utils.VerifyKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.
				Status(config.RR.Internal.Unauthorized.HTTPStatusCode()).
				JSON(config.RR.Internal.Unauthorized.WithLocale(c))
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
	})

	userSchema, err := schemas.ParseSchema("./internal/schemas/schema.gql", resolvers.NewService())
	if err != nil {
		panic(err)
	}
	userHandler := fgg.Handler{Schema: userSchema}
	app.Use(middlewares.AcceptLanguage())
	app.Use(middlewares.Logger())
	app.Post("/graphql", auth, userHandler.ServeHTTP)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
		defer cancel()

		logrus.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	logrus.Infof("Start server on port: %d ...", 8000)
	err = app.Listen(fmt.Sprintf(":%d", 8000))
	if err != nil {
		logrus.Panic(err)
	}
}
