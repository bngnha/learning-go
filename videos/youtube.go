package videos

import (
	"fmt"
	"os"

	"github.com/rylio/ytdl"
)

// DlYt function
func DlYt() {
	vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=BzYpUt1PFFI")
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %s", err.Error())
		return
	}

	file, _ := os.Create(vid.Title + ".mp4")
	defer file.Close()
	err = vid.Download(newFormat(vid.Formats[0]), file)
	if err != nil {
		println(err)
	}
}

func newFormat(itag int) (Format, bool) {
	if f, ok := FORMATS[itag]; ok {
		f.meta = make(map[string]interface{})
		return f, true
	}
	return Format{}, false
}
