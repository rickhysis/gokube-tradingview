package main

import (
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "./tv"
)


// main function to boot up everything
func main() {
    router := mux.NewRouter()
    s := router.PathPrefix("/api/v1").Subrouter()

    s.HandleFunc("/config", tv.GetConfig).Methods("GET")
    s.HandleFunc("/time", tv.GetTime).Methods("GET")
    s.HandleFunc("/symbols", tv.GetSymbol).Methods("GET")
    s.HandleFunc("/history", tv.GetHistory).Methods("GET")
    s.HandleFunc("/quotes", tv.GetQuotes).Methods("GET")

    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    // start server listen
    // with error handling
    log.Fatal(http.ListenAndServe(":6969", handlers.CORS(originsOk, headersOk, methodsOk)(router)))

}
