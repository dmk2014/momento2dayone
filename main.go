package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func main() {
	// TODO
	// Retrieve arguments from command line instead of hard coding
	// Validate that file exists
	// Investiagte image validation, approach if image not found
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
	scanner := bufio.NewScanner(file)

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
					moments[momentCount] = n
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

	os.Exit(0)

	// TODO: Save Moments.
	// test := exec.Command("dayone2", "new", `Text Here`, "-t", `One\ Tag`, "-d", "2017/08/01")
	// if err = test.Run(); err != nil {
	// 	panic(err)
	// }
}

func escapeTags(tags *[]string) {
	// TODO
	// e.g. "Tag Space" -> "Tag\ Space"
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
