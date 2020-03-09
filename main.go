package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/google/go-github/v29/github"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

const (
	versionString      = "0.1.0"
	flagRows      uint = 16 // height
	flagCols      uint = 32 // width
)

// Fetch a map of all available emojis on GitHub,
// using the GitHub API.
// * The keys are names, like "snowman".
// * The values are URLs to PNG images.
func fetchEmojis() (map[string]string, error) {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return fetchEmojisUsingToken(token)
	}
	client := github.NewClient(nil)
	m, _, err := client.ListEmojis(context.Background())
	if err != nil {
		return nil, err
	}
	return m, err
}

// Fetch a map of all available emojis on GitHub,
// using the GitHub API and a token.
func fetchEmojisUsingToken(token string) (map[string]string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	m, _, err := client.ListEmojis(ctx)
	if err != nil {
		return nil, err
	}
	return m, err
}

// This function is based on code from github.com/eliukblau/pixterm,
// which is also licensed under MPL2.
func display(url, description string) error {
	const (
		flagDither uint   = 0         // dither mode, 0, 1 or 2
		flagNoBg   bool   = false     // disable background color?
		flagScale  uint   = 0         // scale method 0, 1 or 2
		matteColor string = "#000000" // matte color
	)
	var (
		pix        *ansimage.ANSImage
		err        error
		isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))
		tx         = 80
		ty         = 24
	)

	// get terminal size
	if isTerminal {
		tx, ty, err = terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return err
		}
	}

	// use custom terminal size (if applies)
	if ty--; flagRows != 0 { // no custom rows? subtract 1 for prompt spacing
		ty = int(flagRows) + 1 // weird, but in this case is necessary to add 1 :O
	}
	if flagCols != 0 {
		tx = int(flagCols)
	}

	// get scale mode from flag
	sm := ansimage.ScaleMode(flagScale)

	// get dithering mode from flag
	dm := ansimage.DitheringMode(flagDither)

	// set image scale factor for ANSIPixel grid
	sfy, sfx := ansimage.BlockSizeY, ansimage.BlockSizeX // 8x4 --> with dithering
	if ansimage.DitheringMode(flagDither) == ansimage.NoDithering {
		sfy, sfx = 2, 1 // 2x1 --> without dithering
	}

	mc, err := colorful.Hex(matteColor) // RGB color from Hex format
	if err != nil {
		return fmt.Errorf("matte color : %s is not a hex-color: %s", matteColor, err)
	}

	// create new ANSImage from url
	if matched, _ := regexp.MatchString(`^https?://`, url); matched {
		pix, err = ansimage.NewScaledFromURL(url, sfy*ty, sfx*tx, mc, sm, dm)
	} else {
		pix, err = ansimage.NewScaledFromFile(url, sfy*ty, sfx*tx, mc, sm, dm)
	}
	if err != nil {
		return err
	}

	// draw ANSImage to terminal
	if isTerminal {
		ansimage.ClearTerminal()
	}

	pix.SetMaxProcs(runtime.NumCPU()) // maximum number of parallel goroutines!

	pix.DrawExt(false, flagNoBg)

	if !isTerminal {
		return errors.New("not a terminal")
	}

	linesUp := strconv.Itoa(int(flagRows)/2 + 2)
	colsRight := strconv.Itoa(int(flagCols) + 5)

	fmt.Println()
	fmt.Print("\033[s") // save cursor position

	fmt.Print("\033[" + linesUp + "A")   // move N lines up
	fmt.Print("\033[" + colsRight + "C") // move 30 to the right
	fmt.Print(description)               // print the description

	// restore cursor position
	fmt.Print("\033[u")

	fmt.Println()

	return nil
}

func usage(arg0 string) {
	fmt.Println(versionString)
	fmt.Fprintln(os.Stderr, "Usage: "+arg0+" [-l] [searchword]")
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Args[0])
		os.Exit(1)
	}
	arg1 := os.Args[1]

	if arg1 == "-l" {
		// List all names

		// Fetch emoji names and URLs
		emojis, err := fetchEmojis()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Collect and sort the names
		names := make([]string, 0, len(emojis))
		for name := range emojis {
			names = append(names, name)
		}
		sort.Strings(names)

		// Output the names
		for _, name := range names {
			fmt.Println(name)
		}
		return
	} else if arg1 == "--version" {
		fmt.Println(versionString)
		return
	} else if arg1 == "--help" {
		usage(os.Args[0])
		return
	} else if strings.HasPrefix(arg1, "-") {
		fmt.Fprintln(os.Stderr, "unrecognized flag: "+arg1)
		os.Exit(1)
	}

	var (
		found      bool   // Signals if a matchin emoji is found or not
		searchword = arg1 // The string to search for
	)

	emojis, err := fetchEmojis()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Does one of the emoji names start with this string?
	for name, url := range emojis {
		if strings.HasPrefix(name, searchword) {
			display(url, ":"+name+":")
			found = true
			break
		}
	}

	// Does one of the emoji names contain this string?
	if !found {
		for name, url := range emojis {
			if strings.Contains(name, searchword) {
				if err := display(url, ":"+name+":"); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Fprintln(os.Stderr, "Not found: "+searchword)
		os.Exit(1)
	}
}
