package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // السماح بجميع المصادر في وضع التطوير
	},
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
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
	hub *WebSocketHub
}

// NewWSHandler - إنشاء معالج WebSocket جديد
func NewWSHandler() *WSHandler {
	hub := &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256), // Buffered channel to prevent blocking
		register:   make(chan *websocket.Conn, 256),
		unregister: make(chan *websocket.Conn, 256),
		maxClients: 100, // Limit concurrent WebSocket connections
	}

	handler := &WSHandler{hub: hub}
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
				log.Printf("⚠️ تم رفض اتصال WebSocket جديد - الحد الأقصى %d", h.maxClients)
				client.Close()
			} else {
				h.clients[client] = true
				log.Println("✅ تم تسجيل عميل WebSocket جديد")
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("❌ تم إلغاء تسجيل عميل WebSocket")

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("❌ خطأ في إرسال WebSocket: %v", err)
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
	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ فشل ترقية الاتصال إلى WebSocket: %v", err)
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
			h.hub.unregister <- conn
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
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
		log.Printf("❌ خطأ في ترميز تحديث الموقع: %v", err)
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
		log.Printf("❌ خطأ في ترميز تحديث الحالة: %v", err)
		return
	}

	h.hub.broadcast <- msgBytes
}