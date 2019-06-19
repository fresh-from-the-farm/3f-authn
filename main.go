package main

import (
	"context"
	"encoding/json"
	"github.com/fresh-from-the-farm/authn/internal/authn"
	"github.com/fresh-from-the-farm/authn/internal/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Starting auth service...")
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Authn up and running."))
	})

	r.Get("/authenticate", authenticate)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", getAccounts)
		r.Post("/", createAccounts)
		r.Get("/search", searchAccounts)
		r.Route("/{username}", func(r chi.Router) {
			r.Use(AccountCtx)
			r.Get("/", GetAccount)
			r.Put("/", UpdateAccount)
			r.Delete("/", DeleteAccount)
		})
	})

	err := http.ListenAndServe(":3333", r)

	if err != nil {
		log.Fatalf("Authn service initialization failed with %v", err)
	}
}

func getAccounts(w http.ResponseWriter, r *http.Request) {

}
func createAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []model.Account
	err := json.NewDecoder(r.Body).Decode(&accounts)

	if err != nil || len(accounts) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	processed := authn.CreateAccounts(accounts)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(processed)
	if err != nil {
		log.Fatalf("Failed to json encode authn account creation response with %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func searchAccounts(w http.ResponseWriter, r *http.Request) {

}

func AccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		account, err := authn.GetAccount(username)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "account", &account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := ctx.Value("account").(*model.Account)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatalf("Failed to json encode authn account response with %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := ctx.Value("account").(*model.Account)
	w.Header().Set("Content-Type", "application/json")

	var accountUpdate model.Account

	err := json.NewDecoder(r.Body).Decode(&accountUpdate)
	if err != nil {
		log.Fatalf("Failed to json decode authn account response with %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(accountUpdate.Username) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = authn.UpdateAccount(*account, accountUpdate)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatalf("Failed to json encode authn account response with %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := ctx.Value("account").(*model.Account)
	w.Header().Set("Content-Type", "application/json")
	err := authn.PurgeAccount(*account)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	log.Printf("Username %v Password %v", username, password)

	if username == "" || password == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	tkn, ok, err := authn.GetAccessToken(username, password)

	if err != nil {
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if tkn.AccessToken == "" && ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if tkn.AccessToken == "" && !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tkn)
	if err != nil {
		log.Fatalf("Failed to json encode authn token response with %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
