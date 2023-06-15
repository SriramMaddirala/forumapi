package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Post struct {
	postId       int64
	posterId     string
	postDate     string
	commId       string
	parentPostId string
	textContent  string
	mediaLinks   string
	eventId      string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func addData(w http.ResponseWriter, r *http.Request) {
	var decodedRequest Post
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&decodedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	sqlStatement := `INSERT INTO forum (postId, posterId, postDate, commId, parentPostId, textContent, mediaLinks, eventId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(sqlStatement, decodedRequest.postId, decodedRequest.postId, decodedRequest.posterId, decodedRequest.postDate, decodedRequest.commId, decodedRequest.parentPostId, decodedRequest.textContent, decodedRequest.mediaLinks, decodedRequest.eventId)
	if err != nil {
		fmt.Println("Issue with DB")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":1025", nil)
}
func connectDB() {
	godotenv.Load()
	host, hostError := os.LookupEnv("DB_HOST")
	if !hostError {
		panic("Couldn't get DB host")
	}
	user, userError := os.LookupEnv("DB_USER")
	if !userError {
		panic("Couldn't get DB username")
	}
	password, passwordError := os.LookupEnv("DB_PASSWORD")
	if !passwordError {
		panic("Couldn't get DB password")
	}
	port, portError := strconv.Atoi(os.Getenv("DB_PORT"))
	if portError != nil {
		panic(portError)
	}
	name, nameError := os.LookupEnv("DB_NAME")
	if !nameError {
		panic("Couldn't get DB name")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
func main() {
	connectDB()
	handleRequests()
}
