package zeaburcaddyextension

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(ZeaburExtension{})
	httpcaddyfile.RegisterHandlerDirective("zeaburextension", func(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
		return &ZeaburExtension{}, nil
	})
	httpcaddyfile.RegisterDirectiveOrder("zeaburextension", httpcaddyfile.After, "header")
}

type ZeaburExtension struct {
	headerConfig  map[string]HeaderConfig
	redirectRules map[string]RedirectRule
}

// Provision implements caddy.Provisioner.
func (z *ZeaburExtension) Provision(ctx caddy.Context) error {
	fsys := ctx.FileSystems().Default()

	if content, err := fs.ReadFile(fsys, "_headers"); err == nil {
		slog.Info("found _headers file")

		hc, err := ParseHeaderConfig(string(content))
		if err != nil {
			return fmt.Errorf("parse headers: %w", err)
		}

		hcm := make(map[string]HeaderConfig, len(hc))
		for _, h := range hc {
			hcm[strings.TrimRight(h.Path, "/")+"/"] = h
		}

		z.headerConfig = hcm
	}

	if content, err := fs.ReadFile(fsys, "_redirects"); err == nil {
		slog.Info("found _redirects file")

		rr, err := ParseRedirects(string(content))
		if err != nil {
			return fmt.Errorf("parse redirects: %w", err)
		}

		rrm := make(map[string]RedirectRule, len(rr))
		for _, r := range rr {
			rrm[strings.TrimRight(r.SourcePath, "/")+"/"] = r
		}

		z.redirectRules = rrm
	}

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (z *ZeaburExtension) ServeHTTP(w http.ResponseWriter, r *http.Request, h caddyhttp.Handler) error {
	path := strings.TrimRight(r.URL.Path, "/") + "/"

	// if the path matches the exact redirect rules, redirect
	if rule, ok := z.redirectRules[path]; ok {
		slog.Info("redirecting", "source", r.URL.Path, "target", rule.TargetPath, "status", rule.StatusCode)
		http.Redirect(w, r, rule.TargetPath, rule.StatusCode)
		return nil
	}

	// if the path matches the prefix redirect rules, redirect
	if rule, ok := z.headerConfig[path]; ok {
		slog.Info("applying headers", "path", r.URL.Path, "headers", rule.Headers)
		for k, v := range rule.Headers {
			w.Header().Set(k, v)
		}
	}

	return h.ServeHTTP(w, r)
}

// CaddyModule returns the Caddy module information.
func (ZeaburExtension) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.zeaburextension",
		New: func() caddy.Module { return new(ZeaburExtension) },
	}
}

var (
	_ caddy.Provisioner           = (*ZeaburExtension)(nil)
	_ caddyhttp.MiddlewareHandler = (*ZeaburExtension)(nil)
)
