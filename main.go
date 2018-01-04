// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "homehttp.html")
}

func main() {

	mkdir(rootdir + "users/log/")
	mkdir(rootdir + "users/groups/")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// open a file
	f, errlog := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if errlog != nil {
		fmt.Printf("error opening file: %v", errlog)
	}

	// don't forget to close it
	defer f.Close()

	// assign it to the standard logger
	// log.SetOutput(f)

	flag.Parse()
	hub := newHub()
	go hub.run()
	// http.HandleFunc("/", serveHome)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	// err := http.ListenAndServeTLS(*addr, "k88kfullchain.pem", "k88kprivkey.pem", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
