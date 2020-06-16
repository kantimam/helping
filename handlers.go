package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"transport-status/pkg"
	"transport-status/pkg/database"

	"github.com/gorilla/mux"
)

func getClaimsFromContext(r *http.Request) (string, string) {
	userID := r.Context().Value("userID").(string)
	role := r.Context().Value("role").(string)
	return userID, role
}

type UserData struct {
	Username string
	Password []byte
}

func CreateGetAllResultsHandler(resultData *database.ResultsData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := resultData.GetAllResults()
		fmt.Println(results)
		if err != nil {
			log.Fatalf("could not get all results from database %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func CreateAddResultHandler(resultData *database.ResultsData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, _ := getClaimsFromContext(r)
		var result database.Result

		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not parse new result body", http.StatusInternalServerError)
			return
		}
		result.Username = username
		t := time.Now()
		formattedString := t.Format("2006-01-02 15:04:05")
		result.ResultTime = formattedString

		if result.GasTankFilled == true && result.ExternalDamage == false && result.WashingNeeded == false && result.TechnicalRepair == false {
			result.VehicleState = true
		} else {
			result.VehicleState = false
		}

		fmt.Printf("%+v", result)

		err = resultData.AddResult(result)
		if err != nil {
			fmt.Printf("could not add result to database %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func CreateGetAllTransportsHandler(transportData *database.VehicleData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transports, err := transportData.GetTransports()
		if err != nil {
			log.Fatalf("could not get all results from database %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(&transports)
		if err != nil {
			http.Error(w, "could not marshal to json", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

func RouteHandler(transportData *database.VehicleData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		transports, err := transportData.GetRoute(id)
		if err != nil {
			log.Fatalf("could not get all results from database %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(&transports)
		if err != nil {
			http.Error(w, "could not marshal to json", http.StatusInternalServerError)
			return
		}
		w.Write(b)

	}
}

func CreateGetAllUsersHandler(userData *database.UserData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userData.GetUsers()
		if err != nil {
			log.Fatalf("could not get all results from database %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

type AddUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func CreateAddUserHandler(userData *database.UserData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, role := getClaimsFromContext(r)
		if role != "admin" {
			fmt.Println("only admin can create a user")
			http.Error(w, "could not parse add user body", http.StatusForbidden)
			return
		}
		var addUserPayload AddUserPayload
		err := json.NewDecoder(r.Body).Decode(&addUserPayload)
		fmt.Println(&addUserPayload)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not parse add user body", http.StatusInternalServerError)
			return
		}

		passwordHash, err := pkg.HashAndSalt([]byte(addUserPayload.Password))
		err = userData.AddUser(addUserPayload.Username, passwordHash, addUserPayload.Role)
		if err != nil {
			log.Fatalf("could not add a new user %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type LoginUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateLoginUser(userData *database.UserData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginUserPayload LoginUserPayload
		err := json.NewDecoder(r.Body).Decode(&loginUserPayload)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not parse payload", http.StatusInternalServerError)
			return
		}

		user, err := userData.GetUser(loginUserPayload.Username)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "user is invalid", http.StatusInternalServerError)
			return
		}
		token, err := pkg.CreateToken(user.Username, user.Role)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(token))
	}
}
