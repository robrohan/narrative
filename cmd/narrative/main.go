/*

# Includes

Here we are importing several packages. Since this application mostly just combines
files and concatenates strings, almost all of these imports are to support those
activities.

The _ardanlabs/conf_ include is a library to help make using command line parameters a
bit easier. More can be seen on their website ...

[@DonaldKnuthLiterateProgramming_2016_WebofStories-LifeStoriesofRemarkablePeople]

*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ardanlabs/conf"
)

/*

This _build_ variable will be overwritten by our build script. This value will be
the git hash of the current, build time commit.

*/
var build = "develop"

/*

# Config Struct

This structure is used to hold the command line parameters passed when the application
was started. Note that only _input_ is required, and _input_ needs to be a line sparated
file that has a list of files to concatenate together.

*/
type Config struct {
	Start  string `conf:"short:s,default:/*"`
	End    string `conf:"short:e,default:*/"`
	Input  string `conf:"short:i,required"`
	Output string `conf:"short:o,default:final.md"`
}

/*

# Parse the NARRATIVE File

*/
func parseHeader(cfg Config, log *log.Logger) {
	// the narrative config file with the list of files we'll parse
	narrativeFile, err := os.Open(cfg.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer narrativeFile.Close()

	// output file
	fout, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

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
		log.Println(inputFile)
		parse(cfg, log, inputFile, fout)
	}
}

/*

# Parse a Markdown File

*/
func parse(cfg Config, log *log.Logger, filePath string, fout io.Writer) {
	code_mode := false
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	extension := strings.ToLower(filepath.Ext(filePath))

	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()

		if extension == ".md" || extension == ".markdown" {
			_, err := fmt.Fprintf(fout, "%s\n", line)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// line = strings.Trim(line, "\n ")
			if line == cfg.Start {
				code_mode = true
				continue
			}
			if line == cfg.End {
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

/*

# Run Wrapper

*/
func run(log *log.Logger) error {
	cfg := Config{}
	if err := conf.Parse(os.Args[1:], "NT", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("NT", &cfg)
			if err != nil {
				return err
			}
			fmt.Println(usage)
			return nil
		}
		return err
	}

	parseHeader(cfg, log)

	return nil
}

/*

# Program Main Entry

*/
func main() {
	log := log.New(os.Stdout, "NT: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
