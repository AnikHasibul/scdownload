# scdownload
[![Go Report Card](https://goreportcard.com/badge/github.com/anikhasibul/scdownload)](https://goreportcard.com/report/github.com/anikhasibul/scdownload)
> A full featured soundcloud downloader, supports downloading playlists and single track. Downloads, even if the track is not downloadable!

**THIS DOWNLOADER ALSO DOWNLOADS THE SONG THAT FLAGGED AS NOT-DOWNLOADABLE BY THE UPLOADER! SO IT'S UP TO YOU TO BE ETHICAL**

# Installation

```sh
go get github.com/anikhasibul/scdownload
```

# Usage 

```
scdownload --help
```

# Example Usage:

Downloads a single track
```sh
$ scdownload -t https://soundcloud.com/rahulcat/give-me-some-sunshine-3idiots
```

Downloads a whole playlist concurrently
```sh
$ scdownload -p https://soundcloud.com/anik-thexplorer/sets/my-fav
``` 

Downloads a whole playlist one-by-one song
```sh
$ scdownload -c="false" -p https://soundcloud.com/anik-thexplorer/sets/my-fav
``` 

