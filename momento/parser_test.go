package momento

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	export :=
		`13 August 2002
==============

13:45
Hello, Day One!
With: Joe Bloggs, John Smith
At: Home: 1 Road Drive, Country (0.00000000, -0.00000000)
At: Work
At: No SemiColon (52.5003935973697, -9.52393457342944)
Tags: Journaling, First Entry
Media: MEDIA_005.mp4
Media: MEDIA_109.jpg`

	reader := strings.NewReader(export)
	result, err := Parse(reader, "/dev/null")
	if err != nil {
		t.Fatalf("Parse error. %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Moment count not equal to expected. %d %d", len(result), 1)
	}

	moment := result[0]

	// Test Basic Properties From Parse
	expectedText := "Hello, Day One!"
	if moment.text != expectedText {
		t.Errorf("Moment text not equal to expected.. %v %v", moment.text, expectedText)
	}

	expectedTags := []string{"Journaling", "First Entry"}
	if !reflect.DeepEqual(moment.tags, expectedTags) {
		t.Errorf("Moment tags not equal to expected. %v %v", moment.tags, expectedTags)
	}

	expectedPeople := []string{"Joe Bloggs", "John Smith"}
	if !reflect.DeepEqual(moment.people, expectedPeople) {
		t.Errorf("Moment people not equal to expected. %v %v", moment.people, expectedPeople)
	}

	expectedPlaces := []string{"Home", "Work", "No SemiColon"}
	if !reflect.DeepEqual(moment.places, expectedPlaces) {
		t.Errorf("Moment places not equal to expected. %v %v", moment.places, expectedPlaces)
	}

	expectedMedia := []string{"/dev/null/MEDIA_005.mp4", "/dev/null/MEDIA_109.jpg"}
	if !reflect.DeepEqual(moment.media, expectedMedia) {
		t.Errorf("Moment media not equal to expected. %v %v", moment.media, expectedMedia)
	}

	// Test Functions
	expectedISODate := "2002-08-13T13:45:00Z"
	if moment.ISODate() != expectedISODate {
		t.Errorf("Moment ISODate not equal to expected. %v %v", moment.ISODate(), expectedISODate)
	}

	if moment.Text() != expectedText {
		t.Errorf("Moment Text not equal to expected.. %v %v", moment.Text(), expectedText)
	}

	expectedCombinedTags := []string{"Journaling", "First Entry", "Joe Bloggs", "John Smith", "Home", "Work", "No SemiColon"}
	if !reflect.DeepEqual(moment.Tags(), expectedCombinedTags) {
		t.Errorf("Moment Tags not equal to expected. %v %v", moment.Tags(), expectedCombinedTags)
	}

	actualMediaJpg := moment.Media(".jpg")
	expectedMediaJpg := []string{"/dev/null/MEDIA_109.jpg"}
	if !reflect.DeepEqual(moment.Media(".jpg"), expectedMediaJpg) {
		t.Errorf("Moment Media not equal to expected. %v %v", actualMediaJpg, expectedMediaJpg)
	}
}

func TestEmptyMoment(t *testing.T) {
	export :=
		`13 August 2002
==============

13:45`

	reader := strings.NewReader(export)
	result, err := Parse(reader, "/dev/null")
	if err != nil {
		t.Fatalf("Parse error. %v", err)
	}

	moment := result[0]

	if moment.text != "" {
		t.Error("Moment text was not empty.")
	}

	if moment.tags != nil {
		t.Error("Moment tags not nil.")
	}

	if moment.people != nil {
		t.Error("Moment people not nil.")
	}

	if moment.places != nil {
		t.Error("Moment places not nil.")
	}

	if moment.media != nil {
		t.Error("Moment media not nil.")
	}

	actualEmptyMedia := moment.Media(".jpg")
	expectedEmptyMedia := []string{}
	if !reflect.DeepEqual(moment.Media(".jpg"), expectedEmptyMedia) {
		t.Errorf("Moment Media not equal to expected. %v %v", actualEmptyMedia, expectedEmptyMedia)
	}
}

func TestInvalidDate(t *testing.T) {
	export :=
		`13 Augusted 2002
==============

13:45`

	reader := strings.NewReader(export)
	_, err := Parse(reader, "/dev/null")
	if err == nil {
		t.Error("Parse error nil.")
	}
}

func TestDiscardBom(t *testing.T) {
	bom := []byte{0xEF, 0xBB, 0xBF}

	reader := bytes.NewReader(bom)
	discardBOM(reader)

	if reader.Len() != 0 {
		t.Errorf("Reader length not equal to expected. %d %d", reader.Len(), 0)
	}
}
