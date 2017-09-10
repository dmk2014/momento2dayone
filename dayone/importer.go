package dayone

import (
	"fmt"
	"os"
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
func Import(entries []DayOne) (err error) {
	log, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	imported, errors := 0, 0

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
		output, err := cmd.CombinedOutput()
		if err != nil {
			writeLog(log, m, output)
			errors++
			continue
		}

		imported++
	}

	fmt.Println()

	return
}

func writeLog(file *os.File, entry DayOne, output []byte) {
	file.WriteString(time.Now().Format("2006-01-02 15:04:05") + " Entry could not be imported. (" + entry.ISODate() + ")")
	file.WriteString("\n\n")
	file.WriteString(string(output))
}
