package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetWebSocketUpgrader يُرجع upgrader مع إعدادات CheckOrigin ديناميكية
func GetWebSocketUpgrader(allowedOrigins string) websocket.Upgrader {
	origins := strings.Split(allowedOrigins, ",")

	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// السماح بطلبات بدون Origin (مثل curl) - مع تحذير
				utils.LogWarning("WebSocket request without Origin header - allowing for tools like curl", map[string]interface{}{
					"origin": origin,
				})
				return true
			}

			// التحقق من الدومينات المسموح بها
			for _, allowed := range origins {
				allowed = strings.TrimSpace(allowed)
				
				// رفض صريح لـ '*' في أي حالة
				if allowed == "*" {
					utils.LogError("Security error: Rejecting '*' in WebSocket CheckOrigin - use specific domains", map[string]interface{}{
						"allowed_origin": allowed,
					})
					continue // تجاهل هذا الإدخال والانتقال للقادم
				}
				
				// السماح بـ localhost و 127.0.0.1 في التطوير
				if strings.Contains(allowed, "localhost") || strings.Contains(allowed, "127.0.0.1") {
					if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
						utils.LogInfo("WebSocket connection allowed from localhost", map[string]interface{}{
							"origin": origin,
						})
						return true
					}
				}
				
				// التحقق من تطابق الدومين
				if strings.HasPrefix(origin, allowed) {
					utils.LogInfo("WebSocket connection allowed", map[string]interface{}{
						"origin": origin,
						"allowed": allowed,
					})
					return true
				}
			}

			utils.LogWarning("WebSocket connection rejected from unauthorized origin", map[string]interface{}{
				"origin": origin,
			})
			return false
		},
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
}

// WebSocketHub - مركز إدارة اتصالات WebSocket
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	maxClients int
}

// WSHandler - معالج WebSocket
type WSHandler struct {
	hub            *WebSocketHub
	upgrader       websocket.Upgrader
	allowedOrigins string
}

// NewWSHandler - إنشاء معالج WebSocket جديد
func NewWSHandler(allowedOrigins string) *WSHandler {
	hub := &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256), // Buffered channel to prevent blocking
		register:   make(chan *websocket.Conn, 256),
		unregister: make(chan *websocket.Conn, 256),
		maxClients: 100, // Limit concurrent WebSocket connections
	}

	handler := &WSHandler{
		hub:            hub,
		upgrader:       GetWebSocketUpgrader(allowedOrigins),
		allowedOrigins: allowedOrigins,
	}
	go hub.run()

	return handler
}

// run - تشغيل مركز WebSocket
func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if len(h.clients) >= h.maxClients {
				utils.LogWarning("WebSocket connection rejected - max clients reached", map[string]interface{}{
					"max_clients": h.maxClients,
					"current_clients": len(h.clients),
				})
				client.Close()
			} else {
				h.clients[client] = true
				utils.LogInfo("WebSocket client registered", map[string]interface{}{
					"total_clients": len(h.clients),
				})
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			utils.LogInfo("WebSocket client unregistered", map[string]interface{}{
				"total_clients": len(h.clients),
			})

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					utils.LogError("WebSocket send error", map[string]interface{}{
						"error": err,
					})
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// HandleWebSocket - معالجة اتصال WebSocket
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	lang := i18n.Detect(c)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.LogError("Failed to upgrade connection to WebSocket", map[string]interface{}{
			"error": err,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	h.hub.register <- conn

	// إرسال رسالة ترحيب
	welcomeMsg := map[string]interface{}{
		"type":      "connected",
		"message":   "تم الاتصال بنجاح",
		"timestamp": time.Now(),
	}
	msgBytes, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)

	// الاستماع للرسائل من العميل
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.LogError("Recovered from panic in WebSocket handler", map[string]interface{}{
					"panic": r,
				})
			}
			h.hub.unregister <- conn
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				utils.LogWarning("WebSocket read error", map[string]interface{}{
					"error": err,
				})
				break
			}
		}
	}()
}

// BroadcastLocationUpdate - إرسال تحديث الموقع لجميع العملاء
func (h *WSHandler) BroadcastLocationUpdate(update map[string]interface{}) {
	message := map[string]interface{}{
		"type":      "location_update",
		"data":      update,
		"timestamp": time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		utils.LogError("Failed to encode location update", map[string]interface{}{
			"error": err,
		})
		return
	}

	h.hub.broadcast <- msgBytes
}

// GetClientCount - عدد العملاء المتصلين
func (h *WSHandler) GetClientCount() int {
	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()
	return len(h.hub.clients)
}

// BroadcastEmployeeStatus - إرسال تحديث حالة الموظف
func (h *WSHandler) BroadcastEmployeeStatus(update map[string]interface{}) {
	message := map[string]interface{}{
		"type":      "employee_status",
		"data":      update,
		"timestamp": time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		utils.LogError("Failed to encode employee status update", map[string]interface{}{
			"error": err,
		})
		return
	}

	h.hub.broadcast <- msgBytes
}
