package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type JsonResponse struct {
	Type    string     `json:"type"`
	Data    []Property `json:"data"`
	Message string     `json:"message"`
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetProperties(w http.ResponseWriter, r *http.Request) {
	printMessage("Getting properties...")

	db := setupDB()

	// Get all properties from table
	rows, err := db.Query("SELECT * FROM properties")

	checkErr(err)

	var properties []Property

	for rows.Next() {
		var id int
		var address string

		err = rows.Scan(&id, &address)

		checkErr(err)

		properties = append(properties, Property{PropertyID: id, Address: address})
	}

	var response = JsonResponse{Type: "success", Data: properties}

	json.NewEncoder(w).Encode(response)
}

// create a property
func CreateProperty(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")

	var response = JsonResponse{}

	if address == "" {
		response = JsonResponse{Type: "error", Message: "You are missing an address"}
	} else {
		db := setupDB()

		printMessage("inserting into db")

		var lastInsertID int
		err := db.QueryRow("INSERT INTO properties(address) VALUES($1) returning id;", address).Scan(&lastInsertID)
		// check errors
		checkErr(err)
		fmt.Println(lastInsertID)
		var createdProperty []Property
		rows, err := db.Query("SELECT * FROM properties WHERE id=($1)", lastInsertID)

		checkErr(err)

		for rows.Next() {
			var id int
			var address string

			err = rows.Scan(&id, &address)

			checkErr(err)

			createdProperty = append(createdProperty, Property{PropertyID: id, Address: address})
		}
		response = JsonResponse{Type: "success", Data: createdProperty, Message: "The property has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}

func main() {
	router := mux.NewRouter()

	// Get all properties
	router.HandleFunc("/properties/", GetProperties).Methods("GET")

	// Create a property
	router.HandleFunc("/property/", CreateProperty).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
