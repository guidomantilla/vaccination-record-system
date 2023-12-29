package endpoints

import (
	"github.com/gin-gonic/gin"
)

type DefaultDrugsEndpoint struct {
}

func NewDefaultDrugsEndpoint() *DefaultDrugsEndpoint {
	return &DefaultDrugsEndpoint{}
}

func (endpoint *DefaultDrugsEndpoint) Create(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultDrugsEndpoint) Update(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultDrugsEndpoint) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultDrugsEndpoint) Find(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
