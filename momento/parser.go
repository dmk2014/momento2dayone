package momento

import (
	"bufio"
	"bytes"
	"io"
	"path"
	"regexp"
	"strings"
	"time"
)

// Moment is a represenation of an entry in a Momento journal.
type Moment struct {
	date   string
	time   string
	text   string
	people []string
	places []string
	tags   []string
	media  []string
}

func (m *Moment) setDate(date, time string) {
	m.date = date
	m.time = time
}

func (m *Moment) setText(text string) {
	m.text = text
}

func (m *Moment) isValid() bool {
	return m.date != "" && m.time != ""
}

// ISODate returns an ISO 8601 date (RFC3339).
func (m Moment) ISODate() string {
	// TODO: Parse could panic.
	momentoTime := m.date + " " + m.time
	t, _ := time.Parse("2 January 2006 03:04", momentoTime)
	return t.Format(time.RFC3339)
}

// Text returns the entry content.
func (m Moment) Text() string {
	return m.text
}

// Tags returns a combined slice of Tags, People and Places.
func (m Moment) Tags() []string {
	tags := make([]string, 0, len(m.tags)+len(m.people)+len(m.places))
	if len(m.tags) > 0 {
		tags = append(tags, m.tags...)
	}
	if len(m.people) > 0 {
		tags = append(tags, m.people...)
	}
	if len(m.places) > 0 {
		tags = append(tags, m.places...)
	}
	return tags
}

// Media returns a slice of all media that ends with the specified suffix.
func (m Moment) Media(suffix string) []string {
	media := make([]string, 0, len(m.media))
	if len(m.media) == 0 {
		return media
	}
	for _, m := range m.media {
		if strings.HasSuffix(m, suffix) {
			media = append(media, m)
		}
	}
	return media
}

// Regular Expressions required during Parse.
var dateRegex = regexp.MustCompile(`[0-9]{1,2}\s[a-zA-Z]{3,9}\s[0-9]{4}`)
var timeRegex = regexp.MustCompile(`[0-9]{2}:[0-9]{2}`)

var dateNextLinePrefix = "=========="

// Parse extracts any Moments from the provided io.Reader and returns
// them in a slice. The media path should be a location containing all
// media files encountered during parse. This location is not validated.
// Err will be non-nil should an error be encounterd while parsing the
// io.Reader contents.
func Parse(reader io.Reader, mediaPath string) (moments []Moment, err error) {
	m := Moment{}
	moments = make([]Moment, 0, 6200)
	currentDate := ""

	// Buffer for string concatenation
	buffer := bytes.Buffer{}

	scanner := bufio.NewScanner(discardBOM(reader))
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
			if m.isValid() {
				m.setText(strings.TrimSpace(buffer.String()))
				moments = append(moments, m)
			}

			// New Moment
			m = Moment{}
			m.setDate(currentDate, text)
			buffer.Reset()

			continue
		}

		// Extract Tags, Media, or append Text
		if found, place := extractPlace(text); found {
			m.places = append(m.places, place)
		} else if found, people := extractPeople(text); found {
			m.people = people
		} else if found, tags := extractTags(text); found {
			m.tags = tags
		} else if found, media := extractMedia(text, mediaPath); found {
			m.media = append(m.media, media)
		} else {
			buffer.WriteString(text)
			buffer.WriteString("\n")
		}
	}

	// Store the last Moment
	if m.isValid() {
		m.setText(strings.TrimSpace(buffer.String()))
		moments = append(moments, m)
	}

	if err = scanner.Err(); err != nil {
		return
	}

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

func extractPlace(text string) (found bool, place string) {
	if !strings.HasPrefix(text, "At:") {
		return
	}

	text = strings.TrimPrefix(text, "At: ")

	i := strings.Index(text, ":")
	if i == -1 {
		return true, text
	}

	return true, text[:i]
}

func extractPeople(text string) (found bool, people []string) {
	if !strings.HasPrefix(text, "With: ") {
		return
	}
	text = strings.TrimPrefix(text, "With: ")
	return true, strings.Split(text, ", ")
}

func extractTags(text string) (found bool, tags []string) {
	if !strings.HasPrefix(text, "Tags: ") {
		return
	}
	text = strings.TrimPrefix(text, "Tags: ")
	return true, strings.Split(text, ", ")
}

func extractMedia(text, mediaPath string) (found bool, media string) {
	if !strings.HasPrefix(text, "Media: ") {
		return
	}
	text = strings.TrimPrefix(text, "Media: ")
	return true, path.Join(mediaPath, text)
}
