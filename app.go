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
	// Set(Key,Value) Set a new item in the cache
	router.HandleFunc("/set", SetHandler)
	// Get(Key) Return the value associated with a key if it exists Otherwise returns 404
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
func SetCache(lru *LRUCache, buf []byte) (bool, error) {
	// TODO Avoid the file to be zero (check r.Body size)
	// TODO Avoid the file to be over 1Mb
	//imgrdsize, err := mbreader.Read(buf)
	//if err == io.EOF {
	//	fmt.Println(err)
	//	fmt.Fprintf(w, "Image size exeeded %v \n", imgrdsize)
	//	return
	//}
	if len(buf) == 0 {
		return false, fmt.Errorf("Image size is empty")
	}
	newfilename := randStr(lrusizelimit)
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
	return true, nil
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
func GetCache(lru *LRUCache, key int) ([]byte, error) {
	if key > lru.size {
		fmt.Printf("LRU Cache is limited to 10 entries\n")
		return nil, fmt.Errorf("LRU Cache is limited to 10 entries\n")
	}
	if lru.l.Len() == 0 {
		return nil, fmt.Errorf("LRU is empty\n")
	}
	readkeyfile := lru.GetLru(key)
	if len(readkeyfile) == 0 {
		fmt.Printf("Empty LRU entry\n")
		return nil, fmt.Errorf("Empty LRU entry\n")
	}
	data, err := TmpfsRead(readkeyfile)
	if err != nil {
		return nil, err
	}
	return data, err
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
			lru.l.MoveToFront(e)
			getfile = fmt.Sprintf("%v", e.Value)
			return getfile
		}
	}
	return ""
}

// DEL
func RmCache(lru *LRUCache, key int) (bool, error) {
	if key > lru.size {
		fmt.Printf("LRU Cache is limited to 10 entries\n")
		return false, fmt.Errorf("LRU Cache is limited to 10 entries\n")
	}
	if lru.l.Len() == 0 {
		return true, fmt.Errorf("LRU is empty\n")
	}
	filename := lru.RmLru(key)
	if filename == "" {
		return true, fmt.Errorf("LRU entry non existing\n")
	}
	rm, err := TmpfsRm(filename)
	return rm, err
}

// LRUDEL
func (lru *LRUCache) RmLru(key int) string {
	i := 0
	var imgpwd string
	for e := lru.l.Front(); e != nil; e = e.Next() {
		i++
		if key == i {
			imgpwd = fmt.Sprintf("%v", e.Value)
			lru.l.Remove(e)
			return imgpwd
		}
	}
	return ""
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
