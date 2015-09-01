package main

import (
	"container/list"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* The structure of the LRU cache
**** We need a doubly linked list to be able to promote and the entries in memory with best performance
**** http://bigocheatsheet.com/
 */
type LRUCache struct {
	size int
	l    *list.List
}

//type entry struct {
//	key       int
//	imagepath string
//}

func main() {
	router := mux.NewRouter()
	// INIT TMPFS

	lru := New(10)
	router.HandleFunc("/", Index)
	// Set(Key,Value) Set a new item in the cache
	router.HandleFunc("/set/{id}", SetCache(lru))
	// Get(Key) Return the value assisiated with a key if it exists Otherwise returns 404
	router.HandleFunc("/get/{id}", GetCache(lru))
	//	// Delete an iem of the cache
	//	router.HandleFunc("/del/{id}", DelCache)
	//	// Reset() Delete all the items of the cache
	//router.HandleFunc("/reset", ResetCache)
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
	}
}

// GET
func GetCache(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		if err != nil {
			// invalid string
		}

		fmt.Fprintf(w, "Size: %v, key: %v", lru.size, key)
		if key > lru.size {
			fmt.Fprintf(w, "Key is superior to cache size")
			// 404 The cache sie if of 10 entries
		}
		fmt.Printf("Lengh: %v", lru.l.Len())
		// I have to search through the Doubly-Linked list at O(n)
		i := 0
		for e := lru.l.Front(); e != nil; e = e.Next() {
			i++
			fmt.Printf("Key %v", key)
			fmt.Printf("i %v", i)
			if key == i {
				//imgpwd := e.Value
				//fmt.Printf(imgpwd)
			}
		}
		/*
			// Read the image from the tmpfs
			image, exists := tmpfs.read(imgpwd)
			if exists == false {
				return nil
			}
			// promote the item
			lru.list.MovetoFront(key)
			return image
		*/
	}
}

// SET
func SetCache(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		if err != nil {
			// invalid string
		}
		buff := []byte("Hi\n")
		fmt.Fprintf(w, "Set Cache %v", key)
		lru.Set(1, buff)
	}
}

// LRUSET
func (lru *LRUCache) Set(key int, image []byte) {
	// First we check if the LRU cache is not full
	fmt.Println(lru.l)
	lru.l.Len()
	//lru.l.PushFront("imagename.jpg")
	//tmpfs.write(image, "imagename.jpg")
	//if lru.cache == nil {
	//}
}

/*
// RESET
func ResetCache(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List flushed")
}

//func (lru *LRUCache) Add(int key, [byte] data) {
//}
*/
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
