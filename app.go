package main

import (
	"code.google.com/p/go-uuid/uuid"
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

const imgfilesizelimit = 1000000
const lrusizelimit = 10

type LRUCache struct {
	size int
	l    *list.List
}

// New Cache creation
func initcache(size int) *LRUCache {
	TmpfsInit()
	return &LRUCache{
		size: size,
		l:    list.New(),
	}

}

func main() {
	lru := initcache(lrusizelimit)
	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	// Set(Key,Value) Set a new item in the cache
	router.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		mbreader := io.LimitReader(r.Body, imgfilesizelimit)
		buf, err := ioutil.ReadAll(mbreader)
		if err != nil {
			fmt.Println(err)
		}
		SetCache(lru, buf)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	// Get(Key) Return the value associated with a key if it exists Otherwise returns 404
	router.HandleFunc("/get/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		if err != nil {
			// invalid string
		}
		GetCache(lru, key)
	})
	// Delete an item of the cache
	router.HandleFunc("/del/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["id"])
		if err != nil {
			// invalid string
		}
		RmCache(lru, key)
	})
	// Reset() Delete all the items of the cache
	router.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		ResetCache(lru)
	})
	//	// Count() Returns the item count of the cache
	router.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		ItemCount(lru)
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

// SETCACHE
func SetCache(lru *LRUCache, buf []byte) {
	// TODO Avoid the file to be zero (check r.Body size)
	// TODO Avoid the file to be over 1Mb
	//imgrdsize, err := mbreader.Read(buf)
	//if err == io.EOF {
	//	fmt.Println(err)
	//	fmt.Fprintf(w, "Image size exeeded %v \n", imgrdsize)
	//	return
	//}
	newfilename := uuid.New() + ".jpg"
	lrulen := lru.l.Len()
	switch {
	case lrulen == lru.size:
		oldest := lru.l.Back()
		// tmpfsDel(oldest.value)
		lru.l.Remove(oldest)
		lru.SetLru(newfilename)
	case lrulen < lru.size:
		lru.SetLru(newfilename)
	}
	TmpfsWrite(buf, newfilename)
}

// LRUSET
func (lru *LRUCache) SetLru(newfilename string) {
	if lru.l.Len() >= 10 {
		fmt.Printf("Cache is full deleting oldest entry\n")
		oldest := lru.l.Back()
		lru.l.Remove(oldest)
	}
	lru.l.PushFront(newfilename)
}

// GET
func GetCache(lru *LRUCache, key int) {
	if key > lru.size {
		fmt.Printf("Cache is limited to 10 entries")
	}
	readkeyfile := lru.GetLru(key)
	TmpfsRead(readkeyfile)
}

// LRUGET
func (lru *LRUCache) GetLru(key int) string {
	i := 0
	var getfile string
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		if key == lru.size {
			break
		}
		if key == i {
			// Return the image
			// TODO tmpfsRead
			fmt.Printf("Image: %v\n", e.Value)
			fmt.Printf("Key: %v\n", key)
			// Promote the accessed entry
			fmt.Printf("Moving %v to Front", key)
			lru.l.MoveToFront(e)
			getfile = fmt.Sprintf("%v", e.Value)
		}
	}
	return getfile
}

// DEL
func RmCache(lru *LRUCache, key int) {
	rmfilename := lru.RmLru(key)
	TmpfsRm(rmfilename)
}

// LRUDEL
func (lru *LRUCache) RmLru(key int) string {
	i := 0
	var imgpwd string
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		if key == i {
			imgpwd = fmt.Sprintf("%v", e.Value)
			fmt.Printf("Image path to delete: %s\n", imgpwd)

			lru.l.Remove(e)
			break
		}
	}
	return imgpwd
}

/*
func (lru *LRUCache) Read(key int) e {
	i := 0
	for e := lru.l.Front(); e != nil; e = e.Next() {
	if key == i {
	}
}
*/

// COUNT
func ItemCount(lru *LRUCache) {
	i := 0
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		fmt.Printf("Key:  %v\n", i)
		fmt.Printf("Value:  %v\n", e.Value)
	}
}

// RESET
func ResetCache(lru *LRUCache) {
	lru.l.Init()
	TmpfsClear()
	fmt.Printf("Cache flushed")
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
