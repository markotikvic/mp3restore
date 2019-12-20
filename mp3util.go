package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type mp3File struct {
	name string
	dir  string
}

func (f *mp3File) fullPath() string {
	return f.dir + "/" + f.name
}

func mp3list(path string) ([]mp3File, error) {
	list := make([]mp3File, 0)

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

		list = append(list, mp3File{
			name: f.Name(),
			dir:  path,
		})
	}
	return list, nil
}

func sanitizeFileName(n string) string {
	out := []byte(n)
	for i := range out {
		if out[i] == '/' || out[i] == '\\' {
			out[i] = '_'
		}
	}
	return string(out)
}

func outPath(outdir, artist, title string, recovered int) string {
	name := ""
	if title != "" && artist != "" {
		name = fmt.Sprintf("%s - %s.mp3", artist, title)
	} else if title != "" {
		name = fmt.Sprintf("%s (%d).mp3", title, recovered)
	} else if artist != "" {
		name = fmt.Sprintf("%s (%d).mp3", artist, recovered)
	}
	return outdir + "/" + sanitizeFileName(name)
}

func copyFile(src, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read error: %s", err.Error())
	}

	if err = ioutil.WriteFile(dest, input, 0644); err != nil {
		return fmt.Errorf("write error: %s", err.Error())
	}

	return nil
}
