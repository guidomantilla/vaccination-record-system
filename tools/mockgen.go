package tools

//go:generate mockgen -package datasource -destination ../pkg/datasource/mocks.go -source ../pkg/datasource/types.go
//go:generate mockgen -package endpoints -destination ../pkg/endpoints/mocks.go -source ../pkg/endpoints/types.go
//go:generate mockgen -package services -destination ../pkg/services/mocks.go -source ../pkg/services/types.go
