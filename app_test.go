package main

import (
	"fmt"
	"testing"
)

var (
	lru *LRUCache
)

func init() {
	//lru := initcache(10)
}

func TestSet(t *testing.T) {
	lru := initcache(lrusizelimit)
	buf := []byte{'i', 'm', 'a', 'g', 'e'}
	set, _ := SetCache(lru, buf)
	if set != true {
		t.Errorf("Success  expected: %d", set)
	}
}

func TestSetEmptyLru(t *testing.T) {
	lru := initcache(lrusizelimit)
	var buf []byte
	set, _ := SetCache(lru, buf)
	if set == true {
		t.Errorf("Success no expected: %d", set)
	}
}

func TestSetUntilFull(t *testing.T) {
	lru := initcache(lrusizelimit)
	var set bool
	for i := 0; i <= (lrusizelimit + 1); i++ {
		fmt.Println(i)
		//si := strconv.Itoa(i)
		set, _ = SetCache(lru, []byte{'i'})
	}
	if set != true {
		t.Errorf("Success expected: %d", set)
	}
}

func TestSetUntilFullCausesLruDemotion(t *testing.T) {
	lru := initcache(lrusizelimit)
	var set bool
	for i := 0; i <= (lrusizelimit + 1); i++ {
		set, _ = SetCache(lru, []byte{'i'})
	}
	if set != true {
		t.Errorf("Success expected: %d", set)
	}
	if lru.l.Back().Value == "9" {
		t.Errorf("Back of the lru data sopposed to be 9: %d", lru.l.Back().Value)
	}
	if lru.l.Front().Value == "11" {
		t.Errorf("Back of the lru data sopposed to be 11: %d", lru.l.Front().Value)
	}
}

func TestGet(t *testing.T) {
	lru := initcache(lrusizelimit)
	buf := []byte{'7', '8', '3', '7'}
	SetCache(lru, buf)
	get, _ := GetCache(lru, 1)
	fmt.Println(get)
	fmt.Println(buf)
	if get == nil {
		t.Errorf("Expected set data to be same as get data %d != %d", get, buf)
	}
}

func TestGetEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	get, _ := GetCache(lru, 1)
	fmt.Println(get)
}

func TestGetOufofBounds(t *testing.T) {
}

func TestGetPromition(t *testing.T) {
}

func TestRmOutOfBounds(t *testing.T) {
}

func TestRmNotExisting(t *testing.T) {
}

func TestRm(t *testing.T) {
}

func CountEmpty(t *testing.T) {
}

func Count(t *testing.T) {
}

func ResetEmpty(t *testing.T) {
}

func Reset(t *testing.T) {
}
