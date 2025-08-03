package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wrestler094/launchpad/internal/services"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	authService    *services.AuthService
	tokenService   *services.TokenService
	presaleService *services.PresaleService
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewHandlers creates new API handlers
func NewHandlers(authService *services.AuthService, tokenService *services.TokenService, presaleService *services.PresaleService) *Handlers {
	return &Handlers{
		authService:    authService,
		tokenService:   tokenService,
		presaleService: presaleService,
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, ErrorResponse{Error: message})
}

// respondSuccess sends a success response
func respondSuccess(w http.ResponseWriter, message string, data interface{}) {
	respondJSON(w, http.StatusOK, SuccessResponse{Message: message, Data: data})
}

// GenerateNonce generates a nonce for MetaMask authentication
func (h *Handlers) GenerateNonce(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		respondError(w, http.StatusBadRequest, "address parameter is required")
		return
	}

	nonce, err := h.authService.GenerateNonce(address)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, "Nonce generated", nonce)
}

// Login handles user authentication
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, "Login successful", response)
}

// VerifyToken verifies a JWT token
func (h *Handlers) VerifyToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.authService.VerifyToken(tokenString)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, "Token is valid", claims)
}

// AuthMiddleware is a middleware for authenticating requests
func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := h.authService.VerifyToken(tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user address to request context
		ctx := r.Context()
		ctx = addUserToContext(ctx, claims.Address)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// CreateToken handles token creation
func (h *Handlers) CreateToken(w http.ResponseWriter, r *http.Request) {
	userAddress := getUserFromContext(r.Context())
	if userAddress == "" {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.tokenService.CreateToken(userAddress, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, "Token created successfully", response)
}

// GetToken handles getting a token by address
func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	if address == "" {
		respondError(w, http.StatusBadRequest, "Token address is required")
		return
	}

	token, err := h.tokenService.GetToken(address)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, "Token retrieved", token)
}

// ListTokens handles listing tokens for a user
func (h *Handlers) ListTokens(w http.ResponseWriter, r *http.Request) {
	userAddress := getUserFromContext(r.Context())
	if userAddress == "" {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	tokens, err := h.tokenService.ListTokens(userAddress)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, "Tokens retrieved", tokens)
}

// CreatePresale handles presale creation
func (h *Handlers) CreatePresale(w http.ResponseWriter, r *http.Request) {
	userAddress := getUserFromContext(r.Context())
	if userAddress == "" {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.CreatePresaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.presaleService.CreatePresale(userAddress, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, "Presale created successfully", response)
}

// GetPresale handles getting a presale by ID
func (h *Handlers) GetPresale(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid presale ID")
		return
	}

	presale, err := h.presaleService.GetPresale(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, "Presale retrieved", presale)
}

// GetPublicPresale handles getting public presale info (for landing pages)
func (h *Handlers) GetPublicPresale(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid presale ID")
		return
	}

	presale, err := h.presaleService.GetPresale(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Get token info as well
	token, err := h.tokenService.GetToken(presale.TokenAddress)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get token info")
		return
	}

	data := map[string]interface{}{
		"presale": presale,
		"token":   token,
	}

	respondSuccess(w, "Presale info retrieved", data)
}

// ListPresales handles listing presales for a user
func (h *Handlers) ListPresales(w http.ResponseWriter, r *http.Request) {
	userAddress := getUserFromContext(r.Context())
	if userAddress == "" {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	presales, err := h.presaleService.ListPresales(userAddress)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, "Presales retrieved", presales)
}

// ParticipateInPresale handles presale participation
func (h *Handlers) ParticipateInPresale(w http.ResponseWriter, r *http.Request) {
	userAddress := getUserFromContext(r.Context())
	if userAddress == "" {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid presale ID")
		return
	}

	var req services.ParticipateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.presaleService.ParticipateInPresale(id, userAddress, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, "Participation recorded", response)
}