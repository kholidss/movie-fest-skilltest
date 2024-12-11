package bootstrap

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/pkg/cdn"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"strings"
)

func RegistryCDN(cfg *config.Config) cdn.CDN {
	var (
		cdnInstance cdn.CDN
		err         error
	)
	switch strings.ToLower(cfg.CDNConfig.Provider) {
	case consts.CDNProviderCloudinary:
		cdnInstance, err = cdn.NewCloudinaryCDN(cfg)
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed initiate CDN cloudinary: %v", err))
		}
	default:
		logger.Fatal(fmt.Sprintf("invalid CDN provider: %v", cfg.CDNConfig.Provider))
	}
	return cdnInstance
}
