package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clientMap            map[string]*Client
	channelToUserListMap map[string]map[string]bool
	pubsub               *redis.PubSub

	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	mgr := &Manager{
		clientMap:            make(map[string]*Client),
		channelToUserListMap: make(map[string]map[string]bool),
		pubsub:               NewRedisPubSub(),
	}

	go mgr.sendNotification()

	return mgr
}

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	// Request has an ID, as in "/task/<id>".
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) < 2 {
		http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
		return
	}

	userId := pathParts[1]

	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Obtain list of channels the user is subscribed to
	membership := fetchMembership(userId)

	// Create New Client
	client := NewClient(userId, membership, conn, m)
	// Add the newly created client to the manager
	m.addClient(userId, membership, client)
	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(userId string, membership []string, client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clientMap[userId] = client

	for _, channel := range membership {
		if _, ok := m.channelToUserListMap[channel]; !ok {
			m.channelToUserListMap[channel] = make(map[string]bool)
			m.pubsub.Subscribe(context.Background(), channel)
		}

		m.channelToUserListMap[channel][userId] = true
	}
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if client == nil {
		return
	}

	userId, membership := client.userId, client.membership

	// Check if Client exists, then delete it
	if _, ok := m.clientMap[userId]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clientMap, userId)
	}

	for _, channel := range membership {
		if _, ok := m.channelToUserListMap[channel]; ok {
			delete(m.channelToUserListMap[channel], userId)
		}

		// Unsubscribe channel if no user are part of it
		if len(m.channelToUserListMap[channel]) == 0 {
			m.pubsub.Unsubscribe(context.Background(), channel)
			delete(m.channelToUserListMap, channel)
		}
	}
}

func (m *Manager) sendNotification() {
	for {
		msg, err := m.pubsub.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		m.RLock()

		go func() {
			defer m.RUnlock()
			fmt.Println(msg.Channel, msg.Payload)

			channel, payload := msg.Channel, msg.Payload
			userListMap := m.channelToUserListMap[channel]

			for userId, _ := range userListMap {
				client := m.clientMap[userId]
				client.egress <- []byte(payload)
			}
		}()
	}
}

func NewRedisPubSub() *redis.PubSub {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pubsub := rdb.Subscribe(context.Background())

	return pubsub
}
