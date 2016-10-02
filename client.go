package main

import (
	"github.com/gorilla/websocket"
	"log"
)

/* To figure out if they wanna broadcast to all or broadcast to all except them */
type Message struct {
	mtype string
	msg   []byte
}

/* Reads and writes messages from client */
type Client struct {
	conn *websocket.Conn
	in   chan []byte
	out  chan Message
}

/* Reads and pumps to out channel */
func (c *Client) ReadLoop() {
	defer close(c.out)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := Message{"ex", message}
		c.out <- msg
	}
}

/* Reads from in channel and pumps to client */
func (c *Client) WriteLoop() {
	defer close(c.in)
	for {
		in := <-c.in
		c.WriteMessage(in)
	}
}

/* Writes a message to the client */
func (c *Client) WriteMessage(msg []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
	}
}

/* Constructor */
func NewClient(conn *websocket.Conn) *Client {
	client := new(Client)
	client.in = make(chan []byte)
	client.out = make(chan Message)
	return client
}
