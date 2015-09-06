package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SetHandler(w http.ResponseWriter, r *http.Request) {
	mbreader := io.LimitReader(r.Body, imgfilesizelimit)
	buf, err := ioutil.ReadAll(mbreader)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
	}
	_, err = SetCache(lru, buf)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	data, err := GetCache(lru, key)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	if _, err := w.Write(data); err != nil {
		log.Println("unable to write image.")
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		// invalid string
	}
	_, err = RmCache(lru, key)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ResetCache(lru)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func CountHandler(w http.ResponseWriter, r *http.Request) {
	i := EntryCount(lru)
	fmt.Fprintf(w, "%s", i)
}
