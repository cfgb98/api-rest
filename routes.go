package main

import (
	"net/http"

	"github.com/gorilla/mux" //enrutador
)

//Route representa una ruta en el navegor web
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

//Routes es un slice de Route
type Routes []Route

//NewRouter retorna una variable que maneja el router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true) // url con /
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{"index", "GET", "/", index},
	Route{"movieList", "GET", "/peliculas", movieList},
	Route{"movieShow", "GET", "/pelicula/{id}", movieShow},
	Route{"movieAdd", "POST", "/pelicula", movieAdd},
	Route{"movieUpdate", "PUT", "/pelicula/{id}", movieUpdate},
	Route{"movieRemove", "DELETE", "/pelicula/{id}", movieRemove},
	Route{"userAdd", "GET", "/signup", userAdd},
}
