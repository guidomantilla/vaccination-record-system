package tools

//go:generate mockgen -package rpc -destination ../pkg/endpoint/rpc/mocks.go -source ../pkg/endpoint/rpc/api_grpc.pb.go
//go:generate mockgen -package rest -destination ../pkg/endpoint/rest/mocks.go -source ../pkg/endpoint/rest/types.go
//go:generate mockgen -package repositories -destination ../pkg/repositories/mocks.go -source ../pkg/repositories/types.go
