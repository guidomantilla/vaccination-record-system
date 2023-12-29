package endpoints

import (
	"github.com/gin-gonic/gin"
)

type DefaultVaccinationsEndpoint struct {
}

func NewDefaultVaccinationsEndpoint() *DefaultVaccinationsEndpoint {
	return &DefaultVaccinationsEndpoint{}
}

func (endpoint *DefaultVaccinationsEndpoint) Create(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultVaccinationsEndpoint) Update(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultVaccinationsEndpoint) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (endpoint *DefaultVaccinationsEndpoint) Find(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
