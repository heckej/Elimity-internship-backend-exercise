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

// command-line flags
var interval time.Duration
var tokenFilePath string
var minStars int

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

func parseFlags() error {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.DurationVar(&interval, "interval", 10*time.Second, "")
	set.StringVar(&tokenFilePath, "token-file", "", "")
	set.IntVar(&minStars, "min-stars", 0, "")
	set.SetOutput(ioutil.Discard)

	args := args[2:]
	if err := set.Parse(args); err != nil {
		return errors.New("got invalid flags")
	}
	return nil
}

func parseInterval() (time.Duration, error) {
	if interval <= 0 {
		return 0, errors.New("got invalid interval")
	}
	return interval, nil
}

func parseTokenFile() (string, error) {
	_, err := os.Stat(tokenFilePath)
	if os.IsNotExist(err) && tokenFilePath != "" {
		return "", errors.New("got non-existing token-file")
	}
	if tokenFilePath != "" {
		token, err := internal.ReadTokenFromFile(tokenFilePath)
		if err != nil {
			message := fmt.Sprintf("failed reading token from %v: %v", tokenFilePath, err)
			return "", usageError{message: message}
		}
		return token, nil
	}
	return "", nil
}

func parseMinStars() (int, error) {
	if minStars < 0 {
		return 0, errors.New("got invalid min-stars")
	}
	return minStars, nil
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
		err := parseFlags()
		if err != nil {
			return err
		}

		interval, err := parseInterval()
		if err != nil {
			return err
		}

		token, err := parseTokenFile()
		if err != nil {
			return err
		}

		minStars, err := parseMinStars()
		if err != nil {
			return err
		}

		if err := internal.Track(interval, token, minStars); err != nil {
			return fmt.Errorf("failed tracking: %v", err)
		}
		return nil

	default:
		return usageError{message: "got invalid command"}
	}
}

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return e.message
}
