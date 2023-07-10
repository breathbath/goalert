package app

import (
	"context"
	"crypto/tls"
	"net"
	_ "net/url"
	"strings"

	"github.com/target/goalert/smtpsrv"
)

func (app *App) initSMTPServer(ctx context.Context) error {
	cfg := smtpsrv.Config{}

	if app.cfg.SMTPListenAddr == "" && app.cfg.SMTPListenAddrTLS == "" {
		return nil
	}

	cfg.AllowedDomains = strings.Split(app.cfg.SMTPAllowedDomains, ",")
	cfg.Domain = ""
	cfg.TLSConfig = app.cfg.TLSConfigSMTP

	if app.cfg.SMTPListenAddrTLS != "" {
		cfg.ListenAddr = app.cfg.SMTPListenAddrTLS
		l, err := tls.Listen("tcp", cfg.ListenAddr, cfg.TLSConfig)
		if err != nil {
			return err
		}
		app.smtpsrv = smtpsrv.NewServer(&cfg)
		go func() { _ = app.smtpsrv.Serve(l) }()
	} else {
		cfg.ListenAddr = app.cfg.SMTPListenAddr
		l, err := net.Listen("tcp", cfg.ListenAddr)
		if err != nil {
			return err
		}
		app.smtpsrv = smtpsrv.NewServer(&cfg)
		go func() { _ = app.smtpsrv.Serve(l) }()
	}
	return nil
}
