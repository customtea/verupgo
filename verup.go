package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func clearPath(path string) string {
	// https://ascii.jp/elem/000/001/430/1430904/
	if len(path) > 1 && path[0:2] == "~/" {
		my, err := user.Current()
		if err != nil {
			panic(err)
		}
		path = my.HomeDir + path[1:]
	}
	path = os.ExpandEnv(path)
	path = filepath.ToSlash(path)
	return filepath.Clean(path)
}

func getFileNameWithoutExt(path string) string {
	// https://qiita.com/KemoKemo/items/d135ddc93e6f87008521
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s Filename\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	targetFilename := clearPath(flag.Arg(0))

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	filename := getFileNameWithoutExt(targetFilename)
	extension := filepath.Ext(targetFilename)[1:]
	dirname := filepath.Dir(targetFilename)

	pattern := dirname + "/*." + extension
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	versionNum := 0
	for _, fname := range files {
		if strings.Contains(fname, filename) && strings.Contains(fname, "_v") {
			//fmt.Println(fname)
			for _, w := range strings.Split(fname, "_v") {
				if regexp.MustCompile(`[0-9][0-9]*\.`).Match([]byte(w)) {
					vnum, err := strconv.Atoi(strings.Split(w, ".")[0])
					if err != nil {
						continue
					}
					versionNum = int(math.Max(float64(versionNum), float64(vnum)))
				}
			}
		}
	}

	newFileName := fmt.Sprintf("%s/%s_v%02d.%s", dirname, filename, versionNum+1, extension)
	newFileName = filepath.FromSlash(newFileName)
	//fmt.Println(newFileName)

	src, err := os.Open(targetFilename)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(newFileName)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s is copied %s", targetFilename, newFileName)
}
