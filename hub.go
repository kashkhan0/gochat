// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Registered clients.
	clientnames map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Inbound messages from the clients.
	sendto chan []byte

	// Register requests from the clients.
	register chan *Client

	login chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		sendto:      make(chan []byte),
		register:    make(chan *Client),
		login:       make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		clientnames: make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		log.Println("hub waiting")

		select {
		case client := <-h.register:
			log.Println("h.register")
			h.clients[client] = true
		case client := <-h.login:
			log.Println("h.login", client.name, len(h.clients))
			h.clientnames[client.name] = client

		case client := <-h.unregister:
			log.Println("h.unregister")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientnames, client.name)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("h.broadcast")
			// log.Println("h.broadcast", getjson(message, "to"))
			// h.clientnames[getjson(message, "to")].send <- message

			for client := range h.clients {
				select {
				case client.send <- message:
					log.Println("client.send <- message", getjson(message, "to"), getjson(message, "uuid"))
					log.Println("client.name", client.name)
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.sendto:
			// log.Println("hub sendto", string(message))

			if h.clientnames[getjson(message, "to")] != nil {
				log.Println("sending", getjson(message, "to"), getjson(message, "uuid"))

				h.clientnames[getjson(message, "to")].send <- message
			} else {
				log.Println("caching", getjson(message, "to"), getjson(message, "uuid"))

				cachemsg(message)

			}

		}
	}
}
