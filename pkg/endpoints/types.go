package endpoints

import "github.com/gin-gonic/gin"

var (
	_ AuthPrincipalEndpoint = (*DefaultAuthPrincipalEndpoint)(nil)
)

type AuthPrincipalEndpoint interface {
	Signup(ctx *gin.Context)
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
