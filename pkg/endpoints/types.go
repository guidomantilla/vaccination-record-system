package endpoints

import (
	"github.com/gin-gonic/gin"
)

type AuthEndpoint interface {
	Login(ctx *gin.Context)
	Signup(ctx *gin.Context)
	Authorize(ctx *gin.Context)
}

type DrugsEndpoint interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Find(ctx *gin.Context)
}

type VaccinationsEndpoint interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Find(ctx *gin.Context)
}
