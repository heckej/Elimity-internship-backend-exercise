# Solution to Elimity backend internship exercise 2021

This repository contains my solution to the backend programming exercise for an internship at Elimity. See the [Elimity exercise template](https://github.com/elimity-com/backend-intern-exercise) for the assignment.

## Installation

The `gothub` CLI of this solution can be installed using `go get` with Go 1.16+:

```sh
$ go get github.com/heckej/Elimity-internship-backend-exercise/cmd/gothub
```

## Usage

```sh
$ gothub help
Simple CLI for tracking public GitHub repositories.

Usage:
  gothub help
  gothub track [-interval=<interval>] [-token-file=<file>] [-min-stars=<integer>]

Commands:
  help  Show usage information
  track Track public GitHub repositories

Options:
  -interval=<interval> Repository update interval, greater than zero [default: 10s]
  -token-file=<file> Path to a file containing a GitHub token to be used for authentication, ignored if empty
  -min-stars=<integer> The minimum positive number of stars that the tracked repositories must have [default: 0]
```

## Dependencies

* To make use of authenticated requests, the `oauth2` library is imported, as described in the manual of the [`go-github` repository](https://github.com/google/go-github#authentication).
* For pretty printing of the table with repositories, the `text/tabwriter` library is used.
