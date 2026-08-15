package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const CookieName = "rbac_admin_session"

type Session struct {
	UserID    uint
	Username  string
	ExpiresAt time.Time
}

type Manager struct {
	secret   []byte
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewManager(secret string) *Manager {
	return &Manager{
		secret:   []byte(secret),
		sessions: make(map[string]Session),
	}
}

func (m *Manager) Create(c *gin.Context, userID uint, username string) error {
	token, err := randomToken()
	if err != nil {
		return err
	}
	expires := time.Now().Add(8 * time.Hour)

	m.mu.Lock()
	m.sessions[token] = Session{UserID: userID, Username: username, ExpiresAt: expires}
	m.mu.Unlock()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    token + "." + m.sign(token),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
	return nil
}

func (m *Manager) Get(c *gin.Context) (Session, bool) {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return Session{}, false
	}
	parts := strings.Split(cookie, ".")
	if len(parts) != 2 || !m.valid(parts[0], parts[1]) {
		return Session{}, false
	}

	m.mu.RLock()
	session, ok := m.sessions[parts[0]]
	m.mu.RUnlock()
	if !ok || session.ExpiresAt.Before(time.Now()) {
		m.Destroy(c)
		return Session{}, false
	}
	return session, true
}

func (m *Manager) Destroy(c *gin.Context) {
	if cookie, err := c.Cookie(CookieName); err == nil {
		parts := strings.Split(cookie, ".")
		if len(parts) > 0 {
			m.mu.Lock()
			delete(m.sessions, parts[0])
			m.mu.Unlock()
		}
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (m *Manager) sign(token string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func (m *Manager) valid(token, sig string) bool {
	expected := m.sign(token)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
