package main

import (
	"encoding/json"
	"fmt"
)

type user struct {
	uid   int
	fname string
	lname string
	age   int
	city  string
}

type userdb struct {
	users map[int]*user
}

func main() {
	var udb userdb
	udb.users = make(map[int]*user)
	u1 := &user{10, "jason", "terry", 100, "delhi"}
	u2 := &user{20, "stephen", "curry", 200, "delhi"}
	u3 := &user{30, "lebron", "james", 300, "delhi"}
	udb.users[u1.uid] = u1
	udb.users[u2.uid] = u2
	udb.users[u3.uid] = u3

	m, _ := json.MarshalIndent(udb, "", "    ")
	fmt.Println(string(m))

}
