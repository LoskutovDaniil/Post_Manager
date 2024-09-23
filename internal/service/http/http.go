package http

import (
	"fmt"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/service/http/graphql"

	gqlgenHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
)

type service struct {
	repo adaptor.Repository
}

func Run(address string, repo adaptor.Repository) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	hand := service{repo}
	r.POST("/api/v1/graphql", hand.graphqlHandler())

	err := r.Run(address)
	if err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

func (srv service) graphqlHandler() gin.HandlerFunc {
	config := graphql.Config{
		Resolvers: graphql.NewResolver(srv.repo),
	}
	schema := graphql.NewExecutableSchema(config)
	handler := gqlgenHandler.NewDefaultServer(schema)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
