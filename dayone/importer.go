package dayone

import (
	"fmt"
	"log"
	"os/exec"
	"time"
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
	imported, errors := 0, 0

	log.Printf("DayOne Import Starting. %d Entries.", len(entries))
	start := time.Now()

	for i, m := range entries {
		fmt.Printf("\rDayOne Import Running - Entry %d of %d.", i+1, len(entries))

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
		output, err := cmd.CombinedOutput()
		if err != nil {
			errors++
			log.Printf("Entry with date %q could not be imported.", m.ISODate())
			log.Printf(string(output))
			continue
		}

		imported++
	}

	fmt.Println()

	duration := time.Since(start)
	log.Printf("DayOne Import Complete. Imported: %d, Errors: %d. Imported in %q.", imported, errors, duration.String())
}
