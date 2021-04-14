package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/heckej/Elimity-internship-backend-exercise/internal"
)

var args = os.Args

var name = makeName()

func log(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, message)
}

func main() {
	if err := run(); err != nil {
		message := err.Error()
		log(message)
		if _, ok := err.(usageError); ok {
			message := fmt.Sprintf("run '%s help' for usage information", name)
			log(message)
		}
	}
}

func makeName() string {
	path := args[0]
	return filepath.Base(path)
}

func parseFlags() (flagSet, error) {
	flags := flagSet{}
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.DurationVar(&flags.interval, "interval", 10*time.Second, "")
	set.StringVar(&flags.tokenFilePath, "token-file", "", "")
	set.IntVar(&flags.minStars, "min-stars", 0, "")
	set.SetOutput(ioutil.Discard)

	args := args[2:]
	if err := set.Parse(args); err != nil {
		return flags, errors.New("got invalid flags")
	}
	return flags, nil
}

func parseInterval(flags flagSet) error {
	if flags.interval <= 0 {
		return errors.New("got invalid interval")
	}
	return nil
}

func parseTokenFile(flags flagSet) (string, error) {
	if flags.tokenFilePath != "" {
		token, err := internal.ReadTokenFromFile(flags.tokenFilePath)
		if err != nil {
			message := fmt.Sprintf("failed reading token from %v: %v", flags.tokenFilePath, err)
			return "", usageError{message: message}
		}
		return token, nil
	}
	return "", nil
}

func parseMinStars(flags flagSet) error {
	if flags.minStars < 0 {
		return errors.New("got invalid min-stars")
	}
	return nil
}

func run() error {
	if nbArgs := len(args); nbArgs < 2 {
		return usageError{message: "missing command"}
	}
	switch args[1] {
	case "help":
		const usage = `
Simple CLI for tracking public GitHub repositories.

Usage:
  %[1]s help
  %[1]s track [-interval=<interval>] [-token-file=<file>] [-min-stars=<integer>]

Commands:
  help  Show usage information
  track Track public GitHub repositories

Options:
  -interval=<interval> Repository update interval, greater than zero [default: 10s]
  -token-file=<file> Path to a file containing a GitHub token to be used for authentication, ignored if empty
  -min-stars=<integer> The minimum positive number of stars that the tracked repositories must have [default: 0]
`
		fmt.Fprintf(os.Stdout, usage, name)
		return nil

	case "track":
		flags, err := parseFlags()
		if err != nil {
			return err
		}

		if err := parseInterval(flags); err != nil {
			return err
		}

		if err := parseMinStars(flags); err != nil {
			return err
		}

		token, err := parseTokenFile(flags)
		if err != nil {
			return err
		}

		if err := internal.Track(flags.interval, token, flags.minStars); err != nil {
			return fmt.Errorf("failed tracking: %v", err)
		}
		return nil

	default:
		return usageError{message: "got invalid command"}
	}
}

type flagSet struct {
	interval      time.Duration
	tokenFilePath string
	minStars      int
}

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return e.message
}
