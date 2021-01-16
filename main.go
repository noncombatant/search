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

type RegexpSlice []regexp.Regexp

func (s *RegexpSlice) String() string {
	return "unused"
}

func (s *RegexpSlice) Set(value string) error {
	// Case-insensitive REs:
	x, e := regexp.Compile("(?i)" + value)
	if e != nil {
		return e
	}
	*s = append(*s, *x)
	return nil
}

func (s *RegexpSlice) MatchContents(pathname string, info os.FileInfo) bool {
	if len(*s) == 0 {
		return false
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
		for _, x := range *s {
			if !x.Match(bytes) {
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

func (s *RegexpSlice) MatchPathname(pathname string, info os.FileInfo) bool {
	if len(*s) == 0 {
		return true
	}

	for _, x := range *s {
		if x.Match([]byte(pathname)) {
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
	var after time.Time
	var beforeString string
	var before time.Time
	var contentPatterns RegexpSlice
	var help bool
	var namePatterns RegexpSlice
	//var size int64
	var types string

	flag.StringVar(&afterString, "a", "", "Show only files after this date (YYYY[MM[DD[ HH[:MM[:SS]]]]]).")
	flag.StringVar(&beforeString, "b", "", "Show only files before this date (YYYY[MM[DD[ HH[:MM[:SS]]]]]).")
	flag.Var(&contentPatterns, "c",
		`Regular expression (case-insensitive) to match file contents. This option
can occur more than once.`)
	flag.BoolVar(&help, "h", false, "Print this help message")
	flag.Var(&namePatterns, "n",
		`Regular expression (case-insensitive) to match file names. This option
can occur more than once.`)
	//flag.StringVar(&size, "s", 0, "File size")
	flag.StringVar(&types, "t", "", "File type: any combination of f (file), d (directory), or both.")
	flag.Parse()

	if help {
		printHelp()
	}

	{
		var e error
		if "" != afterString {
			after, e = ParseDateTime(afterString)
			if e != nil {
				printHelp()
			}
		}
		if "" != beforeString {
			before, e = ParseDateTime(beforeString)
			if e != nil {
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
				return err
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

			// A special case: If both `RegexpSlice`s are empty, just print out all
			// pathnames unconditionally.
			if len(namePatterns) == 0 && len(contentPatterns) == 0 {
				fmt.Println(pathname)
				return nil
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
