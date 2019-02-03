package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var R *mux.Router

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	RootRouter := mux.NewRouter()
	http.Handle("/", RootRouter)
	R = RootRouter.PathPrefix("/rpc").Subrouter()
	RootRouter.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logrus.Debugf("%v %s - %s", time.Now(), r.URL.String(), r.UserAgent())
				next.ServeHTTP(w, r)
			})
		})
	R.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-API-Server", "SenRen2")
				next.ServeHTTP(w, r)
			})
		})
	logrus.Info("Init global router")
}
