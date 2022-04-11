package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
)

var build = "develop"

type Config struct {
	Start     string `conf:"short:s,default:/*"`
	End       string `conf:"short:e,default:*/"`
	MdComment string `conf:"default:'    '"`
	Input     string `conf:"short:i,required"`
	Output    string `conf:"short:o,default:final.md"`
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

	log.Printf("%v", cfg)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
