package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const lrudir = "/tmp/lru"

func TmpfsInit() {
	fmt.Printf("Initializing tmpfs...\n")
	mkdir := exec.Command("mkdir", lrudir)
	err := mkdir.Start()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Command finished with error: %v", err)
	//mount := exec.Command("mount", "-t", "tmpfs", "-o", "size=20m", "tmpfs", "/tmp/lru/")
	//errmount := mount.Start()
	//if errmount != nil {
	//	log.Fatal(errmount)
	//}
	fmt.Printf("Done\n")
}

func TmpfsWrite(buf []byte, filename string) {
	fdest, err := os.Create(lrudir + "/" + filename)
	if err != nil {
		fmt.Printf("Unable to create the file for writing. Check your write access privilege")
		panic(err)
	}
	defer fdest.Close()
	// write in the file
	wrote, err := fdest.Write(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", wrote)
}

// TODO TMPFS READ
func TmpfsRead(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile("/tmp/lru/" + filename)
	if err != nil {
		panic(err)
	}
	return data, err
}

// Remove specified file in the FS
func TmpfsRm(filename string) (bool, error) {
	fmt.Printf("Filename to delete: %s\n", filename)
	mount := exec.Command("rm", "-rf", "/tmp/lru/"+filename)
	errrm := mount.Start()
	if errrm != nil {
		panic(errrm)
	}
	return true, errrm
}

//Clears the tmpfs FS
func TmpfsClear() (bool, error) {
	out, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/lru/*").Output()
	if err != nil {
		panic(err)
		//log.Fatal(err)
		fmt.Printf(" %s\n", out)
	}
	return true, err
}

func TmpfsDestroy() {
	rm := exec.Command("rm", "-rf", lrudir)
	err := rm.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Command finished with error: %v", err)
	umount := exec.Command("umount", lrudir)
	errumount := umount.Start()
	if errumount != nil {
		log.Fatal(errumount)
	}
}
