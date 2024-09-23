//go:build tools

package tooling

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "go.uber.org/mock/mockgen"
)
