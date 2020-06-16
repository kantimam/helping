package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"transport-status/pkg"
	"transport-status/pkg/database"
	"transport-status/pkg/handlers"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Config struct {
	FilesUrl struct {
		ReisaiUrl string `json:"ReisaiUrl"`
	} `json:"FilesUrl"`
}

func LoadConfiguration(filename string) (Config, error) {

	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, errors.Wrap(err, "failed to open config file")
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, errors.Wrap(err, "failed to decode config file")
	}
	return config, err
}

func DownloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to get data file url")
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy data from data file")
	}
	return nil
}

func init() {
	// logging setup
	logFile, err := os.OpenFile("../errorLog.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(errors.Wrap(err, "Failed to create/append errorLog file"))
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// config setup
	config, err := LoadConfiguration("../config.json")
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to load configuration file"))
	}

	ReisaiUrl := config.FilesUrl.ReisaiUrl
	DownloadFile("../reisai.txt", ReisaiUrl)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		auth := r.Header.Get("Authorization")
		if auth == "" {
			fmt.Println(errors.New("Forbidden, auth header not set"))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if strings.HasPrefix(auth, prefix) {
			auth = strings.TrimPrefix(auth, prefix)
		}

		userID, role, err := pkg.ValidateToken(auth)
		if err != nil {
			fmt.Println(errors.Wrap(err, "could not validate token"))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		next(w, r)
	})
}

func main() {

	r := mux.NewRouter()

	// initialize all database connections
	resultsDB := database.ConnectionToResultsDatabase()
	resultData, err := database.CreateResultDatabase(resultsDB)
	if err != nil {
		log.Fatalf("could not setup results databases %w", err)
	}
	transportDB := database.ConnectionToTransportDatabase()
	transportData, err := database.CreateTransportDatabase(transportDB)
	if err != nil {
		log.Fatalf("could not setup transports databases %w", err)
	}
	userDB := database.ConnectionToUsersDatabase()
	userData, err := database.CreateUserDatabase(userDB)
	if err != nil {
		log.Fatalf("could not setup user databases %w", err)
	}

	// NON AUTHENTICATED ROUTES
	r.HandleFunc("/login", CORSMiddleware(handlers.CreateLoginUser(userData))).Methods(http.MethodPost)
	r.HandleFunc("/login", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodOptions)

	// AUTHENTICATED ROUTES
	// results
	r.HandleFunc("/results", AuthMiddleware(CORSMiddleware(handlers.CreateGetAllResultsHandler(resultData)))).Methods(http.MethodGet)
	r.HandleFunc("/results", AuthMiddleware(CORSMiddleware(handlers.CreateAddResultHandler(resultData)))).Methods(http.MethodPost)
	r.HandleFunc("/results", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodOptions)

	// get route_id from busNumber
	r.HandleFunc("/routes/{id}", AuthMiddleware(CORSMiddleware(handlers.RouteHandler(transportData)))).Methods(http.MethodGet)
	r.HandleFunc("/routes/{id}", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodOptions)

	// transport
	r.HandleFunc("/transports", AuthMiddleware(CORSMiddleware(handlers.CreateGetAllTransportsHandler(transportData)))).Methods(http.MethodGet)
	r.HandleFunc("/transports", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodOptions)

	// users
	r.HandleFunc("/admin/users", AuthMiddleware(CORSMiddleware(handlers.CreateGetAllUsersHandler(userData)))).Methods(http.MethodGet)
	r.HandleFunc("/admin/users", AuthMiddleware(CORSMiddleware(handlers.CreateAddUserHandler(userData)))).Methods(http.MethodPost)
	r.HandleFunc("/admin/users", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodOptions)
	fmt.Println("running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}
