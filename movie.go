package main

// Movie es una repesentacion de una pelicula
type Movie struct {
	Name     string `json:"name"`
	Year     int    `json:"year"`
	Category string `json:"category"`
	Director string `json:"director"`
}

//Movies es un slice de Movie
type Movies []Movie
