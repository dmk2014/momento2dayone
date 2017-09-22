package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/dmk2014/momento2dayone/dayone"
	"github.com/dmk2014/momento2dayone/momento"
)

func main() {
	exportPath := flag.String("path",
		"/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04",
		"The Momento export path, containing Export.txt and attachments directory.")
	expected := flag.Int("count",
		6134,
		"The number of entries that should be parsed from the Momento export file. Negative values are ignored.")
	flag.Parse()

	if err := initializeLog(); err != nil {
		log.Fatal("Logger could not be initialized.")
	}
	log.Print("Momento2DayOne Session Beginning.")

	if !isEnvironmentValid() {
		log.Fatal("Invalid runtime environment.")
	}

	moments, err := momento.ParseFile(*exportPath)
	if err != nil {
		log.Fatal("Momento parse failed.")
	}
	if *expected < 0 && *expected != len(moments) {
		log.Fatalf("Moment count mismatch. Expected: %d. Actual: %d.", expected, len(moments))
	}

	entries := convertMomentToDayOne(moments)
	dayone.Import(entries)

	log.Print("Momento2DayOne Session Exiting Successfully.")
	os.Exit(0)
}

func initializeLog() (err error) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	log.SetOutput(file)

	return
}

func isEnvironmentValid() bool {
	if runtime.GOOS != "darwin" {
		log.Println("macOS Required (this is the only platform on which the Day One CLI is available).")
		return false
	}
	if err := exec.Command("dayone2").Run(); err != nil {
		log.Printf("Day One CLI not found. See %q for install instructions.",
			"http://help.dayoneapp.com/day-one-2-0/command-line-interface-cli")
		return false
	}
	return true
}

func convertMomentToDayOne(moments []momento.Moment) []dayone.DayOne {
	// https://npf.io/2014/05/intro-to-go-interfaces/
	// https://stackoverflow.com/questions/12994679/golang-slice-of-struct-slice-of-interface-it-implements
	// TODO: research pointer receivers, conversion and duplication issue when using &m
	entries := make([]dayone.DayOne, len(moments))
	for i, m := range moments {
		entries[i] = dayone.DayOne(m)
	}
	return entries
}
