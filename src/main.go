package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"size:64;not null" json:"firstName"`
	LastName  string `gorm:"size:64;not null" json:"lastName"`
	Email     string `gorm:"unique;size:128;not null" json:"email"`
	Password  string `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Response struct {
	Message string `json:"message"`
}

var (
	sessions = make(map[string]string)
	mu       sync.Mutex
	db       *gorm.DB
)

func initDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Database environment variables are not set properly")
	}
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{"Invalid request body"})
		return
	}
	var existing User
	if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(Response{"User already exists"})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{"Failed to create user"})
		return
	}
	json.NewEncoder(w).Encode(Response{"Register successful"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var reqUser User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{"Invalid request body"})
		return
	}
	var user User
	if err := db.Where("email = ?", reqUser.Email).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{"Invalid credentials"})
		return
	}
	if user.Password != reqUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{"Invalid credentials"})
		return
	}
	sessionID := user.Email + "_session"
	mu.Lock()
	sessions[sessionID] = user.Email
	mu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})
	json.NewEncoder(w).Encode(Response{"Login successful"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{"No session found"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	json.NewEncoder(w).Encode(Response{"Logout successful"})
}

func main() {
	initDB()
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	log.Println("Server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
