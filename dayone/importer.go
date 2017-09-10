package dayone

import (
	"fmt"
	"os/exec"
)

// DayOne is an interface that defines the contract for an entry in
// a DayOne journal.
type DayOne interface {
	Text() string
	ISODate() string
	Tags() []string
	Media(suffix string) []string
}

// Import iterates over the provided entries and utilizes the dayone2
// CLI to add them to DayOne.
func Import(entries []DayOne) {
	for i, m := range entries {
		if i > 10 {
			break
		}

		fmt.Printf("\rDayOne Import Running (%d of %d)", i+1, len(entries))

		args := make([]string, 0, 20)

		args = append(args, "new")
		args = append(args, m.Text())

		args = append(args, "--isoDate")
		args = append(args, m.ISODate())

		args = append(args, "-z")
		args = append(args, "UTC")

		if tags := m.Tags(); len(tags) > 0 {
			args = append(args, "-t")
			args = append(args, tags...)
		}

		// Images (DayOne2 does not support video at present)
		if media := m.Media(".jpg"); len(media) > 0 {
			args = append(args, "-p")
			args = append(args, media...)
		}

		// Ignore Standard In (default behaviour if no text argument provided)
		args = append(args, "--no-stdin")

		cmd := exec.Command("dayone2", args...)
		cmd.CombinedOutput()
	}
}
