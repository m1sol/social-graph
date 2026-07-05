package main

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type Connection interface {
	Type() string
	Weight() int
}

type Friend struct {
	Since string // Дата начала дружбы
}

func (f Friend) Type() string {
	return "friend"
}

func (f Friend) Weight() int {
	return 10
}

type Follower struct {
	Notifications bool
}

func (f Follower) Type() string {
	return "follower"
}

func (f Follower) Weight() int {
	return 5
}

type Blocked struct {
	Reason string
}

func (b Blocked) Type() string {
	return "blocked"
}

func (b Blocked) Weight() int {
	return -1
}

type User struct {
	ID   int
	Name string
}

type Graph struct {
	mu sync.Mutex

	users       map[int]*User
	connections map[int]map[int]Connection
}

func NewGraph() *Graph {
	return &Graph{
		users:       make(map[int]*User),
		connections: make(map[int]map[int]Connection),
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

	if _, ok := g.connections[fromID]; !ok {
		g.connections[fromID] = make(map[int]Connection)
	}

	g.connections[fromID][toID] = Follower{}

	return true
}

func (g *Graph) GetConnections(userID int) []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[userID]) == 0 {
		return nil
	}

	cUsers := make([]*User, 0, len(g.connections[userID]))
	ids := slices.Collect(maps.Keys(g.connections[userID]))

	for _, id := range ids {
		cUsers = append(cUsers, g.users[id])
	}

	return cUsers
}

func (g *Graph) HasConnection(fromID, toID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[fromID]) == 0 || fromID == toID {
		return false
	}

	_, ok := g.connections[fromID][toID]

	return ok
}

func (g *Graph) IsMutual(fromID, toID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[fromID]) == 0 || len(g.connections[toID]) == 0 || fromID == toID {
		return false
	}

	_, fromOk := g.connections[fromID][toID]
	_, toOk := g.connections[toID][fromID]

	return fromOk && toOk
}

func (g *Graph) UserCount() int {
	return len(g.users)
}

func (g *Graph) RemoveConnection(fromID, toID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.connections[fromID][toID]; ok {
		delete(g.connections[fromID], toID)
		return true
	}

	return false
}

func (g *Graph) RemoveUser(id int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.users[id]; !ok {
		return false
	}

	delete(g.users, id)
	delete(g.connections, id)

	for fromID := range g.connections[id] {
		delete(g.connections[fromID], id)
	}

	return false
}

func (g *Graph) ConnectionCount(userID int) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	return len(g.connections[userID])
}

func (g *Graph) CommonConnections(id1, id2 int) []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	var common []*User

	if len(g.connections[id1]) == 0 || len(g.connections[id2]) == 0 {
		return nil
	}

	if len(g.connections[id1]) > len(g.connections[id2]) {
		id1, id2 = id2, id1
	}

	for id := range g.connections[id1] {
		if _, ok := g.connections[id2][id]; !ok {
			common = append(common, g.users[id])
		}
	}

	return common
}

func (g *Graph) SuggestConnections(userID int) []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	var suggest []*User
	seen := make(map[int]struct{})

	for user, _ := range g.connections[userID] {
		for id, _ := range g.connections[user] {
			if _, ok := g.connections[userID][id]; !ok && id != userID {
				if _, exists := seen[id]; !exists {
					suggest = append(suggest, g.users[id])
					seen[id] = struct{}{}
				}
			}
		}
	}

	return suggest
}

func (g *Graph) GetAllUsers() []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	users := make([]*User, 0, len(g.users))

	for _, user := range g.users {
		users = append(users, user)
	}

	return users
}

func (g *Graph) AddTypedConnection(fromID, toID int, conn Connection) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, fromOk := g.users[fromID]
	_, toOk := g.users[toID]
	if !fromOk || !toOk || fromID == toID {
		return false
	}

	if _, ok := g.connections[fromID]; !ok {
		g.connections[fromID] = make(map[int]Connection)
	}

	g.connections[fromID][toID] = conn

	return true
}

func (g *Graph) GetConnectionsByType(userID int, connType string) []*User {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.connections[userID]) == 0 {
		return nil
	}

	cUsers := make([]*User, 0, len(g.connections[userID]))

	for id, connection := range g.connections[userID] {
		if connection.Type() == connType {
			cUsers = append(cUsers, g.users[id])
		}
	}

	return cUsers
}

func (g *Graph) GetConnectionInfo(fromID, toID int) (Connection, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn, ok := g.connections[fromID][toID]
	return conn, ok
}

func main() {
	graph := NewGraph()

	graph.AddUser(1, "Alice")
	graph.AddUser(2, "Bob")
	graph.AddUser(3, "Charlie")

	graph.AddConnection(1, 2) // Alice -> Bob
	//graph.AddConnection(1, 3) // Alice -> Charlie
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

	suggest := graph.SuggestConnections(1)
	for _, user := range suggest {
		fmt.Printf("user: %s\n", user.Name)
	}

	conn, ok := graph.GetConnectionInfo(1, 2)
	fmt.Printf("OK: %v, CONN: %v\n", ok, conn.Weight())
}
