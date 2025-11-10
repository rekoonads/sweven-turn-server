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
	// Get TURN credentials from environment variables
	turnUsername := os.Getenv("TURN_USERNAME")
	turnPassword := os.Getenv("TURN_PASSWORD")
	publicIP := os.Getenv("PUBLIC_IP")
	portStr := os.Getenv("TURN_PORT")

	if turnUsername == "" || turnPassword == "" {
		panic(fmt.Errorf("TURN_USERNAME and TURN_PASSWORD environment variables are required"))
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
		panic(err)
	}

	fmt.Println("TURN server started successfully")

	// Start HTTP health check endpoint for Railway
	healthPort := os.Getenv("PORT")
	if healthPort == "" {
		healthPort = "8080" // Default health check port
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TURN server is running"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"turn-server"}`))
	})

	go func() {
		fmt.Printf("Health check endpoint started on port %s\n", healthPort)
		if err := http.ListenAndServe(":"+healthPort, nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// Block until user sends SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Shutting down TURN server...")

	if err = s.Close(); err != nil {
		log.Panic(err)
	}
}
