package mongodb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte // message chan send to user
	Groups []int       // groups the client belongs to
}

type Hub struct {
	Clients    map[int][]*Client // 建立组id与组成员的映射
	Broadcast  chan WSopt
	Register   chan *Client
	Unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[int][]*Client),
		Broadcast:  make(chan WSopt),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			for _, val := range client.Groups {
				cls := h.Clients[val]
				cls = append(cls, client)
				h.Clients[val] = cls
			}
		case client := <-h.Unregister:
			for _, val := range client.Groups {
				cls := h.Clients[val]
				var opt []*Client
				for j, cli := range cls {
					if cli == client {
						opt = append(opt, cls[:j]...)
						opt = append(opt, cls[j+1:]...)
						h.Clients[val] = opt
					}
				}
			}
			close(client.Send)
		case eventList := <-h.Broadcast:
			groupID := eventList.GroupID
			message, err := json.Marshal(eventList)
			if err != nil {
				fmt.Println("json marshal ", err)
				return
			}
			var leftList []*Client
			for _, client := range h.Clients[groupID] {
				select {
				case client.Send <- message:
					leftList = append(leftList, client)
				default:
					close(client.Send)
				}
			}
			// 在clients中删除对应连接
			h.Clients[groupID] = leftList
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
