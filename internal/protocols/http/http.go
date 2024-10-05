package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/internal/graceful"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http/router"
	"github.com/rs/zerolog/log"

	"net/http/pprof"
	// _ "github.com/nurcahyaari/ecommerce/docs"
	// httpswagger "github.com/swaggo/http-swagger"
)

type Http struct {
	cfg         config.Config
	HttpRouter  *router.HttpRoute
	httpServer  *http.Server
	serverState graceful.ServerState
}

func New(cfg config.Config, httpRouter *router.HttpRoute) *Http {
	return &Http{
		cfg:        cfg,
		HttpRouter: httpRouter,
	}
}

func (p *Http) router(app *chi.Mux) {
	p.HttpRouter.Router(app)
}

func (h *Http) cors(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
}

func (h *Http) swagger(app *chi.Mux) {
	// app.Mount("/swagger", httpswagger.WrapHandler)
	app.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func (h *Http) shutdownStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h.serverState {
		case graceful.StateShutdown:
			// response.Json[string](w, http.StatusInternalServerError, "server is shutting down", "")
			response.Json[string](w,
				response.SetHttpCode[string](http.StatusInternalServerError),
				response.SetMessage[string]("server is shutting down"))
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Http) pprof(r *chi.Mux) {
	if !h.cfg.Application.EnablePprof {
		return
	}

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/heap", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func (p *Http) Listen() {
	app := chi.NewRouter()

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(p.shutdownStateMiddleware)
	p.cors(app)
	p.router(app)
	p.swagger(app)
	p.pprof(app)

	serverPort := fmt.Sprintf(":%d", p.cfg.Application.Transport.Http.PORT)
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	log.Info().Msgf("Server started on Port %s ", serverPort)
	if err := p.httpServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("cannot establish HTTP connection")
	}
}

func (h *Http) Shutdown(ctx context.Context) error {
	h.serverState = graceful.StateShutdown
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
