package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
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

	file, err := os.Open("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04/Export.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Discard BOM
	const (
		bom0 = 0xEF
		bom1 = 0xBB
		bom2 = 0xBF
	)
	reader := bufio.NewReader(file)
	if b, err := reader.Peek(3); err == nil {
		if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
			reader.Discard(3)
		}
	}

	dateRegex := regexp.MustCompile(`[0-9]{1,2}\s[a-zA-Z]{3,9}\s[0-9]{4}`)
	timeRegex := regexp.MustCompile(`[0-9]{2}:[0-9]{2}`)
	placeRegex := regexp.MustCompile(`At: ([^:]+)`)

	dateNextLinePrefix := "=========="

	type Moment struct {
		Date   string
		Time   string
		Text   string
		People []string
		Places []string
		Tags   []string
		Media  []string
	}

	var m Moment
	moments := make([]Moment, 6200)

	momentCount := 0
	currentDate := ""
	expectedMoments := 6134

	// Read file using a Scanner
	// NewScanner creates a Scanner using the default Split Function function, ScanLines.
	// This strips any EOL marker, which is of form `\r?\n`. Thus we do not need to fix the file.
	// Include a note on this in repo, linking to the Go documentation as appropriate. Line returned
	// will simply be empty string if EOL.
	scanner := bufio.NewScanner(reader)

	// Buffer to join strings. Much improved performance over naive concatenation (tested at ~90k lines, 20s - 0.1s)
	buffer := bytes.Buffer{}

	start := time.Now()

	for scanner.Scan() {
		text := scanner.Text()

		// Assumptions: no Time or Date in Moment text (that would match length followed by regex)
		if isDateCandidate(text) {
			if dateRegex.MatchString(text) {
				scanner.Scan()
				nextLine := scanner.Text()
				if strings.HasPrefix(nextLine, dateNextLinePrefix) {
					currentDate = text
					continue
				}
			}
		} else if isTimeCandidate(text) {
			if timeRegex.MatchString(text) {

				// Store Moment in Slice
				if momentCount > 0 {
					m.Text = strings.TrimSpace(buffer.String())
					n := Moment(m)
					moments[momentCount-1] = n
				}

				// New Moment
				m = Moment{}
				m.Date = currentDate
				m.Time = text
				momentCount++
				buffer.Reset()

				continue
			}
		}

		if isPlace(text) {
			text = placeRegex.FindStringSubmatch(text)[1] // submatch -> 0 is entire string, first capture group is at 1
			m.Places = append(m.Places, text)
		} else if isPeople(text) {
			text = strings.Replace(text, "With: ", "", 1)
			m.People = strings.Split(text, ", ")
		} else if isTags(text) {
			text = strings.Replace(text, "Tags: ", "", 1)
			m.Tags = strings.Split(text, ", ")
		} else if isMedia(text) {
			text = strings.Replace(text, "Media: ", "", 1)
			m.Media = append(m.Media, text)
		} else {
			buffer.WriteString(text)
			buffer.WriteString("\n")
		}
	}

	duration := time.Since(start)

	if scanner.Err() != nil {
		panic(err)
	}
	if expectedMoments != momentCount {
		// TODO
	}

	fmt.Printf("Parse Complete (%fs)\n", duration.Seconds())
	fmt.Printf("Moments Found: %d\n", momentCount)

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

func isDateCandidate(text string) bool {
	return len(text) >= 10 && len(text) <= 17
}

func isTimeCandidate(text string) bool {
	return len(text) == 5
}

func isPlace(text string) bool {
	return strings.HasPrefix(text, "At:")
}

func isPeople(text string) bool {
	return strings.HasPrefix(text, "With:")
}

func isTags(text string) bool {
	return strings.HasPrefix(text, "Tags:")
}

func isMedia(text string) bool {
	return strings.HasPrefix(text, "Media:")
}
