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
}

// WebSocketHub - مركز إدارة اتصالات WebSocket
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

// WSHandler - معالج WebSocket
type WSHandler struct {
	hub *WebSocketHub
}

// NewWSHandler - إنشاء معالج WebSocket جديد
func NewWSHandler() *WSHandler {
	hub := &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
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
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("✅ تم تسجيل عميل WebSocket جديد")

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