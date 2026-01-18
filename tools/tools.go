//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/golangci/golines"
	_ "go.uber.org/mock/mockgen"
)
