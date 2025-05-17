package dotosrv

import (
	"net"
	"net/http"

	dotoapi "github.com/DistributedShenanigans/doto/api"
	"github.com/DistributedShenanigans/doto/config"
)

func New(cfg config.Serving, si dotoapi.ServerInterface) *http.Server {
	return &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: dotoapi.HandlerFromMux(si, http.NewServeMux()),
	}
}
