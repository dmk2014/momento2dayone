package momento

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

// Moment is a represenation of an entry in a Momento journal.
type Moment struct {
	Date   string
	Time   string
	Text   string
	People []string
	Places []string
	Tags   []string
	Media  []string
}

func (m *Moment) appendPlace(text string) {
	text = placeRegex.FindStringSubmatch(text)[1] // submatch -> 0 is entire string, first capture group is at 1
	m.Places = append(m.Places, text)
}

func (m *Moment) setPeople(text string) {
	text = strings.Replace(text, "With: ", "", 1)
	m.People = strings.Split(text, ", ")
}

func (m *Moment) setTags(text string) {
	text = strings.Replace(text, "Tags: ", "", 1)
	m.Tags = strings.Split(text, ", ")
}

func (m *Moment) appendMedia(text string) {
	text = strings.Replace(text, "Media: ", "", 1)
	m.Media = append(m.Media, text)
}

// Regular Expressions required during Parse.
var dateRegex = regexp.MustCompile(`[0-9]{1,2}\s[a-zA-Z]{3,9}\s[0-9]{4}`)
var timeRegex = regexp.MustCompile(`[0-9]{2}:[0-9]{2}`)
var placeRegex = regexp.MustCompile(`At: ([^:]+)`)

var dateNextLinePrefix = "=========="

// Parse extracts any Moments from the file at the provided path and
// returns a slice containing all extracted Moments. Error will be
// returned if the file is invalid.
func Parse(p string) (moments []Moment, err error) {
	if path.Ext(p) != ".txt" {
		err = errors.New("file at p must be of type of .txt")
		return
	}

	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	m := Moment{}
	moments = make([]Moment, 0, 6200)
	currentDate := ""

	// Buffer to join strings. Much improved performance over naive concatenation (tested at ~90k lines, 20s - 0.1s)
	buffer := bytes.Buffer{}

	scanner := bufio.NewScanner(discardBOM(file))
	for scanner.Scan() {
		text := scanner.Text()

		// Assumes no Date/Time will exist within Moment text
		switch {
		case isDateCandidate(text):
			if !dateRegex.MatchString(text) {
				break
			}

			scanner.Scan()
			nextLine := scanner.Text()
			if !strings.HasPrefix(nextLine, dateNextLinePrefix) {
				break
			}
			currentDate = text

			continue
		case isTimeCandidate(text):
			if !timeRegex.MatchString(text) {
				break
			}

			// Store Moment
			if m.Time != "" {
				m.Text = strings.TrimSpace(buffer.String())
				moments = append(moments, m)
			}

			// New Moment
			m = Moment{}
			m.Date = currentDate
			m.Time = text
			buffer.Reset()

			continue
		}

		if isPlace(text) {
			m.appendPlace(text)
		} else if isPeople(text) {
			m.setPeople(text)
		} else if isTags(text) {
			m.setTags(text)
		} else if isMedia(text) {
			m.appendMedia(text)
		} else {
			buffer.WriteString(text)
			buffer.WriteString("\n")
		}
	}

	// Store the last Moment
	m.Text = strings.TrimSpace(buffer.String())
	moments = append(moments, m)

	return
}

func discardBOM(reader io.Reader) io.Reader {
	const (
		bom0 = 0xEF
		bom1 = 0xBB
		bom2 = 0xBF
	)

	buffer := bufio.NewReader(reader)

	if b, err := buffer.Peek(3); err == nil {
		if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
			buffer.Discard(3)
		}
	}

	return buffer
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
