package momento

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	export :=
		`1 January 2000
==============

00:00
Hello, millenium!
At: Home: 1 Road Drive, Country (0.00000000, -0.000000)
Tags: New Year, Millenium
Media: MEDIA_001.jpg`

	reader := strings.NewReader(export)

	result, err := Parse(reader, "/dev/null")
	if err != nil {
		t.Fatal("parse error")
	}

	if len(result) != 1 {
		t.Fatal("moment count not equal to expected")
	}

	moment := result[0]

	if moment.date != "1 January 2000" {
		t.Fatal("moment date not equal to expected")
	}
	if moment.time != "00:00" {
		t.Fatal("moment time not equal to expected")
	}
	if moment.text != "Hello, millenium!" {
		t.Fatal("moment text not equal to expected")
	}
}
