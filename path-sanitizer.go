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
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	// Shell must be given.
	if *flagShell != "bash" && *flagShell != "zsh" && *flagShell != "fish" {
		log.Fatal("path-sanitizer: -shell must be one of 'bash', 'zsh' or 'fish'")
	}

	// Emit new PATH setting
	fmt.Println(
		evalString(*flagShell,
			splitPath(
				extendPath(os.Getenv("PATH"), *flagCurrentDir, flag.Args(), *flagPrepend))))
}

// Extends the $PATH setting with the dotdir and other new parts, either pre- or postpending
func extendPath(path string, useDot bool, parts []string, prepend bool) string {
	// Build up the addition to $PATH from the dotdir and all parts (so part/bin and part/sbin if such exist)
	var extra string
	if useDot {
		extra = "."
	}
	for _, arg := range parts {
		for _, bindir := range []string{arg + "/bin", arg + "/sbin"} {
			if isDir(bindir) {
				extra += ":" + bindir
			}
		}
	}

	// Add the path
	if prepend {
		path = extra + ":" + path
	} else {
		path += ":" + extra
	}

	// Avoid double/triple/etc slashes and double/triple/etc/leading/trailing colons
	for strings.Contains(path, "//") {
		path = strings.Replace(path, "//", "/", -1)
	}
	for strings.Contains(path, "::") {
		path = strings.Replace(path, "::", ":", -1)
	}
	for strings.HasPrefix(path, ":") {
		path = strings.TrimPrefix(path, ":")
	}
	for strings.HasSuffix(path, ":") {
		path = strings.TrimSuffix(path, ":")
	}

	return path
}

// Split and deduplicate the parts in a $PATH setting.
func splitPath(path string) []string {
	parts := strings.Split(path, ":")
	newparts := []string{}
	hit := map[string]struct{}{
		"": {}, // avoid empty parts
	}
	for i, part := range parts {
		// avoid parts that we already had
		if _, ok := hit[part]; ok {
			continue
		}
		hit[part] = struct{}{}
		for j := i + 1; j < len(parts); j++ {
			if parts[i] == parts[j] {
				continue
			}
		}
		newparts = append(newparts, part)
	}
	return newparts
}

// Generate a string that a shell can evaluate.
func evalString(shell string, parts []string) string {
	path := strings.Join(parts, ":")
	switch shell {
	case "bash":
		fallthrough
	case "zsh":
		return fmt.Sprintf(`export PATH="%s"`, path)
	case "fish":
		return fmt.Sprintf(`set -gx PATH "%s"`, path)
	default:
		panic("path-sanitizer: shelltype selection failure")
	}
}

// Return true when a filepath is a directory.
func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
