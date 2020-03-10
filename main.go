package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/google/go-github/v29/github"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v2"
	"github.com/xyproto/textoutput"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

const (
	versionString      = "emojiterm 0.3.0"
	height        uint = 16 // height
	width         uint = 32 // width
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
	if ty--; height != 0 { // no custom rows? subtract 1 for prompt spacing
		ty = int(height) + 1 // weird, but in this case is necessary to add 1 :O
	}
	if width != 0 {
		tx = int(width)
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

	linesUp := strconv.Itoa(int(height)/2 + 2)
	colsRight := strconv.Itoa(int(width) + 5)

	fmt.Println()
	fmt.Print("\033[s") // save cursor position

	fmt.Print("\033[" + linesUp + "A")   // move N lines up
	fmt.Print("\033[" + colsRight + "C") // move 30 to the right

	fmt.Print(description) // print the description

	// restore cursor position
	fmt.Print("\033[u")

	return nil
}

func usage(arg0 string) {
	fmt.Println(versionString)
	fmt.Fprintln(os.Stderr, "Usage: "+arg0+" [-l] [-a] [searchword]")
}

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "emojiterm",
		Usage: "list and display GitHub emojis directly on the terminal",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
			&cli.BoolFlag{Name: "long", Aliases: []string{"l"}},
			&cli.BoolFlag{Name: "all", Aliases: []string{"a"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(versionString)
				os.Exit(0)
			}

			// Check if a searchword is given
			searchword := ""
			if c.NArg() != 0 {
				searchword = c.Args().Slice()[0]
			}

			// List all emoji names
			if c.Bool("long") {

				// Fetch emoji names and URLs
				emojis, err := fetchEmojis()
				if err != nil {
					return err
				}

				if searchword != "" {
					// A searchword was supplied

					found := false

					// List the emojis containing this string
					for name := range emojis {
						if strings.Contains(name, searchword) {
							// Highlight the search term
							o.Println(strings.Replace(name, searchword, "<red>"+searchword+"</red>", -1))
							found = true
						}
					}

					if !found {
						fmt.Fprintln(os.Stderr, "Not found: "+searchword)
						os.Exit(1)
					}

					return nil // success
				}

				// Collect and sort the names
				names := make([]string, len(emojis))
				i := 0
				for name := range emojis {
					names[i] = name
					i++
				}
				sort.Strings(names)

				// Output the names
				for _, name := range names {
					fmt.Println(name)
				}
				return nil // success
			}

			// Display all emojis, with names
			if c.Bool("all") {

				// Fetch emoji names and URLs
				emojis, err := fetchEmojis()
				if err != nil {
					return err
				}

				// Collect and sort the names
				names := make([]string, len(emojis))
				i := 0
				for name := range emojis {
					names[i] = name
					i++
				}
				sort.Strings(names)

				// Count all matching emojis
				total := 0
				if searchword != "" {
					for _, name := range names {
						if strings.Contains(name, searchword) {
							total++
						}
					}
				} else {
					total = len(names)
				}

				// Output all emojis, while waiting for a keypress between each one
				counter := 1
				for _, name := range names {
					if searchword == "" || strings.Contains(name, searchword) {
						url := emojis[name]
						display(url, ":"+name+":")
						digits := int(math.Floor(math.Log10(float64(total)) + 1)) // Calculate the number of digits in "total"
						fmt.Printf("[%"+strconv.Itoa(digits)+"d of %"+strconv.Itoa(digits)+"d] Press Enter...\n", counter, total)
						counter++
						fmt.Scanln() // Wait for Enter
					}
				}
				return nil

			}

			// We're not listing emojis, but searching for one. Check if a searchword is given.
			if searchword == "" {
				usage(os.Args[0])
				os.Exit(1)
			}

			found := false // Signals if a matching emoji is found or not

			emojis, err := fetchEmojis()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Check for exact matches first
			for name, url := range emojis {
				if name == searchword {
					if err := display(url, ":"+name+":"); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					found = true
					break
				}
			}

			// Then check for words that start with the given search word
			if !found {
				for name, url := range emojis {
					if strings.HasPrefix(name, searchword) {
						if err := display(url, ":"+name+":"); err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
						found = true
						break
					}
				}
			}

			// Then check for names that contain the given search word
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

			return nil // success
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
