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

var build = "develop"

type Config struct {
	Start  string `conf:"short:s,default:/*"`
	End    string `conf:"short:e,default:*/"`
	Input  string `conf:"short:i,required"`
	Output string `conf:"short:o,default:final.md"`
}

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

func main() {
	log := log.New(os.Stdout, "NT: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
