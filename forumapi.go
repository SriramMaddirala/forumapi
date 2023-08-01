package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Post struct {
	PosterId     string `json:"PosterId"`
	CommId       string `json:"CommId"`
	ParentPostId string `json:"ParentPostId"`
	TextContent  string `json:"TextContent"`
	EventId      string `json:"EventId"`
	MediaLinks   string `json:"MediaLinks"`
}
type PostRow struct {
	PostId       int64
	PosterId     string
	CommId       string
	ParentPostId string
	TextContent  string
	MediaLinks   string
	EventId      string
	PostDate     string
}
type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}
type UserLogin struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type UserRow struct {
	PosterId string
}

func getPost(w http.ResponseWriter, r *http.Request) {
	sqlStatement := `SELECT * FROM forum WHERE postid = $1`
	QueryParams := r.URL.Query()
	postid := QueryParams.Get("postid")
	rows, err := db.Query(sqlStatement, postid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	var rowsData PostRow
	for rows.Next() {
		var (
			PostId       int64
			PosterId     string
			CommId       string
			ParentPostId string
			MediaLinks   string
			TextContent  string
			EventId      string
			PostDate     string
		)
		if err := rows.Scan(&PostId, &PosterId, &PostDate, &CommId, &ParentPostId, &TextContent, &MediaLinks, &EventId); err != nil {
			log.Fatal(err)
		}
		rowsData = PostRow{PostId: PostId, PosterId: PosterId, PostDate: PostDate, CommId: CommId, ParentPostId: ParentPostId, TextContent: TextContent, MediaLinks: MediaLinks, EventId: EventId}
	}
	result, error := json.Marshal(rowsData)
	if error != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Write(result)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	deleteStatement := `DELETE FROM forum WHERE postid = $1`
	selectStatement := `SELECT * FROM forum WHERE postid = $1`
	QueryParams := r.URL.Query()
	postid := QueryParams.Get("postid")
	rows, err := db.Query(selectStatement, postid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	for rows.Next() {
		var (
			PostId       int64
			PosterId     string
			CommId       string
			ParentPostId string
			MediaLinks   string
			TextContent  string
			EventId      string
			PostDate     string
		)
		if err := rows.Scan(&PostId, &PosterId, &PostDate, &CommId, &ParentPostId, &TextContent, &MediaLinks, &EventId); err != nil {
			log.Fatal(err)
		}
		os.Remove(MediaLinks)
	}
	_, err = db.Query(deleteStatement, postid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.WriteHeader(http.StatusNoContent)
}

func getPoster(w http.ResponseWriter, r *http.Request) {
	posterid := r.URL.Query().Get("posterid")
	sqlStatement := `SELECT * FROM forum WHERE posterid = $1`
	rows, err := db.Query(sqlStatement, posterid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	var rowsData []PostRow
	for rows.Next() {
		var (
			PostId       int64
			PosterId     string
			CommId       string
			ParentPostId string
			MediaLinks   string
			TextContent  string
			EventId      string
			PostDate     string
		)
		if err := rows.Scan(&PostId, &PosterId, &PostDate, &CommId, &ParentPostId, &TextContent, &MediaLinks, &EventId); err != nil {
			log.Fatal(err)
		}
		rowsData = append(rowsData, PostRow{PostId: PostId, PosterId: PosterId, PostDate: PostDate, CommId: CommId, ParentPostId: ParentPostId, TextContent: TextContent, MediaLinks: MediaLinks, EventId: EventId})
	}
	result, error := json.Marshal(rowsData)
	if error != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Write(result)
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func uploadBlob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var ErrUnexpectedEOF = errors.New("unexpected EOF")
	body, err := io.ReadAll(r.Body)
	if err != nil && err != ErrUnexpectedEOF {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	var contType string = strings.Split(r.Header.Get("Content-type"), "/")[2]
	var filename string = "/Users/ram/Pictures/" + r.URL.Query().Get("medianame") + "." + contType
	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		http.Error(w, "Failed to save file to filesystem", http.StatusInternalServerError)
		return
	}
	if contType != "webp" {
		convertimg(filename)
	}
	w.WriteHeader(http.StatusNoContent)
}
func addData(w http.ResponseWriter, r *http.Request) {
	var decodedRequest Post
	var media string
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&decodedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	time := time.Now().UTC()
	if decodedRequest.MediaLinks != "" {
		media = "/Users/ram/Pictures/" + decodedRequest.MediaLinks + ".webp"
	}
	sqlStatement := `INSERT INTO forum (PostId, PosterId, PostDate, CommId, ParentPostId, TextContent, MediaLinks, EventId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(sqlStatement, generateSnowflake(time), decodedRequest.PosterId, time.String(), decodedRequest.CommId, decodedRequest.ParentPostId, decodedRequest.TextContent, media, decodedRequest.EventId)
	if err != nil {
		fmt.Println("Issue with DB")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
func getFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	mediapath := r.URL.Query().Get("mediapath")
	http.ServeFile(w, r, mediapath)
}
func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var decodedRequest User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&decodedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	sqlStatement := `INSERT INTO users (PosterId, JoinDate, Username, Pword, Email) VALUES ($1, $2, $3, $4, $5)`
	time := time.Now().UTC()
	var posterid int64 = generateSnowflake(time)
	_, err = db.Exec(sqlStatement, posterid, time.String(), decodedRequest.Username, decodedRequest.Password, decodedRequest.Email)
	if err != nil {
		fmt.Println("Issue with DB")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	result, error := json.Marshal(UserRow{PosterId: strconv.Itoa(int(posterid))})
	if error != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(result)
}
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var decodedRequest UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&decodedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	sqlStatement := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	rows, err := db.Query(sqlStatement, decodedRequest.Username)
	if err != nil {
		fmt.Println("Issue with DB")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	var rowsData UserRow
	for rows.Next() {
		var (
			PosterId string
			JoinDate string
			Username string
			Password string
			Email    string
		)
		if err := rows.Scan(&PosterId, &JoinDate, &Username, &Password, &Email); err != nil {
			log.Fatal(err)
		}
		rowsData = UserRow{PosterId: PosterId}
	}
	result, error := json.Marshal(rowsData)
	if error != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(result)
}
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", addData)
	http.HandleFunc("/getposter", getPoster)
	http.HandleFunc("/getpost", getPost)
	http.HandleFunc("/upload", uploadBlob)
	http.HandleFunc("/getfile", getFile)
	http.HandleFunc("/deletepost", deletePost)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
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
func init() {
	connectDB()
}
func main() {
	handleRequests()
}
