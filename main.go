package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
)

func getHTTPClient() (*http.Client, error) {
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}

	return httpClient, nil
}

func videoDL(destFile string, target string) error {
	httpClient, err := getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := httpClient.Get(target)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("reading answer: non 200[code=%v] status code received: '%v'", resp.StatusCode, err)
		return errors.New("non 200 status code received")
	}
	err = os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(out)
	_, err = io.Copy(mw, resp.Body)
	if err != nil {
		log.Printf("download video err=%s", err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	log.Println(flag.Args())
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/Kent-Recordings", usr.HomeDir)

	arg := flag.Arg(0)

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(arg)
	log.Print(feed.Title)

	feedDir := filepath.Join(currentDir, strings.ReplaceAll(feed.Title, "/", "-"))

	for i, item := range feed.Items {
		temp := strings.ReplaceAll(fmt.Sprintf("%s.mp4", item.Title), "/", "-")
		filename := filepath.Join(feedDir, temp)
		fmt.Printf("Downloading video %d: %s...\n", i, temp)
		videoDL(filename, item.GUID)
	}
}
