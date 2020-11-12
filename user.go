package main

//User es el modelo de datos para los usuarios
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
