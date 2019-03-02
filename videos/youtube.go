package videos

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	config "github.com/micro/go-config"
	"github.com/rylio/ytdl"
	"google.golang.org/api/youtube/v3"
	//github.com/cavaliercoder/grab
)

// Video struct
type Video struct {
	URL      string
	Title    string
	Desc     string
	Category string
	Keywords string
	Privacy  string
}

// Videos struct
type Videos struct {
	Videos []Video
}

// VideoState struct
type VideoState struct {
	video    Video
	fileName string
	status   string
}

// ReupYt function
func ReupYt() {
	err := config.LoadFile("config.json")
	if err != nil {
		fmt.Printf("Load config file failure! %s", err.Error())
		return
	}
	var videos Videos

	config.Scan(&videos)
	if len(videos.Videos) <= 0 {
		fmt.Println("There are no videos to download and upload! Please make sure your configuration is correct!")
		return
	} else if len(videos.Videos) >= 5 {
		fmt.Printf("Maximum videos downloaded and uploaded should be 5 at a time. Found %d videos!", len(videos.Videos))
		return
	}

	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	resources := make([]interface{}, len(videos.Videos))
	for i, v := range videos.Videos {
		resources[i] = v
	}

	pool := NewPool(len(videos.Videos))

	startTime := time.Now()
	fmt.Println("Start processing...")
	pool.start(resources, dlProcessor, ulProcessor)
	fmt.Println("\nTotal time taken ", time.Since(startTime))

	if pool.IsCompleted() {
		fmt.Println("\nAll done!")
		fmt.Println("Press ENTER to exit!")
		fmt.Scanln()
	}
}

func dlProcessor(resource interface{}) (interface{}, error) {
	video := reflect.ValueOf(resource).Interface().(Video)
	videoInfo, err := ytdl.GetVideoInfo(video.URL)
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %s", err.Error())
		return "", err
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

	fileName, err := createFileName(strconv.FormatInt(time.Now().UTC().UnixNano(), 10)+"."+formats[0].Extension, outputFileName{
		Title:         sanitizeFileNamePart(videoInfo.Title),
		Ext:           sanitizeFileNamePart(formats[0].Extension),
		DatePublished: sanitizeFileNamePart(videoInfo.DatePublished.Format("2006-01-02")),
		Resolution:    sanitizeFileNamePart(formats[0].Resolution),
		Author:        sanitizeFileNamePart(videoInfo.Author),
		Duration:      sanitizeFileNamePart(videoInfo.Duration.String()),
	})
	if err != nil {
		err = fmt.Errorf("Unable to parse output file file name: %s", err.Error())
		return "", err
	}

	downloadURL, err := videoInfo.GetDownloadURL(formats[0])
	if err != nil {
		err = fmt.Errorf("Unable to get download url: %s", err.Error())
		return "", err
	}

	var out io.Writer
	var logOut io.Writer = os.Stdout
	fmt.Println("\nDownload video [" + fileName + "] from [" + videoInfo.Title + "]")

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		err = fmt.Errorf("Unable to open output file: %s", err.Error())
		return "", err
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
		return "", err
	}
	contentSize := resp.ContentLength
	pb := pb.New64(contentSize).SetUnits(pb.U_BYTES)
	pb.ShowTimeLeft = true
	pb.ShowSpeed = true
	pb.Prefix("Downloading... [" + fileName + "]:")
	pb.Output = logOut
	pb.Start()
	defer pb.Finish()
	defer resp.Body.Close()

	out = io.MultiWriter(out, pb)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, err
}

func ulProcessor(result Result) error {
	result = reflect.ValueOf(result).Interface().(Result)
	video := result.Job.resource.(Video)
	fileName := result.Extra.(string)
	if fileName == "" {
		err := fmt.Errorf("You must provide a filename of a video file to upload")
		return err
	}

	// delete the file after everything to be done
	defer func(fileName string) {
		var err = os.Remove(fileName)
		if err != nil {
			fmt.Printf("Unable to delete video file: %s, %s", fileName, err.Error())
			return
		} else {
			fmt.Printf("Deleted video file [%s] successfully", fileName)
		}
	}(fileName)

	client := getClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		err = fmt.Errorf("Error creating YouTube client: %s", err.Error())
		return err
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       video.Title,
			Description: video.Desc,
			CategoryId:  video.Category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: video.Privacy},
	}

	if strings.Trim(video.Keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(video.Keywords, ",")
	} else {
		upload.Snippet.Tags = []string{video.Title}
	}

	call := service.Videos.Insert("snippet,status", upload)

	file, err := os.Open(fileName)
	if err != nil {
		err = fmt.Errorf("Error opening file: %s, %s", fileName, err.Error())
		return err
	}
	fileStat, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("Unable to ready file information %s, %s", fileName, err.Error())
		return err
	}
	fmt.Printf("\nUpload video [%s] as name: [%s]\n", fileName, video.Title)

	var logOut io.Writer = os.Stdout
	pb := pb.New64(fileStat.Size()).SetUnits(pb.U_BYTES)
	pb.ShowTimeLeft = true
	pb.ShowSpeed = true
	pb.Prefix("Uploading... [" + fileName + "]:")
	pb.Output = logOut
	pb.Start()
	pb.SetTotal64(fileStat.Size() * 10 / 100) // Init at 10%
	pb.Update()
	defer pb.Finish()
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		err = fmt.Errorf("Upload failure %s", err.Error())
		return err
	}

	pb.SetTotal64(fileStat.Size())
	pb.Update()
	fmt.Printf("\nUpload video [%s] as name: [%s] successfully! Video ID: %v\n", fileName, video.Title, response.Id)

	return err
}
