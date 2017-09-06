package momento

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	export :=
		`1 January 2000
==============

00:00
Hello, millenium!
With: Joe Bloggs, John Smith
At: Home: 1 Road Drive, Country (0.00000000, -0.00000000)
Tags: New Year, Millenium
Media: MEDIA_005.mp4
Media: MEDIA_109.jpg`

	reader := strings.NewReader(export)
	result, err := Parse(reader, "/dev/null")
	if err != nil {
		t.Fatalf("Parse error. %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Moment count not equal to expected. %d %d", 1, len(result))
	}

	moment := result[0]

	// Test Basic Properties From Parse
	expectedDate := "1 January 2000"
	if moment.date != expectedDate {
		t.Errorf("Moment date not equal to expected. %v %v", moment.date, expectedDate)
	}

	expectedTime := "00:00"
	if moment.time != expectedTime {
		t.Errorf("Moment time not equal to expected. %v %v", moment.time, expectedTime)
	}

	expectedText := "Hello, millenium!"
	if moment.text != expectedText {
		t.Errorf("Moment text not equal to expected.. %v %v", moment.text, expectedText)
	}

	expectedTags := []string{"New Year", "Millenium"}
	if !reflect.DeepEqual(moment.tags, expectedTags) {
		t.Errorf("Moment tags not equal to expected. %v %v", moment.tags, expectedTags)
	}

	expectedPeople := []string{"Joe Bloggs", "John Smith"}
	if !reflect.DeepEqual(moment.people, expectedPeople) {
		t.Errorf("Moment people not equal to expected. %v %v", moment.people, expectedPeople)
	}

	expectedPlaces := []string{"Home"}
	if !reflect.DeepEqual(moment.places, expectedPlaces) {
		t.Errorf("Moment places not equal to expected. %v %v", moment.places, expectedPlaces)
	}

	expectedMedia := []string{"/dev/null/MEDIA_005.mp4", "/dev/null/MEDIA_109.jpg"}
	if !reflect.DeepEqual(moment.media, expectedMedia) {
		t.Errorf("Moment media not equal to expected. %v %v", moment.media, expectedMedia)
	}

	// Test Functions
	expectedISODate := "2000-01-1 00:00"
	if moment.ISODate() != expectedISODate {
		t.Errorf("Moment ISODate not equal to expected. %v %v", moment.ISODate(), expectedISODate)
	}

	if moment.Text() != expectedText {
		t.Errorf("Moment Text not equal to expected.. %v %v", moment.Text(), expectedText)
	}

	expectedCombinedTags := []string{"New Year", "Millenium", "Joe Bloggs", "John Smith", "Home"}
	if !reflect.DeepEqual(moment.Tags(), expectedCombinedTags) {
		t.Errorf("Moment Tags not equal to expected. %v %v", moment.Tags(), expectedCombinedTags)
	}

	actualMediaJpg := moment.Media(".jpg")
	expectedMediaJpg := []string{"/dev/null/MEDIA_109.jpg"}
	if !reflect.DeepEqual(moment.Media(".jpg"), expectedMediaJpg) {
		t.Errorf("Moment Media not equal to expected. %v %v", actualMediaJpg, expectedMediaJpg)
	}
}
