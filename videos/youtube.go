package videos

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	config "github.com/micro/go-config"
	"github.com/rylio/ytdl"
	"google.golang.org/api/youtube/v3"
)

// DlYt function
func DlYt() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	// Load json config file
	config.LoadFile("config.json")

	channels := config.Get("channels")
	fmt.Printf("%v", channels)

	info, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=BzYpUt1PFFI")
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

	var fileName string
	fileName, err = createFileName("testing."+formats[0].Extension, outputFileName{
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

	var file *os.File
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		err = fmt.Errorf("Unable to open output file: %s", err.Error())
		return
	}
	defer file.Close()

	err = info.Download(formats[0], file)
	if err != nil {
		println(err)
	}
}

var (
	filename    = flag.String("filename", "", "Name of video file to upload")
	title       = flag.String("title", "Test Title", "Video title")
	description = flag.String("description", "Test Description", "Video description")
	category    = flag.String("category", "22", "Video category")
	keywords    = flag.String("keywords", "", "Comma separated list of video keywords")
	privacy     = flag.String("privacy", "unlisted", "Video privacy status")
)

// UlYt function
func UlYt() {
	flag.Parse()

	//if *filename == "" {
	//	log.Fatalf("You must provide a filename of a video file to upload")
	//}

	client := getClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       *title,
			Description: *description,
			CategoryId:  *category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: *privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(*keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(*keywords, ",")
	}

	call := service.Videos.Insert("snippet,status", upload)

	file, err := os.Open(*filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", *filename, err)
	}

	response, err := call.Media(file).Do()
	handleError(err, "")
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}

func handleError(err error, msg string) {
	fmt.Println(msg)
}
