package ws

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	// "main/internal/services/order"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// type WSManager struct {
// 	mu          sync.Mutex
// 	driverConns map[int]*websocket.Conn
// 	clientConns map[int]*websocket.Conn
// 	repo        Repository
// }

// func NewManager(repo Repository) *WSManager {
// 	return &WSManager{
// 		driverConns: make(map[int]*websocket.Conn),
// 		clientConns: make(map[int]*websocket.Conn),
// 		repo:        repo,
// 	}
// }

// func (m *WSManager) AddDriver(driverID int, conn *websocket.Conn) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.driverConns[driverID] = conn
// 	log.Printf("Driver %d connected\n", driverID)
// }

// func (m *WSManager) AddClient(clientId int, conn *websocket.Conn) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.clientConns[clientId] = conn
// 	log.Printf("Client %d connected\n", clientId)
// }

// func (m *WSManager) RemoveDriver(driverID int) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	delete(m.driverConns, driverID)
// 	log.Printf("Driver %d disconnected\n", driverID)
// }

// func (m *WSManager) RemoveClient(clientId int) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	delete(m.clientConns, clientId)
// 	log.Printf("Client %d disconnected\n", clientId)
// }

// func (m *WSManager) SendToDriver(driverID int, event string, data any) {
// 	m.mu.Lock()
// 	conn, ok := m.driverConns[driverID]
// 	m.mu.Unlock()
// 	if !ok {
// 		return
// 	}

// 	message := map[string]any{
// 		"event": event,
// 		"data":  data,
// 	}
// 	bytes, _ := json.Marshal(message)
// 	conn.WriteMessage(websocket.TextMessage, bytes)
// }

// func (m *WSManager) SendToClient(clientID int, event string, data any) {
// 	m.mu.Lock()
// 	conn, ok := m.clientConns[clientID]
// 	m.mu.Unlock()
// 	if !ok {
// 		log.Printf("Client %d not connected\n", clientID)
// 		return
// 	}

// 	message := map[string]any{
// 		"event": event,
// 		"data":  data,
// 	}

// 	bytes, _ := json.Marshal(message)
// 	conn.WriteMessage(websocket.TextMessage, bytes)
// }

// func (m *WSManager) Broadcast(event string, data any) {
// 	orderData, ok := data.(order.OrderDetail)
// 	if !ok {
// 		log.Println("Broadcast: data is not order.Create")
// 		return
// 	}

// 	fromDistrictId := orderData.FromDistrictId
// 	toDistrictId := orderData.ToDistrictId

// 	driverIds, err := m.repo.GetDriverListByDistrictId(context.Background(), fromDistrictId, toDistrictId)
// 	if err != nil {
// 		log.Printf("Error getting driver list: %v", err)
// 		return
// 	}

// 	bytes, _ := json.Marshal(map[string]any{
// 		"event": event,
// 		"data":  data,
// 	})

// 	m.mu.Lock()

// 	for _, driverId := range driverIds {
// 		if conn, ok := m.driverConns[driverId]; ok {
// 			conn.WriteMessage(websocket.TextMessage, bytes)
// 		}
// 	}
// 	m.mu.Unlock()
// }

// func (m *WSManager) BroadcastExcept(exceptID int, event string, data any) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	for id, conn := range m.driverConns {
// 		if id == exceptID {
// 			continue
// 		}
// 		bytes, _ := json.Marshal(map[string]any{
// 			"event":    event,
// 			"order_id": data,
// 		})
// 		conn.WriteMessage(websocket.TextMessage, bytes)
// 	}
// }
