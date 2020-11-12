package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux" //enrutador
	"gopkg.in/mgo.v2"        //mongodb
	"gopkg.in/mgo.v2/bson"   //bson formato en que se guardan los datos en mongodb
)

//Message representa un mensaje de respuesta de la API
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//metodo de la estructura
func (this *Message) setStatus(data string) {
	this.Status = data
}

func (this *Message) setMessage(data string) {
	this.Message = data
}

var collection = getSession().DB("curso_go").C("movies") //.C es la colecccion, se cre sola la db, la coleccion

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return session
}

func responseMovie(w http.ResponseWriter, status int, results Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func responseMovies(w http.ResponseWriter, status int, results []Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola desde servidor golang")
}

//mostrar todas las peliculas
func movieList(w http.ResponseWriter, r *http.Request) {
	var results []Movie
	err := collection.Find(nil).Sort("-_id").All(&results) //encontrar todos los  registros, se guardan en results, Sort("-_id") los ordena desde el mas reciente al mas viejo

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("resultados: ", results)
	}
	responseMovies(w, 200, results)
}

//mostrar una pelicula
func movieShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //obtener parametros por url
	movieID := params["id"]
	//comprobar si no es un id en hexadecimal porque asi son las claves primarias de mongo
	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404) //recurso no encontrado
		return
	}
	oid := bson.ObjectIdHex(movieID)
	results := Movie{}
	err := collection.FindId(oid).One(&results) //buscar por id

	if err != nil {
		w.WriteHeader(404) //recurso no encontrado
		return
	}
	responseMovie(w, 200, results)
}

//agregar una pelicula
func movieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body) //decodificar json
	var movieData Movie
	err := decoder.Decode(&movieData) //decodificar variable
	if err != nil {
		panic(err) //mostar error y parar servidor
	}
	defer r.Body.Close() //cerrar

	err = collection.Insert(movieData)
	if err != nil {
		w.WriteHeader(500) //error 500 es un fallo del servidor
		return
	}

	responseMovie(w, 200, movieData)
}

func movieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //obtener parametros por url
	movieID := params["id"]

	//comprobar si no es un id en hexadecimal porque asi son las claves primarias de mongo
	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404) //recurso no encontrado
		return
	}
	oid := bson.ObjectIdHex(movieID)
	decoder := json.NewDecoder(r.Body)
	var movieData Movie
	err := decoder.Decode(&movieData) //guardar resultado en movieData

	if err != nil {
		panic(err)
		w.WriteHeader(500) //error 500 es un fallo del servidor
		return
	}
	defer r.Body.Close()
	results := Movie{}
	err = collection.FindId(oid).One(&results) //buscar por id
	document := bson.M{"_id": oid}             //obtener id para actulizar
	change := bson.M{"$set": movieData}
	err = collection.Update(document, change) //ejecutar consulta

	if err != nil {
		panic(err)
		w.WriteHeader(404) //error 404 recurso no encontrado
		return
	}

	responseMovie(w, 200, movieData)
}

//mostrar una pelicula
func movieRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //obtener parametros por url
	movieID := params["id"]
	//comprobar si no es un id en hexadecimal porque asi son las claves primarias de mongo
	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404) //recurso no encontrado
		return
	}
	oid := bson.ObjectIdHex(movieID)
	err := collection.RemoveId(oid) //borrar por id

	if err != nil {
		w.WriteHeader(404) //recurso no encontrado
		return
	}
	message := new(Message) // new retorna puntero a la estructura
	message.setStatus("Success")
	message.setMessage("La película con ID" + movieID + " ha sido borrada correctamente")
	results := message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200) //200 estatus ok, sin errores
	json.NewEncoder(w).Encode(results)
}
