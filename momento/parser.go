package momento

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
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

func (m *Moment) appendPlace(text string) {
	text = placeRegex.FindStringSubmatch(text)[1]
	m.places = append(m.places, text)
}

func (m *Moment) setPeople(text string) {
	text = strings.Replace(text, "With: ", "", 1)
	m.people = strings.Split(text, ", ")
}

func (m *Moment) setTags(text string) {
	text = strings.Replace(text, "Tags: ", "", 1)
	m.tags = strings.Split(text, ", ")
}

func (m *Moment) appendMedia(mediaPath, text string) {
	text = strings.Replace(text, "Media: ", "", 1)
	location := path.Join(mediaPath, text)
	m.media = append(m.media, location)
}

func (m *Moment) isValid() bool {
	return m.date != "" && m.time != ""
}

// ISODate returns an ISO 8601 date.
func (m Moment) ISODate() string {
	// TODO: Currently implements yyyy-mm-dd [hh:mm[:ss]] [AM|PM]
	// Implement as ISODate
	dateParts := strings.Split(m.date, " ")
	day := dateParts[0]
	month := months[dateParts[1]]
	year := dateParts[2]
	return fmt.Sprintf("%s-%s-%s %s", year, month, day, m.time)
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

var months = map[string]string{
	"January":   "01",
	"Feburary":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

// Regular Expressions required during Parse.
var dateRegex = regexp.MustCompile(`[0-9]{1,2}\s[a-zA-Z]{3,9}\s[0-9]{4}`)
var timeRegex = regexp.MustCompile(`[0-9]{2}:[0-9]{2}`)
var placeRegex = regexp.MustCompile(`At: ([^:]+)`)

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

	// Buffer to join strings. Much improved performance over naive concatenation (tested at ~90k lines, 20s - 0.1s)
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

		if isPlace(text) {
			m.appendPlace(text)
		} else if isPeople(text) {
			m.setPeople(text)
		} else if isTags(text) {
			m.setTags(text)
		} else if isMedia(text) {
			m.appendMedia(mediaPath, text)
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
