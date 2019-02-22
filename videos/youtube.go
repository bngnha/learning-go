package videos

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	colorable "github.com/mattn/go-colorable"
	isatty "github.com/mattn/go-isatty"
	config "github.com/micro/go-config"
	"github.com/rylio/ytdl"
	"google.golang.org/api/youtube/v3"
	//github.com/cavaliercoder/grab
)

type Channel struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Category string `json:"category"`
	Keywords string `json:"keywords"`
	Privacy  string `json:"privacy"`
}

type Channels struct {
	Channels []Channel `json:"channels"`
}

var (
	filename    = flag.String("filename", "", "Name of video file to upload")
	title       = flag.String("title", "Test Title", "Video title")
	description = flag.String("description", "Test Description", "Video description")
	category    = flag.String("category", "22", "Video category")
	keywords    = flag.String("keywords", "", "Comma separated list of video keywords")
	privacy     = flag.String("privacy", "unlisted", "Video privacy status")
)

// ReupYt function
func ReupYt() {
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	wg := new(sync.WaitGroup)
	var err error
	var out io.Writer
	var logOut io.Writer = os.Stdout
	if runtime.GOOS == "windows" && isatty.IsTerminal(os.Stdout.Fd()) {
		logOut = colorable.NewColorableStdout()
	}

	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	config.LoadFile("config.json")
	var channels Channels

	config.Scan(&channels)
	fmt.Println("Start download files!")
	for _, channel := range channels.Channels {
		wg.Add(1)
		info, err := ytdl.GetVideoInfo(channel.Url)
		if err != nil {
			err = fmt.Errorf("Unable to fetch video info: %s", err.Error())
			return
		}

		formats := info.Formats
		filters := []string{
			fmt.Sprintf("%s:mp4", ytdl.FormatExtensionKey),
			fmt.Sprintf("!%s:", ytdl.FormatVideoEncodingKey),
			fmt.Sprintf("!%s:", ytdl.FormatAudioEncodingKey),
			fmt.Sprint("best"),
		}

		for _, filter := range filters {
			filter, err := parseFilter(filter)
			if err == nil {
				formats = filter(formats)
			}
		}

		fileName, err := createFileName(channel.Title+"."+formats[0].Extension, outputFileName{
			Title:         sanitizeFileNamePart(info.Title),
			Ext:           sanitizeFileNamePart(formats[0].Extension),
			DatePublished: sanitizeFileNamePart(info.DatePublished.Format("2006-01-02")),
			Resolution:    sanitizeFileNamePart(formats[0].Resolution),
			Author:        sanitizeFileNamePart(info.Author),
			Duration:      sanitizeFileNamePart(info.Duration.String()),
		})
		if err != nil {
			err = fmt.Errorf("Unable to parse output file file name: %s", err.Error())
			return
		}

		downloadURL, err := info.GetDownloadURL(formats[0])
		if err != nil {
			err = fmt.Errorf("Unable to get download url: %s", err.Error())
			return
		}

		fmt.Println("Downloading " + fileName + "...")
		go func() {
			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				err = fmt.Errorf("Unable to open output file: %s", err.Error())
				return
			}
			defer file.Close()
			out = file

			var req *http.Request
			req, err = http.NewRequest("GET", downloadURL.String(), nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
				if err == nil {
					err = fmt.Errorf("Received status code %d from download url", resp.StatusCode)
				}
				err = fmt.Errorf("Unable to start download: %s", err.Error())
				return
			}
			defer resp.Body.Close()
			contentSize := resp.ContentLength
			progressBar := pb.New64(contentSize).SetUnits(pb.U_BYTES)
			progressBar.ShowTimeLeft = true
			progressBar.ShowSpeed = true
			progressBar.SetWidth(75)
			progressBar.Output = logOut
			progressBar.Start()
			defer progressBar.Finish()

			out = io.MultiWriter(out, progressBar)

			size, err := io.Copy(out, resp.Body)
			if contentSize == size {
				time.Sleep(5 * time.Second)
				go upload(fileName, channel.Title, channel.Desc, channel.Category, channel.Keywords, channel.Privacy)
			}
			defer wg.Done()
		}()
		wg.Wait()
	}
}

func upload(filename string, title string, desc string, category string, keywords string, privacy string) {
	if filename == "" {
		log.Fatalf("You must provide a filename of a video file to upload")
	}

	client := getClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: desc,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert("snippet,status", upload)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}

	response, err := call.Media(file).Do()
	handleError(err, "")
	if err != nil {
		fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
	}
}

func handleError(err error, msg string) {
	fmt.Println(err)
}
