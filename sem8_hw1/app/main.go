package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type LogMessage struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
	Pod    string `json:"pod"`
}

var logFile = "/app/logs/app.log"

func init() {
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		if _, err := os.Create(logFile); err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Welcome to the Go app in %s", hostname)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logMsg LogMessage
	if err := json.NewDecoder(r.Body).Decode(&logMsg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s [INFO] %s\n", timestamp, logMsg.Message)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Failed to read log file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/logs", logsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 