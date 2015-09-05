package main

import (
	"bytes"
	"fmt"
	"os"
	"os/Exec"
	"testing"
)

var (
	lru *LRUCache
)

func init() {
	//lru := initcache(lrusizelimit)
}

/*
// To Track data in the LRU stack we use single byte increment
*/

// GET
// We write data in the LRU and compare it to the new file in
// the tmpfs partition
func TestSet(t *testing.T) {
	// TODO Compare with data in the FS
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

// GET
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
	_, err := GetCache(lru, 1)
	if err == nil {
		t.Errorf("Expected emtpy entry error %d", err)
	}
}

func TestGetNonExisting(t *testing.T) {
	lru := initcache(lrusizelimit)
	SetCache(lru, []byte("1"))
	_, err := GetCache(lru, 2)
	if err == nil {
		t.Errorf("Expected Non existing error  %d", err)
	}
}

func TestGetOufofBounds(t *testing.T) {
	lru := initcache(lrusizelimit)
	_, err := GetCache(lru, 11)
	if err == nil {
		t.Errorf("Expected Out of  %d", err)
	}
}

// We fill the LRU with 5 entries, access the oldest entry
// and check if it's in the front of the list
func TestGetPromotion(t *testing.T) {
	lru := initcache(lrusizelimit)
	for i := 0; i <= (lrusizelimit - 6); i++ {
		SetCache(lru, []byte{'i'})
	}
	oldnew, _ := GetCache(lru, 5)
	newone, _ := GetCache(lru, 1)
	if !bytes.Equal(oldnew, newone) {
		t.Errorf("Expected accessed old entry to be promoted %d != %d", oldnew, newone)
	}
}

// DELETE
// We look into the tmpfs for the deleted entry file
func TestRm(t *testing.T) {
	//defer CleanTmpfs()
	lru := initcache(lrusizelimit)
	SetCache(lru, []byte("deleteme"))
	filename := fmt.Sprintf("%d", lru.l.Front().Value)
	fmt.Printf(filename)
	deleted, err := RmCache(lru, 1)
	if err != nil {
		t.Errorf("Expected No Error on Delete: %d", err)
	}
	if deleted == false {
		t.Errorf("Expected Rmcache to return true: %d", deleted)
	}
	if _, err := os.Stat("/tmp/lru" + filename); err == nil {
		t.Errorf("Expected deleted entry's file to be deleted too")
	}
}

func TestRmEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	rm, err := RmCache(lru, 1)
	if rm != true {
		t.Errorf("Expected Deletion on empty LRU: %d", rm)
	}
	if err == nil {
		t.Errorf("Expected error on empty LRU: %d", err)
	}
}

func TestRmNonExisting(t *testing.T) {
	lru := initcache(lrusizelimit)
	SetCache(lru, []byte("1"))
	rm, err := RmCache(lru, 2)
	if rm != true {
		t.Errorf("Expected Deletion on non existing entry: %d", rm)
	}
	if err == nil {
		t.Errorf("Expected Error on non existing entry deletion: %d", err)
	}
}

// We remove an entry id and compare with the same entry id
func TestRmCausesPreviousPromotion(t *testing.T) {
	//defer CleanTmpfs()
	lru := initcache(lrusizelimit)
	for i := 0; i <= (lrusizelimit - 6); i++ {
		SetCache(lru, []byte{'i'})
	}
	tobepromoted, _ := GetCache(lru, 2)
	RmCache(lru, 1)
	promoted, _ := GetCache(lru, 1)
	if !bytes.Equal(promoted, tobepromoted) {
		t.Errorf("Expected previous entry of the deleted to be promoted: %d, %d", tobepromoted, promoted)
	}
}

func TestRmOutOfBounds(t *testing.T) {
	lru := initcache(lrusizelimit)
	_, err := RmCache(lru, lrusizelimit+1)
	if err == nil {
		t.Errorf("Expected Ouf of bounds error")
	}
}

func TestCount(t *testing.T) {
	lru := initcache(lrusizelimit)
	SetCache(lru, []byte("1"))
	count := EntryCount(lru)
	if count >= 0 {
	} else {
		t.Errorf("Expected count to be >= 0 %d", count)
	}
}

func TestCountEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	count := EntryCount(lru)
	if count >= 0 {
	} else {
		t.Errorf("Expected count to be >= 0 %d", count)
	}
}

func TestReset(t *testing.T) {
	//reset, err := ResetCache(lru)
	lru := initcache(lrusizelimit)
	reset, err := ResetCache(lru)
	if reset != true {
		t.Errorf("Expected ResetCache to return true %d", reset)
	}
	if err != nil {
		t.Errorf("Expected ResetCache to run with no error %d", err)
	}
}

func TestResetEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	ResetCache(lru)
}

func CleanTmpfs() {
	out, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/lru/*").Output()
	if err != nil {
		fmt.Printf(" %s\n", out)
		panic(err)
	}
}