/*

# Narrative Package

The narrative package holds most of the code.

## Includes

*/
package narrative

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*

## Config Struct

This structure is used to hold the command line parameters passed when the application
was started. Note that only _input_ is required, and _input_ needs to be a line separated
file that has a list of files to concatenate together.

*/
type Config struct {
	Input   string `conf:"short:i,default:NARRATIVE"`
	Output  string `conf:"short:o,default:final.md"`
	Markers string `conf:"short:m,default:./narrative.yaml"`
}

/*

## Comment Markers

In order to handle comments in different file types, we allow for different _Comment
Markers_. A comment marker is a way to define an area we will use to look for markdown
text. Meaning, it becomes the human readable part of the application.

*/
type CommentMarkers struct {
	Markers []Marker `yaml:"Marker"`
}

/*

A _Marker_ is a single file type's markdown area definition.

*/
type Marker struct {
	Ext   []string `yaml:"Ext"`
	Start string   `yaml:"Start"`
	End   string   `yaml:"End"`
}

/*

## Parse the NARRATIVE File

The Narrative file is used to describe the parse order of the files - and also which
files to include or exclude.

The format of this file is:

* A singe file per line.
* A '#' on the start of a line to denote a single line comment.
* The files will be processed in order they appear in the file..

Some projects choose to create this file dynamically to
include all files within the project. For example, you could add something like the
following before the build step:

```
find ./ -name "*.tf" >> NARRATIVE
```

You can then feed the _NARRATIVE_ file into the narrative process.

*/
func ParseNarrative(cfg Config, log *log.Logger) {
	// open the NARRATIVE input file
	narrativeFile, err := os.Open(cfg.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer narrativeFile.Close()

	// open the output file
	fout, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	log.Printf("Marker config: %s\n", cfg.Markers)
	markers, err := ParseMarkerConfig(cfg.Markers)
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(narrativeFile.Name())
	rd := bufio.NewReader(narrativeFile)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			return
		}

		line = strings.Trim(line, "\n ")
		if line == "" || line[0] == '#' {
			continue
		}
		inputFile := fmt.Sprintf("%s%c%s", dir, filepath.Separator, line)
		log.Printf("Processing: %s\n", inputFile)

		// call the main "markdown finding code" for this file
		Parse(markers, log, inputFile, fout)
	}
}

/*

## Read Marker Config YAML

This is some boilerplate code to read in the YAML file that describes how
comment blocks should start and stop.

*/
func ParseMarkerConfig(markersFile string) (*CommentMarkers, error) {
	filename, _ := filepath.Abs(markersFile)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var markers CommentMarkers

	err = yaml.Unmarshal(yamlFile, &markers)
	if err != nil {
		return nil, err
	}

	return &markers, nil
}

/*

## Find A Marker

This function finds the maker definition based on the passed file extension. While
this function "doesn't scale", the expected amount of configuration data means it
should not be a problem.

---

Note: while looking for an extension match, the _Ext_ array is expected to be in
alpha order. If the first letter of the defined extension is before the first 
letter of the sought after extension, the loop moves on to the next section.

---

The function double loops over the passed in yaml config file looking for the
section that matches the extension. If this becomes a problem, the file could be
indexed by extension instead.

*/
func FindMarker(markers *CommentMarkers, extension string) (*Marker, error) {
	for i := range markers.Markers {
		testExt := markers.Markers[i].Ext
		for m := range testExt {
			if testExt[m][0] > extension[1] {
				break;
			}
			if string("."+testExt[m]) == extension {
				return &markers.Markers[i], nil
			}
		}
	}
	return nil, errors.New("Marker definition not found. Edit narrative.yaml.")
}

/*

## Parse a Code file and Extract the Markdown

While processing the files from the NARRATIVE file, we then look within the code
file and find areas marked as markdown. We also "invert" the rest of the file
to be markdown code blocks.

*/
func Parse(markers *CommentMarkers, log *log.Logger, filePath string, fout io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// try to get the file extension
	extension := strings.ToLower(filepath.Ext(filePath))

	// find the start and end markers for this file type
	marker, err := FindMarker(markers, extension)
	if err != nil {
		log.Fatal(err)
	}

	code_mode := false
	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()

		// if the makers are blank, take the file as is
		if marker.Start == "" && marker.End == "" {
			_, err := fmt.Fprintf(fout, "%s\n", line)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// line = strings.Trim(line, "\n ")
			if strings.Trim(line, "\n ") == marker.Start {
				code_mode = true
				continue
			}
			if strings.Trim(line, "\n ") == marker.End {
				code_mode = false
				continue
			}

			if code_mode {
				_, err := fmt.Fprintf(fout, "%s\n", line)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				_, err := fmt.Fprintf(fout, "     %s\n", line)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
