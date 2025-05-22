package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/KarelKubat/flagnames"
)

var (
	flagCurrentDir = flag.Bool("current-dir", true, "when true, ensure current directory (dot) is in $PATH")
	flagPrepend    = flag.Bool("prepend", true, "when true, prepend to $PATH, otherwise append")
	flagShell      = flag.String("shell", "", "shell type: one of 'bash'`, 'zsh' or 'fish'")
)

const (
	usage = `
Usage: path-sanitizer [FLAGS] [PATHS]
Sanitizes $PATH, optionally adds the current (dot) directory, adds directories.
Emits a PATH environment setting that can be sourced, most useful in a shell startup file.
The PATHS arguments must point to directories just above bin/ or sbin/ (e.g., /usr/local).

Examples:
  path-sanitizer -s bash /opt/local  # emits export PATH=... with /opt/local/{bin,sbin} when these exist
  path-sanitizer -s bash -c          # emits export PATH=... with the current (dot) directory present

Usage:
  source <(path-sanitizer ...)  # bash
  eval "$(path-sanitzer ...)"   # zsh or fish

Flags, which may be abbreviated (e.g. '-s' for '-shell'):
`
)

func main() {
	flagnames.Patch()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	// Shell must be given.
	if *flagShell != "bash" && *flagShell != "zsh" && *flagShell != "fish" {
		log.Fatal("path-sanitizer: -shell must be one of 'bash', 'zsh' or 'fish'")
	}

	// Fetch path.
	path := os.Getenv("PATH")
	// fmt.Println("initial PATH:", path)

	// Add any bin/ or sbin/ paths under all args
	for _, arg := range flag.Args() {
		for _, sub := range []string{"bin", "sbin"} {
			target := arg + "/" + sub
			if !isdir(target) {
				continue
			}
			if *flagPrepend {
				path = target + ":" + path
			} else {
				path = path + ":" + target
			}
		}
	}

	// Remove duplicates and clean up
	parts := strings.Split(path, ":")
	// fmt.Println("initial parts:", parts)
	newparts := []string{}
	for i, part := range parts {
		if part == "" {
			continue
		}
		for j := i + 1; j < len(parts); j++ {
			if parts[i] == parts[j] {
				continue
			}
		}
		newparts = append(newparts, part)
		// fmt.Println("appended", part, "now it's", newparts)
	}

	// Emit new PATH setting
	path = strings.Join(newparts, ":")
	switch *flagShell {
	case "bash":
		fallthrough
	case "zsh":
		fmt.Printf("export PATH=\"%s\"\n", path)
	case "fish":
		fmt.Printf("set -gx PATH \"%s\"\n", path)
	default:
		panic("path-sanitizer: shelltype selection failure")
	}
}

func isdir(path string) bool {
	fi, err := os.Stat(path)
	return err != nil && fi.IsDir()
}
