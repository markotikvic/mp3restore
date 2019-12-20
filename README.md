## mp3restore
mp3restore is an mp3 file name recovery tool

For example, if you have a failing USB flash drive and you use photorec (or similar tools)
to recover it's content you might end up with a list of files that look like this:

	f11000444667.mp3
	f11000444668.mp3
	f11000444669.mp3
	f11000444670.mp3
	...

mp3restore will try to read ID3v1 metadata and find the song's title and artist, if it succeeds in
doing so you will end up with a list of files looking more like this:

	Artist - Song title.mp3
	Artist - Song title.mp3
	Artist - Song title.mp3
	...

For files with missing artist or title tags mp3restore will produce output similar to this:

	- for files with missing artist: Title <some number>.mp3
	- for files with missing title:  Artist <some number>.mp3

mp3restore will not remove or modify the original file and it will try to replicate the
source directory tree.

### Usage
	$ mp3restore -i=<source directory> -o=<target directory>
		-i - directory containing files to be processed (default: .)
		-o - output directory (default: ./recovered)

### Download

### Build from source
You will need Go installed on your machine in ourder to build project.
	$ git clone https://github.com/markotikvic/mp3restore.git
	$ cd mp3restore
	$ go get -u update
	$ go build

### Dependencies
mp3restore relies on [bogem's ID3 tag library](https://github.com/bogem/id3v2).

### License
See [license](LICENSE).


### TODO
- Include build script for all platforms
- Include releases
