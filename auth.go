package main

import (
	"net/http"
)

const (
	eschecker = "eschecker"
)

/*
** Functions used to implement a basic authentification for the server displaying
** the queries stats (see stats.go). It's a basic HTTP authentication, the
** credentials to access the page are to be set in the YAML.
 */

type BasicAuth struct {
	Login    string
	Password string
	//this is the function that will be launched in case of login success
	Display func(http.ResponseWriter, *http.Request)
}

func NewBasicAuth(login, passwd string) *BasicAuth {
	return &BasicAuth{Login: login, Password: passwd}
}

func (a *BasicAuth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+eschecker+`"`)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (a *BasicAuth) ValidAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return username == a.Login && password == a.Password
}

func (a *BasicAuth) BasicAuthHandler(w http.ResponseWriter, r *http.Request) {
	if !a.ValidAuth(r) {
		a.Authenticate(w, r)
	} else {
		a.Display(w, r)
	}
}

func (a *BasicAuth) setDisplayFunc(display func(http.ResponseWriter, *http.Request)) {
	a.Display = display
}
