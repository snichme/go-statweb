package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func wrap(action func(req *http.Request) Response) func(http.ResponseWriter, *http.Request) {
	return func(out http.ResponseWriter, req *http.Request) {
		res := action(req)
		if res == nil {
			res = ServerError
		}
		res.WriteTo(out)
	}
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultValue
}

func ShowIndexPage(req *http.Request) Response {
	page := &Page{
		name: "index",
	}
	return page.Render()
}

func ShowNamedPage(req *http.Request) Response {
	vars := mux.Vars(req)
	page := &Page{
		name: PageName(vars["page"]),
	}
	return page.Render()
}

func main() {
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/", wrap(ShowIndexPage))
	mainRouter.HandleFunc("/{page:[a-zA-Z0-9\\/]+}", wrap(ShowNamedPage))

	n := negroni.Classic()
	n.UseHandler(mainRouter)
	n.Run("127.0.0.1:" + getEnv("PORT", "3000"))
}
