package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	mp3 "github.com/mikkyang/id3-go"
)

type fileInfo struct {
	name      string
	parentDir string
}

func (f *fileInfo) path() string {
	return f.parentDir + "/" + f.name
}

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

	if err = os.MkdirAll(o, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory %s: %s\n", o, err.Error())
		return
	}

	scanned, recovered := 0, 0
	for _, f := range files {
		scanned++
		metadata, err := mp3.Open(f.path())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %s\n", f.name, err.Error())
			continue
		}

		title := metadata.Title()
		artist := metadata.Artist()
		metadata.Close()

		if title == "" && artist == "" {
			fmt.Printf("can't recover '%s's name: missing ID3 tags\n", f.name)
			continue
		}

		recovered++

		fmt.Printf("%d. %s - %s\n", recovered, artist, title)

		outdir := o + strings.TrimPrefix(f.parentDir, i)
		if err = os.MkdirAll(outdir, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "error creating directory %s: %s\n", outdir, err.Error())
			continue
		}

		outfile := outdir + "/" + recoveredName(artist, title, recovered)
		if err = copyFile(f.path(), outfile); err != nil {
			fmt.Fprintf(os.Stderr, "error creating new file %s: %s\n", outfile, err.Error())
		}
	}

	perc := 0.0
	if scanned != 0 {
		perc = float64(recovered) / float64(scanned) * 100.0
	}
	fmt.Printf("\nscanned %d file(s), recovered %d name(s) (%.2f%%)\n", scanned, recovered, perc)
}

func mp3list(path string) ([]fileInfo, error) {
	list := make([]fileInfo, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		if f.IsDir() {
			subfiles, err := mp3list(path + "/" + f.Name())
			if err != nil {
				return list, err
			}
			if len(subfiles) > 0 {
				list = append(list, subfiles...)
			}
			continue
		}

		if !strings.HasSuffix(f.Name(), ".mp3") {
			continue
		}

		list = append(list, fileInfo{
			name:      f.Name(),
			parentDir: path,
		})
	}
	return list, nil
}

func sanitizeFileName(n string) string {
	return strings.Replace(n, "/", "_", -1)
}

func recoveredName(artist, title string, recovered int) string {
	name := ""
	if title != "" && artist != "" {
		name = fmt.Sprintf("%s - %s.mp3", artist, title)
	} else if title != "" {
		name = fmt.Sprintf("%s (%d).mp3", title, recovered)
	} else if artist != "" {
		name = fmt.Sprintf("%s (%d).mp3", artist, recovered)
	}
	return sanitizeFileName(name)
}

func copyFile(src, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dest, input, 0644); err != nil {
		return err
	}

	return nil
}
