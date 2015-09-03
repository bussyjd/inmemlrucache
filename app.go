package main

import (
	"container/list"
	"fmt"
	"html"
	"io"
	"io/ioutil"
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
	TmpfsInit()
	lru := New(10)
	router.HandleFunc("/", Index)
	// Set(Key,Value) Set a new item in the cache
	router.HandleFunc("/set/{id}", SetCache(lru))
	// Get(Key) Return the value assisiated with a key if it exists Otherwise returns 404
	router.HandleFunc("/get/{id}", GetCache(lru))
	// Delete an iem of the cache
	router.HandleFunc("/del/{id}", DelCache(lru))
	// Reset() Delete all the items of the cache
	router.HandleFunc("/reset", ResetCache(lru))
	//	// Count() Returns the item count of the cache
	router.HandleFunc("/count", ItemCount(lru))
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

// SET
func SetCache(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		// Set ID cannot be > 10
		if key >= lru.size {
			fmt.Fprintf(w, "Cache size limited to 10\n")
			return
		}
		// Use io.LimitReader instead to limit data loading
		//buf, _ := ioutil.ReadAll(r.Body)
		//reader := bytes.NewReader([]byte(r.Body))
		mbreader := io.LimitReader(r.Body, 1000000)
		buf, err := ioutil.ReadAll(mbreader)
		if err != nil {
			fmt.Println(err)
		}
		// TODO Avoid the file to be zero (check r.Body size)
		// TODO Avoid the file to be over 1Mb
		//imgrdsize, err := mbreader.Read(buf)
		//if err == io.EOF {
		//	fmt.Println(err)
		//	fmt.Fprintf(w, "Image size exeeded %v \n", imgrdsize)
		//	return
		//}
		// update the doubly linked list
		if err != nil {
			// invalid string
		}
		// Retireve the image data from POST
		fmt.Fprintf(w, "Set Cache %v", key)
		// Check if the status of the cache
		switch {
		case lru.l.Len() == lru.size:
			oldest := lru.l.Back()
			// tmpfsDel(oldest.value)
			lru.l.Remove(oldest)
			lru.Set(vars["id"], buf)
		case lru.l.Len() < lru.size:
			lru.Set(vars["id"], buf)
		}

		i := 0
		for e := lru.l.Front(); e != nil; e = e.Next() {
			i++
			if key == lru.size {
				fmt.Printf("Cache is full deleting oldest entry\n")
				lru.l.Remove(e)
				break
			}
		}
	}
}

// LRUSET
func (lru *LRUCache) Set(key string, image []byte) {
	// First we check if the LRU cache is not full
	lru.l.PushFront("imagename" + key + ".jpg\n")
	fmt.Println(lru.l.Len())
	//tmpfs.write(image, "imagename.jpg")
	TmpfsWrite(image, key)
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
		// I have to search through the Doubly-Linked list at O(n)
		i := 0
		for e := lru.l.Front(); e != nil; e = e.Next() {
			if key == lru.size {
				break
			}
			if key == i {
				// Return the image
				// TODO tmpfsRead
				fmt.Fprintf(w, "Image: %v\n", e.Value)
				// Promote the accessed entry
				fmt.Fprintf(w, "Moving %v to Front", key)
				lru.l.MoveToFront(e)
			}
			i++
		}
	}
}

// DEL
func DelCache(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		if err != nil {
			// invalid string
		}
		i := 0
		for e := lru.l.Front(); e != nil; e = e.Next() {
			i++
			if key == i {
				imgpwd := e.Value
				fmt.Printf("Image path to delete: %s", imgpwd)
				// Remove the image in the tmpfs
				// TmpfsDel(e.Value)
				// Remove the image in the double-linked-list
				lru.l.Remove(e)
				break
			}
		}

		// Remove the lru entry
		//lru.Remove()
		fmt.Println(key)
	}
}

// COUNT
func ItemCount(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := 0
		for e := lru.l.Front(); e != nil; e = e.Next() {
			i++
			fmt.Printf("Key:  %v\n", i)
			fmt.Printf("Value:  %v\n", e.Value)
		}
	}
}

// RESET
func ResetCache(lru *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lru.l.Init()
		fmt.Fprintf(w, "Cache flushed")
	}
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
