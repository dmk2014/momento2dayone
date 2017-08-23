package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

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

	fmt.Println(moments[6132])
	os.Exit(0)

	expectedMoments := 6134
	if expectedMoments != len(moments) {
		// TODO
	}

	// Save Moments
	months := make(map[string]string)
	months["January"] = "01"
	months["Feburary"] = "02"
	months["March"] = "03"
	months["April"] = "04"
	months["May"] = "05"
	months["June"] = "06"
	months["July"] = "07"
	months["August"] = "08"
	months["September"] = "09"
	months["October"] = "10"
	months["November"] = "11"
	months["December"] = "12"

	start = time.Now()

	for i := 0; i < 1; i++ {
		d := moments[i]

		args := make([]string, 0, 10)

		// Text
		args = append(args, "new")
		args = append(args, d.Text)

		// Date (yyyy-mm-dd [hh:mm])
		dateParts := strings.Split(d.Date, " ")
		day := dateParts[0]
		month := months[dateParts[1]]
		year := dateParts[2]
		date := fmt.Sprintf("%s-%s-%s %s", year, month, day, d.Time)
		args = append(args, "-d")
		args = append(args, date)

		// Tags
		if len(d.Tags) > 0 || len(d.People) > 0 || len(d.Places) > 0 {
			args = append(args, "-t")
			args = append(args, d.Tags...)
			args = append(args, d.People...)
			args = append(args, d.Places...)
		}

		// Images (DayOne2 does not support video at present)
		if len(d.Media) > 0 {
			images := make([]string, 0, len(d.Media))

			for _, media := range d.Media {
				absPath := path.Join("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04/Attachments", media)
				if path.Ext(absPath) != ".jpg" {
					// TODO: Log media not added
					continue
				}
				images = append(images, absPath)
			}

			if len(images) > 0 {
				args = append(args, "-p")
				args = append(args, images...)
			}
		}

		// Ignore Standard In (default behaviour if no text argument provided)
		args = append(args, "--no-stdin")

		cmd := exec.Command("dayone2", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// TODO: We will continue, log output/error
			fmt.Println(string(output))
			fmt.Println(err)
		}
	}

	duration = time.Since(start)
	fmt.Printf("Import Complete (%fs)\n", duration.Seconds())

	os.Exit(0)
}
