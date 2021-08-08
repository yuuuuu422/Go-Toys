package main

import (
	pb"g-talk/proto"
	"sync"
	"errors"
)

type token string

type User struct {
	Token  token
	Name   string
	Stream pb.Chat_StreamServer
	Done   chan struct{}
}

type Center struct {
	users map[token]*User
	mutex   sync.RWMutex
}

func NewCenter() *Center {
	return &Center{
		users: make(map[token]*User),
	}
}

func (c *Center) Users() map[token]*User {
	return c.users
}

func (c *Center) Register(tkn token, username string) {
	c.mutex.Lock()
	c.users[tkn] = &User{
		Token: tkn,
		Name:  username,
		Done:  make(chan struct{}), //ctrl+c exit
	}
	c.mutex.Unlock()
}

func (c *Center) Get(tkn token) (*User, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if user, ok := c.users[tkn]; ok {
		return user, nil
	}
	return nil, errors.New("token not valid")
}

func (c *Center) LoginOut(tkn token){
	c.mutex.Lock()
	close(c.users[tkn].Done) // close for <- 
	delete(c.users,tkn)
	c.mutex.Unlock()
}

func (c *Center) Broadcast(resp *pb.StreamResponse){
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, user := range c.users {
		if user.Stream == nil {
			continue
		}
		if err := user.Stream.Send(resp); err != nil {
			c.LoginOut(user.Token)
			continue
		}
	}
}