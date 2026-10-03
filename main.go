package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/amonaco/goauth/lib/auth"
	"github.com/amonaco/goauth/lib/cache"
	"github.com/amonaco/goauth/lib/config"
	"github.com/amonaco/goauth/lib/middleware"
	"github.com/amonaco/goauth/lib/session"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.Read("config/config.yml")
	conf := config.Get()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cache.Start()
	defer cache.Close()

	router := chi.NewRouter()
	router.Use(middleware.RateLimit(60, time.Minute))
	router.Use(middleware.Middleware)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		auth.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserID(r.Context())
		if err != nil {
			slog.Warn("secure route denied", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		auth.JSONResponse(w, http.StatusOK, map[string]any{"user_id": userID})
	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			UserID    uint32   `json:"user_id"`
			CompanyID uint32   `json:"company_id"`
			Roles     []string `json:"roles"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if payload.UserID == 0 {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		sess, token, err := auth.IssueSession(w, r, payload.UserID, payload.CompanyID, payload.Roles)
		if err != nil {
			slog.Error("login failed", "error", err.Error())
			http.Error(w, "login failed", http.StatusInternalServerError)
			return
		}

		auth.JSONResponse(w, http.StatusOK, map[string]any{"token": token, "session": sess})
	})

	router.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		if err := auth.Logout(w, r); err != nil {
			slog.Warn("logout failed", "error", err.Error())
			http.Error(w, "logout failed", http.StatusInternalServerError)
			return
		}
		auth.JSONResponse(w, http.StatusOK, map[string]string{"status": "logged_out"})
	})

	router.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if payload.Token == "" {
			token, err := r.Cookie(auth.TokenCookieName)
			if err == nil {
				payload.Token = token.Value
			}
		}
		if payload.Token == "" {
			http.Error(w, "token is required", http.StatusBadRequest)
			return
		}

		refreshed, err := auth.RefreshSession(payload.Token)
		if err != nil {
			slog.Warn("refresh failed", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		auth.SetSessionCookie(w, refreshed, conf.SessionTTL, conf.CookieSecure)
		auth.JSONResponse(w, http.StatusOK, map[string]string{"token": refreshed})
	})

	router.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(session.ContextKey("session")).(session.Session)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		auth.JSONResponse(w, http.StatusOK, map[string]any{"user_id": sess.UserID, "company_id": sess.CompanyID, "roles": sess.Roles})
	})

	if conf.Listen == "" {
		conf.Listen = "0.0.0.0:80"
	}
	if conf.Redis == "" {
		panic(errors.New("redis connection string is required"))
	}
	if conf.JWTSecret == "" || conf.JWTSecret == "change-me-in-production" {
		panic(errors.New("jwt_secret must be set to a secure value"))
	}

	slog.Info("starting server", "listen", conf.Listen)
	if err := http.ListenAndServe(conf.Listen, router); err != nil {
		slog.Error("server failed", "error", err.Error())
		os.Exit(1)
	}
	_ = strconv.Itoa(0)
}
