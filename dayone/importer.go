package dayone

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/dmk2014/momento2dayone/momento"
)

var months = make(map[string]string)

func init() {
	months["January"] = "01"
	months["Feburary"] = "02"
	months["March"] = "03"
	months["April"] = "04"
	months["May"] = "05"
	months["June"] = "06"
	months["July"] = "07"
	months["August"] = "08"
	months["September"] = "09"
	months["October"] = "10"
	months["November"] = "11"
	months["December"] = "12"
}

// Import iterates over the provided moments and utilizes the dayone2
// CLI to add them to DayOne.
func Import(moments []momento.Moment) {
	for i := 0; i < 10; i++ {
		m := moments[i]

		args := make([]string, 0, 20)

		// Text
		args = append(args, "new")
		args = append(args, m.Text)

		// Date (yyyy-mm-dd [hh:mm])
		dateParts := strings.Split(m.Date, " ")
		day := dateParts[0]
		month := months[dateParts[1]]
		year := dateParts[2]
		date := fmt.Sprintf("%s-%s-%s %s", year, month, day, m.Time)
		args = append(args, "-d")
		args = append(args, date)

		// Tags
		if len(m.Tags) > 0 || len(m.People) > 0 || len(m.Places) > 0 {
			args = append(args, "-t")
			args = append(args, m.Tags...)
			args = append(args, m.People...)
			args = append(args, m.Places...)
		}

		// Images (DayOne2 does not support video at present)
		if len(m.Media) > 0 {
			argsSet := false
			for _, media := range m.Media {
				absPath := path.Join("/Users/darren/Desktop/Momento Export 2017-08-13 16_27_04/Attachments", media)
				if path.Ext(absPath) != ".jpg" {
					// TODO: Log media not added
					continue
				}
				if !argsSet {
					args = append(args, "-p")
					argsSet = true
				}
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
