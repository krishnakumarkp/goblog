package apputil

import "github.com/gorilla/sessions"

type Session struct {
	SessionStore *sessions.CookieStore
}

var AppSession *Session

func NewSession(store *sessions.CookieStore) *Session {
	appSession := new(Session)
	appSession.SessionStore = store
	appSession.SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	return appSession
}
