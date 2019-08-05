package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type User struct {
    Name, Email string
}

type Db struct {
	v map[int]User
}

var userCount = 0
var Users Db
var mx sync.RWMutex

func initDb() Db {
	return Db{
		v: make(map[int]User),
	}
}

func parseUserFromJson(r *http.Request) (User, error) {
	var u User

	if r.Body == nil {
		return User{}, errors.New("request body is missing")
	}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return User{}, errors.New("error parsing request body")
	}
	return u, nil
}

func parseId(r *http.Request) (int, error) {
	id, ok := r.URL.Query()["id"]

    if !ok || len(id) < 1 {
        return 0, errors.New("parameter 'id' is missing")
    }

	i, _ := strconv.Atoi(id[0])
	if i < 1 || i > userCount {
		return 0, errors.New("user not found")
	}
	return i, nil
}

func printUser(w http.ResponseWriter, i int, u User) {
	_, _ = fmt.Fprintf(w, "id: %v, Name: %v, E-mail: %v\n", i, u.Name, u.Email)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		i, err := parseId(r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		mx.RLock()
		user, ok := Users.v[i]
		if !ok {
			fmt.Fprintln(w, "user not found")
			return
		}
		fmt.Fprintln(w, "created a new user!")
		printUser(w, i, user)
		mx.RUnlock()
	case "POST":
		u, err := parseUserFromJson(r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		mx.Lock()
		userCount++
		time.Sleep(1*time.Millisecond)
		Users.v[userCount] = u
		fmt.Fprintln(w, "user created!")
		printUser(w, userCount, u)
		mx.Unlock()
	case "UPDATE":
		i, err := parseId(r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		u, err := parseUserFromJson(r)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		mx.Lock()
		Users.v[i] = u
		fmt.Fprintln(w, "user updated!")
		printUser(w, i, u)
		mx.Unlock()
	case "DELETE":
		i, err := parseId(r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		u := Users.v[i]
		mx.Lock()
		delete(Users.v, i)
		fmt.Fprintln(w, "user deleted!")
		printUser(w, i, u)
		mx.Unlock()
	default:
		fmt.Fprintln(w, "unknown method")
	}

}

func main() {
	Users = initDb()
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}