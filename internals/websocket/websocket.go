package websocket

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type WSService struct {
	redis *redis.Client
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

type RiderEvent struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type BaseWSMessage struct {
	Type string `json:"type"`
}

func NewWSService(redis *redis.Client) *WSService {
	return &WSService{
		redis: redis,
		conns: make(map[string]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (ws *WSService) addConn(riderID string, conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, ok := ws.conns[riderID]; ok {
		return
	}

	ws.conns[riderID] = conn
}

func (ws *WSService) removeConn(riderID string) {
	ws.mu.Lock()
	delete(ws.conns, riderID)
	ws.mu.Unlock()
}
func (ws *WSService) WSHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade failed:", err)
		return
	}

	riderID := r.URL.Query().Get("rider_id")
	if riderID == "" {
		log.Println("missing rider_id")
		conn.Close()
		return
	}

	ws.addConn(riderID, conn)
	defer ws.removeConn(riderID)

	ws.subscribeRider(ctx, riderID)

	go ws.subscribeTripLocation(ctx, riderID, riderID)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("ws disconnected:", riderID)
			return
		}
	}
}

func (ws *WSService) send(riderID string, payload []byte) {
	ws.mu.RLock()
	conn := ws.conns[riderID]
	ws.mu.RUnlock()

	if conn == nil {
		return
	}

	_ = conn.WriteMessage(websocket.TextMessage, payload)
}

func ensureRiderGroup(ctx context.Context, rds *redis.Client, stream string) {
	err := rds.XGroupCreateMkStream(ctx, stream, "ws-group", "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Println("[xgroup create]", err)
	}
}

func (ws *WSService) subscribeRider(ctx context.Context, riderID string) {
	stream := "RIDER:" + riderID
	group := "ws-group"
	consumer := uuid.NewString()

	ensureRiderGroup(ctx, ws.redis, stream)

	go func() {
		for {

			res, err := ws.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{stream, ">"},
				Count:    5,
				Block:    0,
			}).Result()

			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Println("[xreadgroup]", err)
				continue
			}

			for _, s := range res {
				for _, msg := range s.Messages {

					raw, ok := msg.Values["data"]
					if !ok {
						continue
					}

					var data string
					switch v := raw.(type) {
					case string:
						data = v
					case []byte:
						data = string(v)
					default:
						continue
					}

					ws.send(riderID, []byte(data))

					ws.redis.XAck(ctx, stream, group, msg.ID)
				}
			}
		}
	}()
}

func (ws *WSService) subscribeTripLocation(
	ctx context.Context,
	tripID string,
	riderID string,
) {
	channel := "trip:" + tripID

	sub := ws.redis.Subscribe(ctx, channel)
	defer sub.Close()

	ch := sub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				return
			}

			// forward live location to rider
			ws.send(riderID, []byte(msg.Payload))

		case <-ctx.Done():
			return
		}
	}
}
