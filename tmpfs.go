package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const lrudir = "/tmp/lru"

func TmpfsInit() {
	fmt.Printf("Initializing tmpfs")
	mkdir := exec.Command("mkdir", lrudir)
	err := mkdir.Start()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Command finished with error: %v", err)
	//mount := exec.Command("mount", "-t", "tmpfs", "-o", "size=20m", "tmpfs", "/tmp/lru/")
	//errmount := mount.Start()
	//if err != nil {
	//	log.Fatal(errmount)
	//}
}

func TmpfsWrite(buf []byte, id string) {
	fdest, err := os.Create(lrudir + "/uploadedfile" + id + ".jpg")
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

// Remove specified file in the FS
func TmpfsRm() {
}

//Clears the tmpfs FS
func TmpfsClear() {
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
