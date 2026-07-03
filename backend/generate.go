// This file hosts repository-wide go:generate directives.
//
// Run `go generate ./...` from the backend module root to regenerate the
// OpenAPI docs under ./docs from the swag annotations on the handlers.
// It must run from the module root so swag scans internal/ for annotations.
package tools

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.4 init --generalInfo cmd/api/main.go --output docs --parseInternal --parseDependency --parseDepth 2
