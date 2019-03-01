package videos

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	config "github.com/micro/go-config"
	"github.com/rylio/ytdl"
)

// WorkPool function
func WorkPool() {
	config.LoadFile("config.json")
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
	pool.start(resources, resourceProcessor, resultProcessor)

	if pool.IsCompleted() {
		fmt.Println("All done!")
	}
}

// Video struct
type Video struct {
	URL      string
	Title    string
	Desc     string
	Category string
	Keywords string
	Privacy  string
	FileName string
}

// Videos struct
type Videos struct {
	Videos []Video
}

// ProcessorFunc callback function
type ProcessorFunc func(resource interface{}) error

// ResultProcessorFunc callback function
type ResultProcessorFunc func(result Result) error

// Job struct
type Job struct {
	id       int
	resource interface{}
}

// Result struct
type Result struct {
	Job Job
	Err error
}

// Pool struct
type Pool struct {
	numRoutines int
	jobs        chan Job
	results     chan Result
	done        chan bool
	completed   bool
}

// NewPool is a Pool constructor
func NewPool(numRoutines int) *Pool {
	r := &Pool{numRoutines: numRoutines}
	r.jobs = make(chan Job, numRoutines)
	r.results = make(chan Result, numRoutines)
	r.done = make(chan bool)

	return r
}

// Start function
func (p *Pool) start(resource []interface{}, proFunc ProcessorFunc, resFunc ResultProcessorFunc) {
	startTime := time.Now()
	fmt.Println("worker pool starting...")

	go p.allocate(resource)
	go p.collect(resFunc)
	go p.workerPool(proFunc)

	<-p.done

	fmt.Println("total time taken ", time.Since(startTime))
}

func (p *Pool) allocate(jobs []interface{}) {
	defer close(p.jobs)
	for i, v := range jobs {
		p.jobs <- Job{id: i, resource: v}
	}
}

func (p *Pool) work(wg *sync.WaitGroup, processor ProcessorFunc) {
	defer wg.Done()
	fmt.Println("gRountine work starting")
	for job := range p.jobs {
		fmt.Printf("working on job ID %d\n", job.id)
		p.results <- Result{job, processor(job.resource)}
		fmt.Printf("done with job ID %d\n", job.id)
	}
	fmt.Println("gRountine work done")
}

func (p *Pool) workerPool(processor ProcessorFunc) {
	defer close(p.results)
	fmt.Printf("Worker Pool spawning new goRoutines, total: %d\n", p.numRoutines)
	var wg sync.WaitGroup
	for i := 0; i < p.numRoutines; i++ {
		wg.Add(1)
		go p.work(&wg, processor)
		fmt.Printf("Spawned work goRoutine %d\n", i)
	}
	fmt.Println("Worker Pool done spawning work gRoutines")
	wg.Wait()

	fmt.Println("All work goRoutine done processing")
}

func (p *Pool) collect(proc ResultProcessorFunc) {
	fmt.Println("goRoutine collect starting")
	for result := range p.results {
		outcome := proc(result)
		fmt.Printf("Job with id: %d completed, outcome: %s", result.Job.id, outcome)
	}
	fmt.Println("goRoutine collect done, setting channel done as completed")
	p.done <- true
	p.completed = true
}

// IsCompleted function
func (p *Pool) IsCompleted() bool {
	return p.completed
}

func resourceProcessor(resource interface{}) (string, error) {
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
	fmt.Println("\nDownloading video [" + videoInfo.Title + "]...")

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
	defer resp.Body.Close()
	contentSize := resp.ContentLength
	progressBar := pb.New64(contentSize).SetUnits(pb.U_BYTES)
	progressBar.ShowTimeLeft = true
	progressBar.ShowSpeed = true
	progressBar.SetWidth(80)
	progressBar.Output = logOut
	progressBar.Start()
	defer progressBar.Finish()

	out = io.MultiWriter(out, progressBar)
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, err
}

func resultProcessor(result Result) error {
	//fmt.Println("Result processor", result)
	return nil
}
