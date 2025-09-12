package mcp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// SessionManager manages MCP sessions for HTTP transport
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// Session represents an active MCP session
type Session struct {
	ID        string
	CreatedAt time.Time
	LastSeen  time.Time
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a new session and returns its ID
func (sm *SessionManager) CreateSession() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := sm.generateSessionID()
	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	sm.sessions[sessionID] = session
	log.Printf("Created new MCP session: %s", sessionID)
	return sessionID
}

// ValidateSession checks if a session exists and updates its last seen time
func (sm *SessionManager) ValidateSession(sessionID string) bool {
	if sessionID == "" {
		return false
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return false
	}

	session.LastSeen = time.Now()
	return true
}

// RemoveSession removes a session
func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, sessionID)
	log.Printf("Removed MCP session: %s", sessionID)
}

// CleanExpiredSessions removes sessions older than the specified duration
func (sm *SessionManager) CleanExpiredSessions(maxAge time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for id, session := range sm.sessions {
		if now.Sub(session.LastSeen) > maxAge {
			delete(sm.sessions, id)
			log.Printf("Expired MCP session: %s", id)
		}
	}
}

// GetActiveSessionCount returns the number of active sessions
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// generateSessionID creates a cryptographically secure session ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to time-based ID if random fails
		return fmt.Sprintf("mcp-session-%d", time.Now().UnixNano())
	}
	return "mcp-session-" + hex.EncodeToString(bytes)
}

// SessionMiddleware creates HTTP middleware for session management
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("MCP Session Middleware: %s %s", r.Method, r.URL.Path)

		// For GET requests (SSE connections), create a new session
		if r.Method == http.MethodGet {
			sessionID := sm.CreateSession()
			w.Header().Set("Mcp-Session-Id", sessionID)
			w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
			log.Printf("Setting Mcp-Session-Id header: %s", sessionID)
		}

		// For POST requests, validate session ID
		if r.Method == http.MethodPost {
			sessionID := r.Header.Get("Mcp-Session-Id")
			if sessionID == "" {
				sessionID = sm.CreateSession()
				w.Header().Set("Mcp-Session-Id", sessionID)
				w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
				log.Printf("No Mcp-Session-Id header found; created new session: %s", sessionID)
			}

			if !sm.ValidateSession(sessionID) {
				log.Printf("Invalid session ID: %s", sessionID)
				http.Error(w, "session not found or expired", http.StatusNotFound)
				return
			}

			log.Printf("Validated session ID: %s", sessionID)
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
