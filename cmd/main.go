package main

import (
	"edgeturn"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	threadNum = 4
	realm     = "thinkmay.net"

	min = 60000
	max = 65535
)

var (
	proj     = "https://supabase.thinkmay.net"
	anon_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ewogICJyb2xlIjogImFub24iLAogICJpc3MiOiAic3VwYWJhc2UiLAogICJpYXQiOiAxNjk0MDE5NjAwLAogICJleHAiOiAxODUxODcyNDAwCn0.EpUhNso-BMFvAJLjYbomIddyFfN--u-zCf0Swj9Ac6E"
)

func init() {
	project := os.Getenv("TM_PROJECT")
	key := os.Getenv("TM_ANONKEY")
	if project != "" {
		proj = project
	}
	if key != "" {
		anon_key = key
	}
}

func main() {
	// Start HTTP health check endpoint FIRST for Railway
	// This ensures the service appears healthy even during initialization
	healthPort := os.Getenv("PORT")
	if healthPort == "" {
		healthPort = "8080" // Default health check port
	}

	serverReady := false
	serverError := false
	var serverErrorMsg string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if serverError {
			w.Write([]byte("TURN server failed to start: " + serverErrorMsg))
		} else if serverReady {
			w.Write([]byte("TURN server is running"))
		} else {
			w.Write([]byte("TURN server is initializing..."))
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		status := "initializing"
		if serverError {
			status = "error"
		} else if serverReady {
			status = "healthy"
		}
		w.Write([]byte(fmt.Sprintf(`{"status":"%s","service":"turn-server"}`, status)))
	})

	go func() {
		fmt.Printf("Health check endpoint started on port %s\n", healthPort)
		if err := http.ListenAndServe(":"+healthPort, nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// Give health check server time to start
	time.Sleep(100 * time.Millisecond)

	// Get TURN credentials from environment variables
	turnUsername := os.Getenv("TURN_USERNAME")
	turnPassword := os.Getenv("TURN_PASSWORD")
	publicIP := os.Getenv("PUBLIC_IP")
	portStr := os.Getenv("TURN_PORT")

	// Check if required environment variables are missing
	if turnUsername == "" || turnPassword == "" {
		errorMsg := "TURN_USERNAME and TURN_PASSWORD must be set"
		log.Println(errorMsg)
		serverError = true
		serverErrorMsg = errorMsg
		
		// Keep the application running so health check endpoint stays alive
		select {} // Block forever
	}

	if publicIP == "" {
		publicIP = "0.0.0.0" // Default to all interfaces
	}

	port := 3478 // Default TURN port
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	fmt.Printf("Starting TURN server on %s:%d\n", publicIP, port)
	fmt.Printf("TURN credentials: %s / %s\n", turnUsername, turnPassword)

	// Optional: Supabase ping (if worker ID is provided)
	workerID := os.Getenv("WORKER_ID")
	if workerID != "" {
		go func() {
			agent := edgeturn.NewSupabaseAgent(
				fmt.Sprintf("https://%s/", proj),
				anon_key)
			for {
				err := agent.Ping(workerID)
				if err != nil {
					fmt.Println("Supabase ping error:", err.Error())
				}
				time.Sleep(10 * time.Second)
			}
		}()
	}

	s, err := edgeturn.SetupTurn(publicIP, turnUsername, turnPassword, port, min, max)
	if err != nil {
		errorMsg := fmt.Sprintf("TURN server initialization error: %v", err)
		log.Println(errorMsg)
		log.Println("Health check will continue to run, but TURN functionality is disabled")
		serverError = true
		serverErrorMsg = errorMsg
		
		// Keep the application running so health check endpoint stays alive
		select {} // Block forever
	} else {
		fmt.Println("TURN server started successfully")
		serverReady = true
		
		// Block until user sends SIGINT or SIGTERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		fmt.Println("Shutting down TURN server...")

		if err = s.Close(); err != nil {
			log.Panic(err)
		}
	}
}