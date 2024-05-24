package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TalisonK/TalisonContabil/src/util"
	"github.com/markbates/goth/gothic"
)

func AuthProviderCallback(w http.ResponseWriter, r *http.Request) {

	provider := r.PathValue("provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		util.LogHandler("Erro ao autenticar usuário", err, "AuthProviderCallback")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	fmt.Println(user)

	util.LogHandler("Usuário autenticado com sucesso", nil, "AuthProviderCallback")

	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutProvider(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.Logout(w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}

func AuthProvider(w http.ResponseWriter, r *http.Request) {

	provider := r.PathValue("provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {

		fmt.Println(gothUser)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}