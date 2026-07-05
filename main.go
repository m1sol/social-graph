package main

import "fmt"

type User struct {
	ID   int
	Name string
}

type Graph struct {
	// TODO: добавьте поля для хранения пользователей и связей
}

func NewGraph() *Graph {
	// TODO
	return nil
}

func (g *Graph) AddUser(id int, name string) {
	// TODO
}

func (g *Graph) GetUser(id int) (*User, bool) {
	// TODO
	return nil, false
}

func (g *Graph) AddConnection(fromID, toID int) bool {
	// TODO: вернуть false если один из пользователей не существует
	return false
}

func (g *Graph) GetConnections(userID int) []*User {
	// TODO: вернуть слайс указателей на пользователей
	return nil
}

func (g *Graph) HasConnection(fromID, toID int) bool {
	// TODO
	return false
}

func (g *Graph) UserCount() int {
	return 0
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
