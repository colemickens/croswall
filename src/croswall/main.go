package main

import (
	"encoding/xml"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var downloadFlag = flag.String("downloadDir", "", "download location")
var downloadDir string

type contents struct {
	Key string
}

type listBucketResult struct {
	Name     string
	Contents []contents
}

func init() {
	flag.Parse()

	if *downloadFlag != "" {
		downloadDir = *downloadFlag
	} else if len(os.Args) > 1 {
		downloadDir = os.Args[1]
	}

	if downloadDir == "" {
		panic("You failed to specify download dir. Sorry, you must be explicit.")
	}
}

func main() {
	response, err := http.Get("https://storage.googleapis.com/chromeos-wallpaper-public/")
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var res listBucketResult
	xml.Unmarshal(contents, &res)

	for _, r := range res.contents {
		if strings.Contains(r.Key, "high_resolution") {
			download(r.Key)
		}
	}
}

func download(key string) {
	filename := filepath.Join(downloadDir, key)

	if _, err := os.Stat(filename); err == nil {
		log.Printf("Skipping %s as it exists: %s\n", key, filename)
		return
	}

	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	log.Printf("Downloading %s to %s\n", key, filename)

	resp, err := http.Get("https://storage.googleapis.com/chromeos-wallpaper-public/" + key)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
