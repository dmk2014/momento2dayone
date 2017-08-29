package dayone

import (
	"fmt"
	"os/exec"
	"path"
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

		args := make([]string, 0, 20)

		args = append(args, "new")
		args = append(args, m.Text())

		args = append(args, "-d")
		args = append(args, m.ISODate())

		if tags := m.Tags(); len(tags) > 0 {
			args = append(args, "-t")
			args = append(args, tags...)
		}

		// Images (DayOne2 does not support video at present)
		if media := m.Media(".jpg"); len(media) > 0 {
			args = append(args, "-p")
			for _, image := range media {
				// TODO: Should Moment expose Path or just filename?
				absPath := path.Join("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04/Attachments", image)
				args = append(args, absPath)
			}
		}

		// Ignore Standard In (default behaviour if no text argument provided)
		args = append(args, "--no-stdin")

		cmd := exec.Command("dayone2", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// TODO: We will continue, log output/error
			fmt.Println(string(output))
			fmt.Println(err)
		}
	}
}
