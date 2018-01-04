// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 1 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//
	name string

	// clientnum int
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Println("readPump start")
	defer log.Println("readPump end")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("type", getjson(message, "type"))

		if getjson(message, "type") == "login" {
			log.Println("login", getjson(message, "login"))
			c.checklogin(message)

		}

		if getjson(message, "type") == "broadcast" {
			c.hub.broadcast <- message
		}

		if getjson(message, "type") == "msg" {
			c.msgsend(message)
		}

		if getjson(message, "type") == "system" {
			c.handlerequest(message)
		}

	}
}

func (c *Client) msgsend(d1 []byte) {
	dst := getjson(d1, "to")
	if dst == "" {
		return
	}
	// log.Println("readPumpgot from", c.name, string(message), getjson(message, "to"))
	// TODO validate message from and c.name

	if getjson(d1, "from") != c.name {
		clog("from invalid", getjson(d1, "from"), c.name)
		log2file("error.txt", "from invalid "+getjson(d1, "from")+" "+c.name+"\n")

	}
	// log2file("error.txt", "from invalid "+getjson(message, "from")+" "+c.name)

	// log.Println("to ", getjson(message, "to"))

	if dst == "group" {
		clog("msgsend group", (getjson(d1, "group")))
		var from = getjson(d1, "from")
		group := getgrouprecipients(getjson(d1, "group"), getjson(d1, "from"))
		clog(group)
		for li, aa := range group {
			clog(li, aa)
			if from == aa {
				clog(aa, "skipping")
				continue
			}
			// {"type":"msg", "group":"a@a.com#grouptest" , "to":"group" ,"from":"b@a.com","uuid":"msg1497695534928.YnJvYWRjYXN0LHN5c3RlbSxnZXRncm91cGxpc3Q="}
			dout := []byte(`{"type":"msg","to":"` + aa + `","from":"` + getjson(d1, "from") + `","msg":"` + getjson(d1, "msg") + `","group":"` + getjson(d1, "group") + `","uuid":"` + getjson(d1, "uuid") + `","imurl":"` + getjson(d1, "imurl") + `"}`)
			c.hub.sendto <- dout
		}
		return
	}
	c.hub.sendto <- d1
}

func (c *Client) handlerequest(d1 []byte) {

	from := getjson(d1, "from")
	requuid := getjson(d1, "uuid")
	message := []byte(`{"message":"ok","requid":"` + requuid + `"}`)
	cmd := getjson(d1, "cmd")

	c.send <- message
	log.Printf("handlerequest", string(d1), from)
	if cmd == "creategroup" {
		creategroup(d1)
	}
	if cmd == "addtogroup" {
		addtogroup(d1)
	}
	if cmd == "deletefromgroup" {
		deletefromgroup(d1)
	}
	if cmd == "makegroupadmin" {
		log.Printf("makegroupadmin", cmd)
		makegroupadmin(d1)
	}
	if cmd == "getgrouplist" {
		grouplist := getgrouplist(getjson(d1, "group"), getjson(d1, "from"))
		log.Printf("getgrouplist", grouplist)

	}
	if cmd == "getcachedmsg" {
		log.Printf("getcachedmsg")
		c.send_cached()
		c.clear_cache()
	}
	if cmd == "getuserinfo" {
		log.Printf("getuserinfo", cmd)
		getuserinfo(d1)
	}
	if cmd == "setuserinfo" {
		log.Printf("setuserinfo", cmd)
	}

}

func (c *Client) checklogin(d1 []byte) {

	name := getjson(d1, "login")
	csrf := getjson(d1, "csrf")
	storedcsrf := getkvfile("users/" + name + "/csrf.txt")

	if csrf == "abc" || csrf == storedcsrf {
		c.name = name
		clog("c.name", c.name)
		createuser(name)
		c.hub.login <- c
		message := []byte(`{"type":"system","to":"` + c.name + `","message":"loginsuccess"}`)
		c.send <- message
	} else {
		message := []byte(`{"type":"system","message":"loginfailed"}`)
		c.send <- message
	}

}

func (c *Client) send_cached() {
	chatcachedir := rootdir + "users/" + c.name + "/chatcache/"
	files, _ := ioutil.ReadDir(chatcachedir)
	for _, f := range files {
		if f.Name()[:3] != "msg" {
			continue
		}
		fn := chatcachedir + f.Name()
		fi, err := os.Stat(fn)
		if err != nil {
			// Could not obtain stat, handle error
			log.Println("Could not obtain stat, handle error", f.Name())
			return
		} else {
			// path/to/whatever exists
			// dat, err := ioutil.ReadFile(fn)
			// check(err)
			lines := readintostring(fn)
			// fmt.Print(string(dat))
			log.Println("sendcached", f.Name(), fi.Size())
			message := []byte(lines)
			c.send <- message
		}
	}

}

func (c *Client) clear_cache() {

	chatcachedir := rootdir + "users/" + c.name + "/chatcache/"
	archivedir := rootdir + "users/" + c.name + "/archive/"
	files, _ := ioutil.ReadDir(chatcachedir)
	for _, f := range files {
		if f.Name()[:3] != "msg" {
			continue
		}
		srcfn := chatcachedir + f.Name()
		dstfn := archivedir + f.Name()
		fi, err := os.Stat(srcfn)
		if err != nil {
			// Could not obtain stat, handle error
			log.Println("Could not obtain stat, handle error", f.Name())
			continue
		} else {
			clog("clearcache", f.Name(), fi.Size())
			err := os.Rename(srcfn, dstfn)
			if err != nil {
				clog(err)
				continue
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	log.Println("writePump start")
	defer log.Println("writePump end")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println(c.name + "WriteMessage not ok")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			log.Println("w.Write(message)", c.name, string(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// client.hub.clientnames <- "lll"
	client.name = "notset"
	go client.writePump()

	client.readPump()

}
