package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Pattern struct {
	// Whether or not we are trying to match (true) or trying to not-match (false).
	Affirmative bool
	Regexp      regexp.Regexp
}

type Patterns []Pattern

func (s *Patterns) String() string {
	return "unused"
}

func (s *Patterns) Set(value string) error {
	// If the *very first character* is '!', this is a negative RE. In that
	// position specifically, '!' is not a Go RE metacharacter. (Whew.) If it's
	// present, account for it and then remove it.
	affirmative := true
	if value[0] == '!' {
		affirmative = false
		value = value[1:]
	}

	// Case-insensitive REs:
	x, e := regexp.Compile("(?i)" + value)
	if e != nil {
		return e
	}

	*s = append(*s, Pattern{affirmative, *x})
	return nil
}

func (s Patterns) MatchContents(pathname string, info os.FileInfo) bool {
	if len(s) == 0 {
		return true
	}

	file, e := os.Open(pathname)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return false
	}
	defer file.Close()

	r := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		// Require all patterns to match:
		matched := true
		for _, p := range s {
			if p.Affirmative != p.Regexp.Match(bytes) {
				matched = false
				break
			}
		}
		if matched {
			fmt.Printf("%s:%s\n", pathname, scanner.Text())
			r = true
		}
	}
	return r
}

func (s Patterns) MatchPathname(pathname string, info os.FileInfo) bool {
	if len(s) == 0 {
		return true
	}

	for _, p := range s {
		if p.Affirmative == p.Regexp.Match([]byte(pathname)) {
			return true
		}
	}
	return false
}

func matchFileType(info os.FileInfo, types string) bool {
	if info.IsDir() && strings.Contains(types, "d") {
		return true
	}
	if !info.IsDir() && strings.Contains(types, "f") {
		return true
	}
	return "" == types
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  search -h")
	fmt.Println("  search [-c expression] [-n expression] [pathnames...]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var afterString string
	var beforeString string
	var contentPatterns Patterns
	var help bool
	var namePatterns Patterns
	var sizeString string
	var types string

	flag.StringVar(&afterString, "a", "", "Show only files after this date (YYYY[MM[DD[ HH[:MM[:SS]]]]]).")
	flag.StringVar(&beforeString, "b", "", "Show only files before this date (YYYY[MM[DD[ HH[:MM[:SS]]]]]).")
	flag.Var(&contentPatterns, "c",
		`Regular expression (case-insensitive) to match file contents. This option
can occur more than once. Use a leading '!' to invert the match.`)
	flag.BoolVar(&help, "h", false, "Print this help message")
	flag.Var(&namePatterns, "n",
		`Regular expression (case-insensitive) to match file names. This option
can occur more than once. Use a leading '!' to invert the match.`)
	flag.StringVar(&sizeString, "s", "",
		`Show only files larger than a given size. Size can have a scale factor of
K, Ki, M, Mi, G, Gi, T, Ti.`)
	flag.StringVar(&types, "t", "", "File type: any combination of f (file), d (directory), or both.")
	flag.Parse()

	if help {
		printHelp()
	}

	var after time.Time
	var before time.Time
	{
		var e error
		if "" != afterString {
			after, e = ParseDateTime(afterString)
			if e != nil {
				fmt.Fprintln(os.Stderr, e)
				printHelp()
			}
		}
		if "" != beforeString {
			before, e = ParseDateTime(beforeString)
			if e != nil {
				fmt.Fprintln(os.Stderr, e)
				printHelp()
			}
		}
	}

	var size int64
	{
		if "" != sizeString {
			var e error
			size, e = ParseSize(sizeString)
			if e != nil {
				fmt.Fprintln(os.Stderr, e)
				printHelp()
			}
		}
	}

	roots := []string{"."}
	if flag.NArg() > 0 {
		roots = flag.Args()
	}

	for _, root := range roots {
		e := filepath.Walk(root, func(pathname string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}

			if !matchFileType(info, types) {
				return nil
			}

			if "" != afterString {
				if info.ModTime().Before(after) {
					return nil
				}
			}
			if "" != beforeString {
				if info.ModTime().After(before) {
					return nil
				}
			}

			if "" != sizeString {
				if info.Size() < size {
					return nil
				}
			}

			if namePatterns.MatchPathname(pathname, info) {
				if len(contentPatterns) == 0 {
					fmt.Println(pathname)
				} else {
					contentPatterns.MatchContents(pathname, info)
				}
			}
			return nil
		})
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
		}
	}
}
