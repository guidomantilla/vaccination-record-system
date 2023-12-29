package datasource

import (
	"log/slog"

	feather_commons_log "github.com/guidomantilla/go-feather-commons/pkg/log"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(logger feather_commons_log.Logger, url string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(url), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: slogGorm.New(
			slogGorm.WithLogger(logger.RetrieveLogger().(*slog.Logger)),
			slogGorm.WithTraceAll(), slogGorm.WithRecordNotFoundError(),
		),
	})
}
