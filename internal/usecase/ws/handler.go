package ws

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"main/internal/services/order"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"github.com/redis/go-redis/v9"
// )

// type Cache struct {
// 	host    string
// 	db      int
// 	expires time.Duration
// }

// func NewCache(host string, db int, expires time.Duration) *Cache {
// 	return &Cache{
// 		host:    host,
// 		db:      db,
// 		expires: expires}
// }

// func GetClient() *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	})
// }

// func (ch *Cache) Set(ctx context.Context, key string, value interface{}) error {
// 	client := GetClient()
// 	v, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	err = client.Set(ctx, key, v, ch.expires*time.Second).Err()
// 	return err
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// func (m *WSManager) HandleDriverWS(ctx *gin.Context) {
// 	driverIDStr := ctx.Query("driver_id")

// 	driverID, err := strconv.Atoi(driverIDStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver_id"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		log.Println("upgrade error:", err)
// 		return
// 	}

// 	m.AddDriver(driverID, conn)
// 	defer func() {
// 		m.RemoveDriver(driverID)
// 		conn.Close()
// 	}()

// 	cacheStore := NewCache("localhost:6379", 0, 3600)

// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		var location order.Location

// 		if err := json.Unmarshal(msg, &location); err != nil {
// 			log.Println("invalid location data:", err)
// 			continue
// 		}

// 		err = cacheStore.Set(context.Background(),
// 			fmt.Sprintf("driver:%d:location", driverID),
// 			location,
// 		)

// 		if err != nil {
// 			log.Println("error setting location:", err)
// 		}
// 	}
// }

// func (m *WSManager) HandleClientWS(ctx *gin.Context) {
// 	clientStr := ctx.Query("client_id")
// 	clientId, err := strconv.Atoi(clientStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_id"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade error:", err)
// 		return
// 	}

// 	m.AddClient(clientId, conn)
// 	defer func() {
// 		m.RemoveClient(clientId)
// 		conn.Close()
// 		log.Printf("Client %d disconnected", clientId)
// 	}()

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Client read error:", err)
// 			break
// 		}
// 		log.Printf("Message from client %d: %s", clientId, message)
// 	}
// }
