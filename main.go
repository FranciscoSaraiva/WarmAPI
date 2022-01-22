package main

import (
	"log"
	"net/http"
	"os"

	"WarmAPI/config"
	"WarmAPI/handlers"
)

func main() {
	router := config.NewServerAPI()

	/***
	** Endpoints
	***/
	router.GET("/candidate", handlers.RetrieveCandidates)      //Get all the candidates that exist, communicates with the mongo db
	router.GET("/candidate/:id", handlers.RetrieveCandidateID) //Get  a specific candidate via his id, communicates with the mongo db
	router.POST("/candidate", handlers.RetrieveCandidate)      //Get a specific candidate with a name, communicates with  Greenhouse API
	router.POST("/routine", handlers.RetrieveCandidateRoutine)
	// Start server
	log.Println("[ServerAPI] ServerAPI is running in http://localhost:3000/")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal("Error: %v", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
