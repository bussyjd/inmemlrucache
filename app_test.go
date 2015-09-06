package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/Exec"
	"testing"
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
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	buf := []byte{'i', 'm', 'a', 'g', 'e'}
	fmt.Println(len(buf))
	id, err := SetCache(lru, buf)
	if !(len(id) > 0) {
		t.Errorf("Expected an Id of the data entry: %d", id)
	}
	if err != nil {
		t.Errorf("Expected Set file to run with no error: %d", err)
	}
	data, err := ioutil.ReadFile("/tmp/lru/" + fmt.Sprintf("%v", lru.l.Front().Value))
	if err != nil {
		t.Errorf("Expected Set file to be Readable in tpmfs: %d", err)
	}
	if !bytes.Equal(data, buf) {
		t.Errorf("Expected Set data to be the same in tmpfs %d != %d", data, buf)
	}
}

func TestSetEmptyLru(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	var buf []byte
	id, err := SetCache(lru, buf)
	if len(id) > 0 {
		t.Errorf("NOT Expected and Id of the data entry: %d", id)
	}
	if err == nil {
		t.Errorf("Expected Set file to run with an error: %d", err)
	}
}

func TestSetUntilFull(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	var id string
	var err error
	for i := 0; i <= (lrusizelimit + 1); i++ {
		id, err = SetCache(lru, []byte{'i', 'i', 'i', 'i', 'i'})
		if err != nil {
			t.Errorf("Expected Set file to run with no error: %d", err)
		}
	}
	if !(len(id) > 0) {
		t.Errorf("Expected an Id of the data entry: %d", id)
	}
}

func TestSetUntilFullCausesLruDemotion(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	var id string
	var err error
	for i := 0; i <= (lrusizelimit + 1); i++ {
		id, err = SetCache(lru, []byte{'i', 'i', 'i', 'i', 'i'})
		if err != nil {
			t.Errorf("Expected Set file to run with no error: %d", err)
		}
	}
	if !(len(id) > 0) {
		t.Errorf("Expected an Id of the data entry: %d", id)
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
	defer CleanTmpfs()
	fdest, _ := os.Create(lrudir + "/test5")
	buf := []byte{'1', '3', '3', '7'}
	fdest.Write(buf)
	fdest.Close()
	lru.l.PushFront("test5")
	get, err := GetCache(lru, "test5")
	if err != nil {
		t.Errorf("Expected Set file to run with no error: %d", err)
	}
	if !(len(get) > 0) {
		t.Errorf("Expected Get to return some data")
	}
	if !bytes.Equal(get, buf) {
		t.Errorf("Expected Get data to be the same in tmpfs %d != %d", get, buf)
	}
}

func TestGetEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	get, err := GetCache(lru, "test5")
	if err == nil {
		t.Errorf("Expected emtpy entry error %d", err)
	}
	if len(get) != 0 {
		t.Errorf("Expected empty data %d", err)
	}
}

func TestGetNonExisting(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	fdest, _ := os.Create(lrudir + "/testE")
	buf := []byte{'i', 'm', 'a', 'g', 'e'}
	fdest.Write(buf)
	fdest.Close()
	lru.l.PushFront("testE")
	get, err := GetCache(lru, "getno")
	if err == nil {
		t.Errorf("Expected Non existing error  %d", err)
	}
	if len(get) != 0 {
		t.Errorf("Expected empty data %d", err)
	}
}

func TestGetWrongIdSize(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	fdest, _ := os.Create(lrudir + "/testE")
	buf := []byte{'i', 'm', 'a', 'g', 'e'}
	fdest.Write(buf)
	fdest.Close()
	lru.l.PushFront("testE")
	get, err := GetCache(lru, "unneccesarylongid")
	if err == nil {
		t.Errorf("Expected Out of  %d", err)
	}
	if len(get) != 0 {
		t.Errorf("Expected empty data %d", err)
	}
}

// We fill the LRU with 5 entries, access the oldest entry
// and check if it's in the front of the list
func TestGetPromotion(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	for i := 0; i <= (lrusizelimit - 6); i++ {
		SetCache(lru, []byte{'i', 'i', 'i', 'i', 'i'})
	}
	oldnew, _ := GetCache(lru, "test5")
	newone, _ := GetCache(lru, "test1")
	if !bytes.Equal(oldnew, newone) {
		t.Errorf("Expected accessed old entry to be promoted %d != %d", oldnew, newone)
	}
}

// DELETE
func TestRm(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	id, _ := SetCache(lru, []byte("delet"))
	deleted, err := RmCache(lru, id)
	if err != nil {
		t.Errorf("Expected No Error on Delete: %d", err)
	}
	if deleted == false {
		t.Errorf("Expected Rmcache to return true: %d", deleted)
	}
	if _, err := os.Stat("/tmp/lru/delet"); err == nil {
		t.Errorf("Expected deleted entry's file to be deleted too")
	}
}

func TestRmEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	rm, err := RmCache(lru, "11111")
	if rm != true {
		t.Errorf("Expected Deletion on empty LRU: %d", rm)
	}
	if err == nil {
		t.Errorf("Expected error on empty LRU: %d", err)
	}
}

func TestRmNonExisting(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	SetCache(lru, []byte("1"))
	rm, err := RmCache(lru, "11111")
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
	defer CleanTmpfs()
	for i := 0; i <= (lrusizelimit - 6); i++ {
		SetCache(lru, []byte{'i', 'i', 'i', 'i', 'i'})
	}
	tobepromoted, _ := GetCache(lru, "11111")
	RmCache(lru, "11111")
	promoted, _ := GetCache(lru, "22222")
	if !bytes.Equal(promoted, tobepromoted) {
		t.Errorf("Expected previous entry of the deleted to be promoted: %d, %d", tobepromoted, promoted)
	}
}

func TestRmWrongIdSize(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	rm, err := RmCache(lru, "unnecessarylongid")
	if err == nil {
		t.Errorf("Expected Ouf of bounds error")
	}
	if rm != true {
		t.Errorf("Expected true on RmCache")
	}
}

func TestCount(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	SetCache(lru, []byte("11111"))
	count := EntryCount(lru)
	if count >= 0 {
	} else {
		t.Errorf("Expected count to be >= 0 %d", count)
	}
}

func TestCountEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	count := EntryCount(lru)
	if count != 0 {
		t.Errorf("Expected count to be >= 0 %d", count)
	}

}

func TestReset(t *testing.T) {
	//reset, err := ResetCache(lru)
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	SetCache(lru, []byte("11111"))
	reset, err := ResetCache(lru)
	if reset != true {
		t.Errorf("Expected ResetCache to return true %d", reset)
	}
	if err != nil {
		t.Errorf("Expected ResetCache to run with no error %d", err)
	}
	dir, err := ioutil.ReadDir("/tmp/lru")
	if err != nil {
		t.Errorf("Expected %s to be readable: %d", lrudir, err)
	}
	dirlen := len(dir)
	if dirlen >= 0 {
	} else {
		t.Errorf("Expected >= 0 files in %s:  %d", lrudir, dirlen)
	}
}

func TestResetEmpty(t *testing.T) {
	lru := initcache(lrusizelimit)
	defer CleanTmpfs()
	ResetCache(lru)
	dir, err := ioutil.ReadDir("/tmp/lru")
	if err != nil {
		t.Errorf("Expected %s to be readable: %d", lrudir, err)
	}
	dirlen := len(dir)
	if dirlen != 0 {
		t.Errorf("Expected zero files in %s:  %d", lrudir, dirlen)
	}
}

func CleanTmpfs() {
	out, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/lru/*").Output()
	if err != nil {
		fmt.Printf(" %s\n", out)
		panic(err)
	}
}
