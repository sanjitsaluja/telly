package telly

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"net/http"
)

// Route configuration
//
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{"Index", "GET", "/", serveHome},
	Route{"Index", "GET", "/channels/{id}.m3u8", serveVideoMpg},
	Route{"Index", "GET", "/channels/{id}.ts", serveVideoTS},
	Route{"Index", "POST", "/channels/{id}/tune", serveTuneChannel},
	Route{"Index", "GET", "/channels/{id}/status", serveChannelStatus},
	Route{"Index", "GET", "/channels/{id}", serveChannelDetails},
}

func InitHttpHandlers() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	return router
}

func StartServer() {
	r := InitHttpHandlers()
	n := negroni.New()

	// Panic recovery
	n.Use(negroni.NewRecovery())

	// Logging
	n.Use(negronilogrus.NewMiddleware())

	// Router goes last
	n.UseHandler(r)

	// Run http server
	n.Run(fmt.Sprintf(":%s", AppConfig.Port))
}
