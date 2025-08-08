package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/corstijank/mekstrike/gamemaster/internal/config"
	"github.com/corstijank/mekstrike/gamemaster/internal/handlers"
	"github.com/corstijank/mekstrike/gamemaster/internal/middleware"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	"github.com/corstijank/mekstrike/gamemaster/internal/services"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer client.Close()
	
	// Initialize repository
	gameRepo := repository.NewMemoryGameRepository()
	
	// Initialize services
	gameService := services.NewGameService(client, gameRepo)
	eventService := services.NewEventService(client, gameRepo)
	
	// Initialize handlers
	gameHandlers := handlers.NewGameHandlers(client, gameRepo)
	eventHandlers := handlers.NewEventHandlers(client, gameRepo, cfg)
	
	// Setup router
	r := mux.NewRouter()
	
	// Apply middleware
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	
	// Game routes
	r.HandleFunc("/games/{id}", gameHandlers.GetGame).Methods("GET")
	r.HandleFunc("/games/{id}/availableActions", gameHandlers.GetAvailableActions).Methods("GET")
	r.HandleFunc("/games/{id}/board", gameHandlers.GetBoard).Methods("GET")
	r.HandleFunc("/games/{id}/units/{uid}", gameHandlers.GetUnit).Methods("GET")
	r.HandleFunc("/games/{id}/advanceTurn", gameHandlers.AdvanceTurn).Methods("POST")
	r.HandleFunc("/games", gameHandlers.NewGame).Methods("POST")
	r.HandleFunc("/games", gameHandlers.GetGames).Methods("GET")
	
	// Dapr pub/sub event handlers
	r.HandleFunc("/dapr/subscribe", eventHandlers.GetDaprSubscriptions).Methods("GET")
	r.HandleFunc("/unit-movement-completed", eventHandlers.HandleMovementCompleted).Methods("POST")
	r.HandleFunc("/unit-attack-completed", eventHandlers.HandleAttackCompleted).Methods("POST")
	r.HandleFunc("/unit-end-phase-completed", eventHandlers.HandleEndPhaseCompleted).Methods("POST")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	log.Printf("Starting Mekstrike gamemaster on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
	
	// Note: gameService and eventService are available but not used directly in main
	// They can be used for future enhancements or testing
	_ = gameService
	_ = eventService
}