package main

import (
	"log"
	"net/http"
)

func main() {
	// статичні файли
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// маршрути для головної сторінки
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	
	// маршрути для калькуляторів
	SetupCalculator1Routes()
	SetupCalculator2Routes()
	
	// сервер
	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}