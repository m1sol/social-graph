package main

import (
	"fmt"
	"sync"
)

type User struct {
	ID   int
	Name string
}

type Graph struct {
	mu sync.Mutex

	users       map[int]*User
	connections map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		users:       make(map[int]*User),
		connections: make(map[int][]int),
	}
}

func (g *Graph) AddUser(id int, name string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.users[id] = &User{
		ID:   id,
		Name: name,
	}
}

func (g *Graph) GetUser(id int) (*User, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	user, ok := g.users[id]
	return user, ok
}

func (g *Graph) AddConnection(fromID, toID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, fromOk := g.users[fromID]
	_, toOk := g.users[toID]
	if !fromOk || !toOk || fromID == toID {
		return false
	}

	g.connections[fromID] = append(g.connections[fromID], toID)
	g.connections[toID] = append(g.connections[toID], fromID)

	return true
}

func (g *Graph) GetConnections(userID int) []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[userID]) == 0 {
		return nil
	}

	cUsers := make([]*User, 0, len(g.connections[userID]))

	for _, id := range g.connections[userID] {
		cUsers = append(cUsers, g.users[id])
	}

	return cUsers
}

func (g *Graph) HasConnection(fromID, toID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[fromID]) == 0 || len(g.connections[toID]) == 0 || fromID == toID {
		return false
	}

	//микрооптимизация - ищем по меньшему массиву
	if len(g.connections[fromID]) > len(g.connections[toID]) {
		fromID, toID = toID, fromID
	}

	for _, id := range g.connections[fromID] {
		if id == toID {
			return true
		}
	}

	return false
}

func (g *Graph) UserCount() int {
	return len(g.users)
}

func main() {
	graph := NewGraph()

	graph.AddUser(1, "Alice")
	graph.AddUser(2, "Bob")
	graph.AddUser(3, "Charlie")

	graph.AddConnection(1, 2) // Alice -> Bob
	graph.AddConnection(1, 3) // Alice -> Charlie
	graph.AddConnection(2, 3) // Bob -> Charlie

	if user, ok := graph.GetUser(1); ok {
		fmt.Printf("User: %s\n", user.Name)
		friends := graph.GetConnections(1)
		fmt.Printf("Friends: %d\n", len(friends))
		for _, friend := range friends {
			fmt.Printf("  - %s\n", friend.Name)
		}
	}

	fmt.Printf("Alice and Bob connected: %v\n",
		graph.HasConnection(1, 2))
}
