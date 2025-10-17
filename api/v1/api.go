package v1

// Use github.com/deepmap/oapi-codegen/v2 version v2.2.0

//go:generate oapi-codegen -package v1 -generate types,chi-server,client -o api.gen.go api.swagger.yaml
