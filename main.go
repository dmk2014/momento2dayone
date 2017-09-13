package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/dmk2014/momento2dayone/dayone"
	"github.com/dmk2014/momento2dayone/momento"
)

func main() {
	// TODO: Retrieve arguments from command line instead of hard coding

	if err := initializeLog(); err != nil {
		log.Fatal("Logger could not be initialized.")
	}
	log.Print("Momento2DayOne Session Beginning.")

	validateRuntime()

	moments := parseMomentoExport("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04")
	expected := 6134
	if expected != len(moments) {
		log.Fatalf("Moment count mismatch. Expected: %d. Actual: %d.", expected, len(moments))
	}

	importToDayOne(moments)

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

func validateRuntime() {
	if runtime.GOOS != "darwin" {
		log.Fatal("macOS Required (this is the only platform on which the Day One CLI is available).")
	}
	if err := exec.Command("dayone2").Run(); err != nil {
		log.Fatalf("Day One CLI not found. See %q for install instructions.",
			"http://help.dayoneapp.com/day-one-2-0/command-line-interface-cli")
	}
}

func parseMomentoExport(basePath string) []momento.Moment {
	exportPath := path.Join(basePath, "Export.txt")
	mediaPath := path.Join(basePath, "Attachments")

	if _, err := os.Stat(mediaPath); err != nil {
		log.Fatalf("Attachments path (%s) could not be verified.", mediaPath)
	}

	file, err := os.Open(exportPath)
	if err != nil {
		log.Fatal("Momento export could not be opened. Verify path and try again.")
	}
	defer file.Close()

	moments, err := momento.Parse(file, mediaPath)
	if err != nil {
		log.Fatal(err)
	}

	return moments
}

func importToDayOne(moments []momento.Moment) {
	// https://npf.io/2014/05/intro-to-go-interfaces/
	// https://stackoverflow.com/questions/12994679/golang-slice-of-struct-slice-of-interface-it-implements
	// TODO: research pointer receivers, conversion and duplication issue when using &m
	entries := make([]dayone.DayOne, len(moments))
	for i, m := range moments {
		entries[i] = dayone.DayOne(m)
	}

	dayone.Import(entries)
}
