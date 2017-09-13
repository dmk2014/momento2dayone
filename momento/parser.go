package momento

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"path"
	"regexp"
	"strings"
	"time"
)

// Date buffer is used to join strings in Moment.setISODate to improve parsing performance.
var dateBuffer = bytes.Buffer{}
var momentoDateLayout = "_2 January 2006 15:04"

// Moment is a represenation of an entry in a Momento journal.
type Moment struct {
	date   time.Time
	text   string
	people []string
	places []string
	tags   []string
	media  []string
}

func (m *Moment) setDate(d, t string) (err error) {
	dateBuffer.WriteString(d)
	dateBuffer.WriteString(" ")
	dateBuffer.WriteString(t)
	defer dateBuffer.Reset()

	date, err := time.ParseInLocation(momentoDateLayout, dateBuffer.String(), time.UTC)
	m.date = date

	return
}

func (m *Moment) setText(text string) {
	m.text = text
}

func (m *Moment) isValid() bool {
	return !m.date.IsZero()
}

// ISODate returns an ISO 8601 date (RFC3339).
func (m Moment) ISODate() string {
	return m.date.Format(time.RFC3339)
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

	log.Print("Momento Parse Starting.")
	start := time.Now()

	for scanner.Scan() {
		text := scanner.Text()

		// Assumes a Date/Time will never be found within Moment text
		switch {
		case isDateCandidate(text):
			if !dateRegex.MatchString(text) {
				break
			}

			currentDate = text

			// Skip "=======" line
			scanner.Scan()

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
			if err = m.setDate(currentDate, text); err != nil {
				log.Printf("Momento Parse Failed. Unable to read date %q, time %q.", currentDate, text)
				log.Print(err)
				return
			}
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
		log.Print("Momento Parse Failed. Error encountered while scanning io.Reader.")
		log.Print(err)
		return
	}

	duration := time.Since(start)
	log.Printf("Momento Parse Complete. %d entries found. Parsed in %q.", len(moments), duration.String())

	return
}

func isDateCandidate(text string) bool {
	return len(text) >= 10 && len(text) <= 17
}

func isTimeCandidate(text string) bool {
	return len(text) == 5
}

// The below helper functions are used to extract metadata from the Momento export.
var placePrefix = "At: "
var peoplePrefix = "With: "
var tagsPrefix = "Tags: "
var mediaPrefix = "Media: "
var commaSeparator = ", "
var semicolon = ":"

func extractPlace(text string) (found bool, place string) {
	if !strings.HasPrefix(text, placePrefix) {
		return
	}

	text = strings.TrimPrefix(text, placePrefix)

	i := strings.Index(text, semicolon)
	if i == -1 {
		return true, text
	}

	return true, text[:i]
}

func extractPeople(text string) (found bool, people []string) {
	if !strings.HasPrefix(text, peoplePrefix) {
		return
	}
	text = strings.TrimPrefix(text, peoplePrefix)
	return true, strings.Split(text, commaSeparator)
}

func extractTags(text string) (found bool, tags []string) {
	if !strings.HasPrefix(text, tagsPrefix) {
		return
	}
	text = strings.TrimPrefix(text, tagsPrefix)
	return true, strings.Split(text, commaSeparator)
}

func extractMedia(text, mediaPath string) (found bool, media string) {
	if !strings.HasPrefix(text, mediaPrefix) {
		return
	}
	text = strings.TrimPrefix(text, mediaPrefix)
	return true, path.Join(mediaPath, text)
}

// discardBOM removes the BOM, if present, from the provided io.Reader.
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
