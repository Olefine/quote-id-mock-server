package domain

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

// Broadcaster to broadcast to log events into []*ws.Conn
type Broadcaster struct {
	clients map[string]*websocket.Conn
	sync.Mutex
}

// Join room
func (caster *Broadcaster) Join(ws *websocket.Conn) {
	// TODO remove this shit
	caster.Lock()
	clientUUID := ws.Request().URL.Query().Get("uuid")
	_, e := caster.clients[clientUUID]
	if !e {
		caster.clients[clientUUID] = ws
	}
	caster.Unlock()
}

// Broadcast writes data to clients
func (caster *Broadcaster) Broadcast(p []byte) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(caster.clients))

	for name := range caster.clients {
		fmt.Println(name)
	}

	for name, conn := range caster.clients {
		go func(nm string, c *websocket.Conn) {
			defer waitGroup.Done()
			_, err := c.Write(p)

			if err != nil && strings.Contains(err.Error(), "broken pipe") {
				caster.Lock()
				if _, e := caster.clients[nm]; e {
					caster.clients[nm].Close()
					delete(caster.clients, nm)
				}
				caster.Unlock()
			}
		}(name, conn)

	}

	waitGroup.Wait()
}

// NewBroadcaster constructs broadcaster to hold clients
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{clients: make(map[string]*websocket.Conn)}
}
