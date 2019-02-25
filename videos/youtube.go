package videos

import (
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
	config "github.com/micro/go-config"
	"github.com/rylio/ytdl"
	"google.golang.org/api/youtube/v3"
	//github.com/cavaliercoder/grab
)

// Video struct
type Video struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Category string `json:"category"`
	Keywords string `json:"keywords"`
	Privacy  string `json:"privacy"`
}

// Videos struct
type Videos struct {
	Videos []Video `json:"videos"`
}

// ReupYt function
func ReupYt() {
	config.LoadFile("config.json")
	var videos Videos

	config.Scan(&videos)
	if len(videos.Videos) <= 0 {
		fmt.Println("There are no videos to download and upload! Please make sure your configuration is correct!")
		return
	}

	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	wg := new(sync.WaitGroup)
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	uploadedFiles := make(chan string, len(videos.Videos))

	fmt.Println("Start download videos!")
	fmt.Println("=======================")
	for _, video := range videos.Videos {
		wg.Add(1)
		go dlYt(video, uploadedFiles, wg)
	}
	wg.Wait()

	for ulFilename := range uploadedFiles {
		f := ulFilename
		go func(fileName string) {
			fmt.Println("Deleting uploaded video file: " + fileName + "...")
			var err = os.Remove(fileName)
			if err != nil {
				return
			}
		}(f)
	}

	fmt.Println("Press ENTER to exit!")
	fmt.Scanln()
}

func dlYt(video Video, uploadedFiles chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	videoInfo, err := ytdl.GetVideoInfo(video.URL)
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %s", err.Error())
		return
	}

	formats := videoInfo.Formats
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

	fileName, err := createFileName(time.Now().Format("20060102T150405.0700")+"."+formats[0].Extension, outputFileName{
		Title:         sanitizeFileNamePart(videoInfo.Title),
		Ext:           sanitizeFileNamePart(formats[0].Extension),
		DatePublished: sanitizeFileNamePart(videoInfo.DatePublished.Format("2006-01-02")),
		Resolution:    sanitizeFileNamePart(formats[0].Resolution),
		Author:        sanitizeFileNamePart(videoInfo.Author),
		Duration:      sanitizeFileNamePart(videoInfo.Duration.String()),
	})
	if err != nil {
		err = fmt.Errorf("Unable to parse output file file name: %s", err.Error())
		return
	}

	downloadURL, err := videoInfo.GetDownloadURL(formats[0])
	if err != nil {
		err = fmt.Errorf("Unable to get download url: %s", err.Error())
		return
	}

	var out io.Writer
	var logOut io.Writer = os.Stdout
	fmt.Println("\nDownloading video [" + videoInfo.Title + "]...")

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
	time.Sleep(5 * time.Second)

	if contentSize == size && err == nil {
		go upload(fileName, video, uploadedFiles)
	}
}

func upload(filename string, video Video, uploadedFiles chan string) {
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
			Title:       video.Title,
			Description: video.Desc,
			CategoryId:  video.Category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: video.Privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(video.Keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(video.Keywords, ",")
	}

	call := service.Videos.Insert("snippet,status", upload)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}

	fmt.Println("\nUploading a video as name: [" + video.Title + "]...")
	response, err := call.Media(file).Do()
	if err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(10 * time.Second)
		fmt.Printf("\nUpload video name ["+video.Title+"] successful! Video ID: %v\n", response.Id)
		uploadedFiles <- filename
	}
}
