package serve

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	feather_commons_config "github.com/guidomantilla/go-feather-commons/pkg/config"
	feather_commons_environment "github.com/guidomantilla/go-feather-commons/pkg/environment"
	feather_commons_log "github.com/guidomantilla/go-feather-commons/pkg/log"
	feather_security "github.com/guidomantilla/go-feather-security/pkg/security"
	feather_web_rest "github.com/guidomantilla/go-feather-web/pkg/rest"
	feather_web_server "github.com/guidomantilla/go-feather-web/pkg/server"
	"github.com/qmdx00/lifecycle"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/guidomantilla/vaccination-record-system/pkg/config"
	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/endpoints"
	"github.com/guidomantilla/vaccination-record-system/pkg/services"
)

func ExecuteCmdFn(_ *cobra.Command, args []string) {

	var err error
	ctx := context.Background()
	logger := feather_commons_log.Custom()
	appName, version := config.Application, config.Version

	osArgs := os.Environ()
	environment := feather_commons_environment.NewDefaultEnvironment(feather_commons_environment.WithArrays(osArgs, args))

	var cfg config.Config
	if err = feather_commons_config.Process(ctx, environment, &cfg); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	url := strings.Replace(*cfg.DatasourceUrl, ":username", *cfg.DatasourceUsername, 1)
	url = strings.Replace(url, ":password", *cfg.DatasourcePassword, 1)
	url = strings.Replace(url, ":server", *cfg.DatasourceServer, 1)
	url = strings.Replace(url, ":service", *cfg.DatasourceService, 1)

	var db *gorm.DB
	if db, err = datasource.Open(logger, url); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	transactionHandler := datasource.NewTransactionHandler(db)
	passwordEncoder := feather_security.NewBcryptPasswordEncoder()
	passwordGenerator := feather_security.NewDefaultPasswordGenerator()
	passwordManager := feather_security.NewDefaultPasswordManager(passwordEncoder, passwordGenerator)

	var tokenTimeout time.Duration
	if tokenTimeout, err = time.ParseDuration(*cfg.TokenTimeout); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}
	tokenManager := services.NewJwtTokenManager(services.WithIssuer(appName),
		services.WithTimeout(tokenTimeout),
		services.WithSigningKey([]byte(*cfg.TokenSignatureKey)),
		services.WithVerifyingKey([]byte(*cfg.TokenVerificationKey)))

	authService := services.NewDefaultAuthService(transactionHandler, tokenManager, passwordManager)
	authEndpoint := endpoints.NewDefaultAuthEndpoint(authService)

	drugsService := services.NewDefaultDrugsService(transactionHandler)
	drugsEndpoint := endpoints.NewDefaultDrugsEndpoint(drugsService)
	vaccinationsEndpoint := endpoints.NewDefaultVaccinationsEndpoint()

	// Rest Endpoints

	recoveryFilter := gin.Recovery()
	loggerFilter := sloggin.New(logger.RetrieveLogger().(*slog.Logger).WithGroup("http"))

	publicHandler := gin.New()
	publicHandler.Use(recoveryFilter, loggerFilter)
	publicHandler.POST("/login", authEndpoint.Login)
	publicHandler.POST("/signup", authEndpoint.Signup)
	publicHandler.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	publicHandler.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, feather_web_rest.NotFoundException("resource not found"))
	})
	publicHandler.GET("/info", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"appName": appName})
	})

	privateHandler := publicHandler.Group("/api", authEndpoint.Authorize)
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

	// Application lifecycle

	app := lifecycle.NewApp(
		lifecycle.WithName(appName),
		lifecycle.WithVersion(version),
		lifecycle.WithSignal(syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL),
	)
	app.Attach("HttpServer", feather_web_server.BuildHttpServer(httpServer))
	if err = app.Run(); err != nil {
		feather_commons_log.Fatal(err.Error())
	}
}
