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
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/robrohan/narrative/internal"
	"log"
	"os"
)

/*

This _build_ variable will be overwritten by our build script. This value will be
the git hash of the current, build time commit.

*/
var build = "develop"

/*

# Run Wrapper

*/
func run(log *log.Logger) error {
	cfg := narrative.Config{}
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

	narrative.ParseNarrative(cfg, log)

	return nil
}

/*

# Program Main Entry

You've got to start somewhere.

Here we just setup the logger and kick off the main _run_ method.

*/
func main() {
	log := log.New(os.Stdout, "NT: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
