package serve

import (
	"context"
	"fmt"
	"github.com/guidomantilla/vaccination-record-system/pkg/endpoints"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	feather_commons_config "github.com/guidomantilla/go-feather-commons/pkg/config"
	feather_commons_environment "github.com/guidomantilla/go-feather-commons/pkg/environment"
	feather_commons_log "github.com/guidomantilla/go-feather-commons/pkg/log"
	feather_security "github.com/guidomantilla/go-feather-security/pkg/security"
	feather_web_rest "github.com/guidomantilla/go-feather-web/pkg/rest"
	feather_web_server "github.com/guidomantilla/go-feather-web/pkg/server"
	"github.com/qmdx00/lifecycle"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/cobra"

	"github.com/guidomantilla/vaccination-record-system/pkg/config"
)

func ExecuteCmdFn(_ *cobra.Command, args []string) {

	ctx := context.Background()
	logger := feather_commons_log.Custom()
	appName, version := config.Application, config.Version

	app := lifecycle.NewApp(
		lifecycle.WithName(appName),
		lifecycle.WithVersion(version),
		lifecycle.WithSignal(syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL),
	)

	osArgs := os.Environ()
	environment := feather_commons_environment.NewDefaultEnvironment(feather_commons_environment.WithArrays(osArgs, args))

	var cfg config.Config
	if err := feather_commons_config.Process(ctx, environment, &cfg); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	//Database

	passwordEncoder := feather_security.NewBcryptPasswordEncoder()
	passwordGenerator := feather_security.NewDefaultPasswordGenerator()
	passwordManager := feather_security.NewDefaultPasswordManager(passwordEncoder, passwordGenerator)
	principalManager := feather_security.NewInMemoryPrincipalManager(passwordManager) //TODO make DB based
	tokenManager := feather_security.NewJwtTokenManager(feather_security.WithIssuer(appName),
		feather_security.WithSigningKey([]byte(*cfg.TokenSignatureKey)),
		feather_security.WithVerifyingKey([]byte(*cfg.TokenVerificationKey)))
	authenticationService := feather_security.NewDefaultAuthenticationService(passwordManager, principalManager, tokenManager)
	authenticationEndpoint := feather_security.NewDefaultAuthenticationEndpoint(authenticationService)

	authorizationService := feather_security.NewDefaultAuthorizationService(tokenManager, principalManager)
	authorizationFilter := feather_security.NewDefaultAuthorizationFilter(authorizationService)
	authPrincipalEndpoint := endpoints.NewDefaultAuthPrincipalEndpoint(principalManager)
	drugsEndpoint := endpoints.NewDefaultDrugsEndpoint()
	vaccinationsEndpoint := endpoints.NewDefaultVaccinationsEndpoint()
	//

	recoveryFilter := gin.Recovery()
	loggerFilter := sloggin.New(logger.RetrieveLogger().(*slog.Logger).WithGroup("http"))
	customFilter := func(ctx *gin.Context) {
		feather_security.AddApplicationToContext(ctx, appName)
		ctx.Next()
	}

	publicHandler := gin.New()
	publicHandler.Use(loggerFilter, recoveryFilter, customFilter)
	publicHandler.POST("/login", authenticationEndpoint.Authenticate)
	publicHandler.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	publicHandler.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, feather_web_rest.NotFoundException("resource not found"))
	})
	publicHandler.GET("/info", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"appName": appName})
	})

	privateHandler := publicHandler.Group("/api", authorizationFilter.Authorize)
	privateHandler.POST("/signup", authPrincipalEndpoint.Signup)

	privateHandler.POST("/drugs", drugsEndpoint.Create)
	privateHandler.PUT("/drugs/:id", drugsEndpoint.Update)
	privateHandler.DELETE("/drugs/:id", drugsEndpoint.Delete)
	privateHandler.GET("/drugs/", drugsEndpoint.Find)

	privateHandler.POST("/vaccinations", vaccinationsEndpoint.Create)
	privateHandler.PUT("/vaccinations/:id", vaccinationsEndpoint.Update)
	privateHandler.DELETE("/vaccinations/:id", vaccinationsEndpoint.Delete)
	privateHandler.GET("/vaccinations/", vaccinationsEndpoint.Find)

	//

	httpServer := &http.Server{
		Addr:              net.JoinHostPort(*cfg.Host, *cfg.HttpPort),
		Handler:           publicHandler,
		ReadHeaderTimeout: 60000,
	}
	app.Attach("HttpServer", feather_web_server.BuildHttpServer(httpServer))
	err := app.Run()
	if err != nil {
		feather_commons_log.Fatal(err.Error())
	}
}
