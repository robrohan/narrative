package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ardanlabs/conf"
)

var build = "develop"

type Config struct {
	Start  string `conf:"short:s,default:/*"`
	End    string `conf:"short:e,default:*/"`
	Input  string `conf:"short:i,required"`
	Output string `conf:"short:o,default:final.md"`
}

func parse(cfg Config, file io.Reader, fout io.Writer) {
	code_mode := false

	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()

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
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func run() error {
	log := log.New(os.Stdout, "NT: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

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

	{
		file, err := os.Open(cfg.Input)
		if err != nil {
			log.Fatal(err)
		}
		//fout, err := os.Create(cfg.Output)
		//if err != nil {
		//	log.Fatal(err)
		//}
		fout, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer fout.Close()

		parse(cfg, file, fout)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
