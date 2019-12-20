package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	mp3 "github.com/bogem/id3v2"
)

func main() {
	var i, o string

	flag.StringVar(&i, "i", ".", "Source directory")
	flag.StringVar(&o, "o", "./recovered", "Output directory")

	flag.Parse()

	i = strings.TrimSuffix(i, "/")
	o = strings.TrimSuffix(o, "/")

	files, err := mp3list(i)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory %s: %s\n", i, err.Error())
		return
	}

	recovered := 0
	for _, f := range files {
		tag, err := mp3.Open(f.fullPath(), mp3.Options{Parse: true})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %s\n", f.name, err.Error())
			continue
		}

		title, artist := tag.Title(), tag.Artist()
		tag.Close()

		if title == "" && artist == "" {
			//fmt.Printf("can't recover '%s's name: missing ID3 tags\n", f.name)
			continue
		}
		recovered++

		fmt.Printf("%d. %s - %s\n", recovered, artist, title)

		outdir := o + strings.TrimPrefix(f.dir, i)
		if err = os.MkdirAll(outdir, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "error creating directory %s: %s\n", outdir, err.Error())
			continue
		}

		outfile := outPath(outdir, artist, title, recovered)
		if err = copyFile(f.fullPath(), outfile); err != nil {
			fmt.Fprintf(os.Stderr, "error creating new file %s: %s\n", outfile, err.Error())
		}
	}

	scanned := len(files)
	perc := 0.0
	if scanned != 0 {
		perc = float64(recovered) / float64(scanned) * 100.0
	}
	fmt.Printf("\nscanned %d file(s), recovered %d name(s) (%.2f%%)\n", scanned, recovered, perc)
}
