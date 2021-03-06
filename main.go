package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yanatan16/golang-soundcloud/soundcloud"
)

var api = &soundcloud.Api{
	ClientId: "LvWovRaJZlWCHql0bISuum8Bd2KX79mb",
}
var (
	log = logger{}
	wg  = sync.WaitGroup{}
)

var (
	pluri = flag.String(
		"p",
		"",
		"The playlist url you want to download.",
	)
	turi = flag.String(
		"t",
		"",
		"The track url you want to download.",
	)
	conc = flag.Bool(
		"c",
		true,
		"Enables concurrent downloading when downloading a playlist.",
	)
	verbose = flag.Bool(
		"v",
		true,
		"Enables verbose logging.",
	)
)

func init() {
	flag.Parse()
	log.verbose = *verbose
}

func main() {
	var uri string

	if *pluri != "" {
		uri = *pluri
	}
	if *turi != "" {
		uri = *turi
	}
	if uri == "" {
		log.Fatal("Please provide a valid url. For usage please use --help flag!")
	}
	res, err := api.Resolve(uri)
	if err != nil {
		log.Fatal(err)
	}
	val := strings.Replace(
		filepath.Base(res.Path),
		".json",
		"",
		-1,
	)
	if *pluri != "" {
		PlayListDl(val)
	}
	if *turi != "" {
		TrackDl(val)
	}
	wg.Wait()
	log.Println("Saving all files...")
	time.Sleep(30 * time.Second)
}

// PlayListDl downloads a playlist
func PlayListDl(val string) {
	plid, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	pl := api.Playlist(uint64(plid))
	p, err := pl.Get(url.Values{})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range p.Tracks {
		if *conc {
			go downloadTrack(v)
		} else {
			downloadTrack(v)
		}
	}

}

// TrackDl downloads a track
func TrackDl(val string) {
	tid, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	tr := api.Track(uint64(tid))
	t, err := tr.Get(url.Values{})
	if err != nil {
		log.Fatal(err)
	}
	downloadTrack(t)

}

// downloadTrack manages the perfect way to download a track
func downloadTrack(t *soundcloud.Track) {
	wg.Add(1)
	var n int64
	defer func() {
		log.Println("[SAVED]", n/1024, "kb", t.Title)
		wg.Done()
	}()
	if t.Downloadable {
		log.Println("[DOWNLOADING]", t.Title)
		n = saveTrack(t.Title, t.DownloadUrl)
		return
	}
	if t.Streamable {
		log.Println("[DOWNLOADING](BYPASS)", t.Title)
		n = saveTrack(t.Title, t.StreamUrl)
		return
	}
	log.Println("[FAIL] Can't download", t.Title)
}

// saveTrack saves the track to a file
func saveTrack(name, uri string) int64 {
	out, err := os.Create(
		prepareName(name) + ".mp3")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer out.Close()
	resp, err := http.Get(uri + "?client_id=" + api.ClientId)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return 0
	}
	return n
}

type logger struct {
	verbose bool
}

func (l *logger) Println(p ...interface{}) {
	if l.verbose {
		fmt.Println(p...)
	}
}
func (l *logger) Fatal(p ...interface{}) {
	if l.verbose {
		fmt.Println(p...)
		os.Exit(1)
	}
	fmt.Println("An error occurred!")
	os.Exit(1)
}

func prepareName(name string) string {
	v1 := strings.Replace(
		name,
		" ",
		"_",
		-1,
	)
	v2 := strings.Replace(
		v1,
		"/",
		"-",
		-1,
	)
	v3 := strings.Replace(
		v2,
		`\`,
		"-",
		-1,
	)
	v4 := strings.Replace(
		v3,
		"*",
		"-",
		-1,
	)
	return v4
}
