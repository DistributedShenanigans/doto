package dotosrv

import (
	"fmt"
	"net/http"

	"github.com/DistributedShenanigans/doto/config"
)

func New(cfg config.Serving, h http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: h,
	}
}
