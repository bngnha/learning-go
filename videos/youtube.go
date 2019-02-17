package videos

import (
	"fmt"
	"os"

	"github.com/rylio/ytdl"
)

// DlYt function
func DlYt() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

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
