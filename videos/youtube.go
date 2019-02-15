package videos

import (
	"fmt"
	"os"

	"github.com/rylio/ytdl"
)

// DlYt function
func DlYt() {
	vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=Zap_pTv2mgA")
	if err != nil {
		fmt.Println("error")
	}
	file, _ := os.Create(vid.Title + ".mp4")
	defer file.Close()
	err = vid.Download(ytdl.FORMATS[18], file)
	if err != nil {
		println(err)
	}
}
