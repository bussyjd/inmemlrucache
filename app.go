package main

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/* The structure of the LRU cache
**** We need a doubly linked list to be able to promote and the entries in memory with best performance
**** http://bigocheatsheet.com/
 */

const imgfilesizelimit = 1000000
const lrusizelimit = 10
const idsize = 5

type LRUCache struct {
	size int
	l    *list.List
}

var lru *LRUCache

// New Cache creation
func initcache(size int) *LRUCache {
	TmpfsInit()
	return &LRUCache{
		size: size,
		l:    list.New(),
	}
}

func main() {
	lru = initcache(lrusizelimit)
	router := mux.NewRouter()
	// Set(Value) Set a new item in the cache, Returns id
	router.HandleFunc("/set", SetHandler)
	// Get(id) Return the value associated with an id if it exists Otherwise returns 404
	router.HandleFunc("/get/{id}", GetHandler)
	// Delete an item of the cache
	router.HandleFunc("/del/{id}", DeleteHandler)
	// Reset() Delete all the items of the cache
	router.HandleFunc("/reset", ResetHandler)
	//	// Count() Returns the item count of the cache
	router.HandleFunc("/count", CountHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// -- ERROR HANDLING --

// SETCACHE
// Add a new entry in the lru cache and write the picture in
// tmpfs
func SetCache(lru *LRUCache, buf []byte) (string, error) {
	newfilename := randStr(idsize)
	lrulen := lru.l.Len()
	if len(buf) == 0 {
		return "", fmt.Errorf("Image size is empty\n")
	} else if len(buf) != idsize {
		return "", fmt.Errorf("Invalid id\n")
	}
	switch {
	case lrulen == lru.size:
		oldest := lru.l.Back()
		lru.l.Remove(oldest)
		lru.SetLru(newfilename)
	case lrulen < lru.size:
		lru.SetLru(newfilename)
	}
	TmpfsWrite(buf, newfilename)
	return newfilename, nil
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
func GetCache(lru *LRUCache, key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("Image size is empty\n")
	} else if len(key) != idsize {
		return nil, fmt.Errorf("Invalid id\n")
	}
	if lru.l.Len() == 0 {
		return nil, fmt.Errorf("LRU is empty\n")
	}
	validid := lru.GetLru(key)
	if validid == false {
		return nil, fmt.Errorf("Empty LRU entry\n")
	}
	data, err := TmpfsRead(key)
	if err != nil {
		return nil, err
	}
	return data, err
}

// LRUGET Check for the id in the double linked list and promotes it on success
func (lru *LRUCache) GetLru(key string) bool {
	i := 0
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		if e.Value == key {
			lru.l.MoveToFront(e)
			return true
		}
	}
	return false
}

// DEL
func RmCache(lru *LRUCache, key string) (bool, error) {
	if lru.l.Len() == 0 {
		return true, fmt.Errorf("LRU is empty\n")
	}
	if len(key) == 0 {
		return false, fmt.Errorf("Image size is empty\n")
	} else if len(key) != idsize {
		return false, fmt.Errorf("Invalid id\n")
	}
	filename := lru.RmLru(key)
	if filename == false {
		fmt.Println(filename)
		return filename, fmt.Errorf("LRU entry non existing\n")
	}
	rm, err := TmpfsRm(key)
	return rm, err
}

// LRUDEL
func (lru *LRUCache) RmLru(key string) bool {
	i := 0
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		if key == e.Value {
			lru.l.Remove(e)
			return true
		}
	}
	return false
}

// COUNT
func EntryCount(lru *LRUCache) int {
	return lru.l.Len()
}

// For debugging
func DescribeLRU(lru *LRUCache) {
	i := 0
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		fmt.Printf("Key:  %v\n", i)
		fmt.Printf("Value:  %v\n", e.Value)
	}
}

// RESET
func ResetCache(lru *LRUCache) (bool, error) {
	lru.l.Init()
	rm, err := TmpfsClear()
	fmt.Printf("FileSystem flushed\n")
	return rm, err
}

//UUID GENERATION
func randStr(strSize int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
