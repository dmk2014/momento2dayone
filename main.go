package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dmk2014/momento2dayone/dayone"
	"github.com/dmk2014/momento2dayone/momento"
)

func main() {
	// TODO
	// Retrieve arguments from command line instead of hard coding
	// Validate that file exists
	// Investigate image validation, approach if image not found
	// Check for images that exist and are not used?

	if runtime.GOOS != "darwin" {
		fmt.Printf("macOS Required...")
		os.Exit(1)
	}
	if err := exec.Command("dayone2").Run(); err != nil {
		fmt.Printf("dayone2 not found. Append more info, link to install instructions...")
		fmt.Println(err)
		os.Exit(2)
	}

	start := time.Now()
	moments, err := momento.Parse("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04/Export.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	duration := time.Since(start)
	fmt.Printf("Parse Complete (%fs)\n", duration.Seconds())

	fmt.Printf("Moments Found: %d\n", len(moments))
	expectedMoments := 6134
	if expectedMoments != len(moments) {
		// TODO
	}

	// https://npf.io/2014/05/intro-to-go-interfaces/
	// https://stackoverflow.com/questions/12994679/golang-slice-of-struct-slice-of-interface-it-implements
	// TODO: research pointer receivers, conversion and duplication issue when using &m
	entries := make([]dayone.DayOne, len(moments))
	for i, m := range moments {
		entries[i] = dayone.DayOne(m)
	}

	start = time.Now()
	dayone.Import(entries)
	duration = time.Since(start)
	fmt.Printf("Import Complete (%fs)\n", duration.Seconds())

	os.Exit(0)
}
