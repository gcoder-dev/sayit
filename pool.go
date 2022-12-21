package main

import "github.com/gorilla/websocket"

type Pool struct {
	pool map[*websocket.Conn]bool
}

func (p *Pool) register(c *websocket.Conn) {
	p.pool[c] = true
}

func (p *Pool) unregister(c *websocket.Conn) {
	if ok, _ := p.pool[c]; ok == true {
		delete(p.pool, c)
	}
}
