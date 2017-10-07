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

	for i, entry := range entries {
		fmt.Printf("\rDayOne Import Running - Entry %d of %d.", i+1, len(entries))

		args := getArgs(entry)

		cmd := exec.Command("dayone2", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errors++
			log.Printf("Entry with date %q could not be imported.", entry.ISODate())
			log.Printf(string(output))
			continue
		}

		imported++

		if (i+1)%100 == 0 {
			time.Sleep(time.Second * 10)
		}
	}

	fmt.Println()

	duration := time.Since(start)
	log.Printf("DayOne Import Complete. Imported: %d, Errors: %d. Imported in %q.", imported, errors, duration.String())
}

func getArgs(entry DayOne) []string {
	args := make([]string, 0, 20)

	args = append(args, "new")
	args = append(args, entry.Text())

	args = append(args, "--isoDate")
	args = append(args, entry.ISODate())

	args = append(args, "--time-zone")
	args = append(args, "UTC")

	if tags := entry.Tags(); len(tags) > 0 {
		args = append(args, "--tags")
		args = append(args, tags...)
	}

	// Images (DayOne2 does not support video at present)
	if media := entry.Media(".jpg"); len(media) > 0 {
		args = append(args, "--photos")
		args = append(args, media...)
	}

	// Ignore Standard In (default behaviour if no text argument provided)
	args = append(args, "--no-stdin")

	return args
}
