package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	mp3 "github.com/mikkyang/id3-go"
)

func main() {
	var i, o string

	flag.StringVar(&i, "i", ".", "Source directory")
	flag.StringVar(&o, "o", "./recovered", "Output directory")

	flag.Parse()

	files, err := ioutil.ReadDir(i)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory %s: %s\n", i, err.Error())
		return
	}
	files = filterMP3files(files)

	if err = os.MkdirAll(o, 0566); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory %s: %s\n", o, err.Error())
		return
	}

	scanned, recovered := 0, 0
	for _, f := range files {
		infile := i + "/" + f.Name()
		scanned++
		metadata, err := mp3.Open(infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %s\n", f.Name(), err.Error())
			continue
		}

		title := metadata.Title()
		artist := metadata.Artist()
		metadata.Close()

		if title == "" && artist == "" {
			fmt.Printf("can't recover '%s's name: missing ID3 tags\n", f.Name())
			continue
		}

		recovered++

		fmt.Printf("%d. %s - %s\n", recovered, artist, title)

		/*
			outfile := o + "/" + filename(artist, title, recovered)
			if err = copyfile(infile, outfile); err != nil {
				fmt.Fprintf(os.Stderr, "error creating new file %s: %s\n", outfile, err.Error())
			}
		*/
	}

	perc := 0.0
	if scanned != 0 {
		perc = float64(recovered) / float64(scanned) * 100.0
	}
	fmt.Printf("\nscanned %d file(s), recovered %d name(s) (%.2f%%)\n", scanned, recovered, perc)
}

func filterMP3files(files []os.FileInfo) []os.FileInfo {
	r := make([]os.FileInfo, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".mp3") {
			r = append(r, f)
		}
	}
	return r
}

func filename(artist, title string, recovered int) string {
	name := ""
	if title != "" && artist != "" {
		name = fmt.Sprintf("%s - %s.mp3", artist, title)
	} else if title != "" {
		name = fmt.Sprintf("%s (%d).mp3", title, recovered)
	} else if artist != "" {
		name = fmt.Sprintf("%s (%d).mp3", artist, recovered)
	}
	return name
}

func copyfile(src, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dest, input, 0644); err != nil {
		return err
	}

	return nil
}
