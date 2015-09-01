package main

import (
	"container/list"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// The structure of the LRU cache
type LRUCache struct {
	size int
	l    *list.List
	// implement cache
}

type entry struct {
	key       int
	imagepath string
}

func main() {
	router := mux.NewRouter()
	//l := list.New()
	//l.Init()
	// Default index page
	router.HandleFunc("/", Index)
	// Get(Key) Return the value assisiated with a key if it exists Otherwise returns 404
	router.HandleFunc("/get/{id}", GetCache)
	// Set(Key,Value) Set a new item in the cache
	//	router.HandleFunc("/set/{cache}", SetCache)
	//	// Delete an iem of the cache
	//	router.HandleFunc("/del/{id}", DelCache)
	//	// Reset() Delete all the items of the cache
	router.HandleFunc("/reset", ResetCache)
	//	// Count() Returns the item count of the cache
	//	router.HandleFunc("/count", ItemClount)
	// Remove oldest
	//  router.HandleFunc("/del/oldest", DelOld)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// New Cache creation
func New(size int) *LRUCache {
	return &LRUCache{
		size: size,
		l:    list.New(),
		// We build a hash table
		////cache: make(map[])
	}
}

func (lru *LRUCache) Add(key int, image string) {

	//if lru.cache == nil {
	//}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func GetCache(w http.ResponseWriter, r *http.Request) {
	// promote the item
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func ResetCache(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List flushed")
}

//func (lru *LRUCache) Add(int key, [byte] data) {
//}
