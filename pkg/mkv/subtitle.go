package mkv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
)

// ExtractSubtitle extracts a subtitle from a given mkv
func ExtractSubtitle(filePath string, logger *logrus.Logger) *astisub.Subtitles {
	filePieces := strings.Split(filePath, "/")
	fileNameWithExtension := filePieces[len(filePieces)-1]
	fileNamePieces := strings.Split(fileNameWithExtension, ".")
	fileName := strings.Join(fileNamePieces[0:len(fileNamePieces)-1], ".")

	logger.Infof("[MKV] Extracting subtitle for %s...", fileName)

	outputSub := fmt.Sprintf("/tmp/mkv-seat/%s_temp_sub", fileName)

	input := []string{filePath, "tracks", fmt.Sprintf("2:%s", outputSub)}

	cmd := exec.Command("mkvextract", input...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Println(string(out))
		logger.Panicf("[MKV] Failed to extract subtitle for %s: %s", filePath, err.Error())
	}
	strOut := string(out)

	format, err := detectFromat(strOut)
	if err != nil {
		logger.Panic(err.Error())
	}

	// Rename the file with proper format
	newName := fmt.Sprintf("%s.%s", fileName, format)
	newFilePath := fmt.Sprintf("/tmp/mkv-seat/%s", newName)
	err = os.Rename(outputSub, newFilePath)
	if err != nil {
		logger.Panic(err)
	}

	// Open the file using the astibsub lib
	subtitle, err := astisub.OpenFile(newFilePath)
	if err != nil {
		logger.Panic(err)
	}

	// Remove the file from tmp
	err = os.Remove(newFilePath)
	if err != nil {
		logger.Errorf("[MKV] Failed to remove tmp file: %s", err.Error())
	}

	// Parse the subtitle
	for i, item := range subtitle.Items {
		for j, line := range item.Lines {
			for k := range line.Items {
				// Fix breaklines
				text := subtitle.Items[i].Lines[j].Items[k].Text
				text = strings.Replace(text, "\\N", " \n ", -1)

				subtitle.Items[i].Lines[j].Items[k].Text = text
			}
		}
	}

	logger.Info("[MKV] Subtitle extracted successfully")

	return subtitle
}

func detectFromat(output string) (string, error) {
	re := regexp.MustCompile("with the CodecID \\'(.*)\\' to the file")
	match := re.FindStringSubmatch(output)
	if len(match) == 0 {
		return "", errors.New("[MKV] Format not detected from the output")
	}

	format := match[1]

	switch format {
	case "S_TEXT/SSA":
		return "ass", nil
	case "S_TEXT/ASS":
		return "ass", nil
	case "S_SSA":
		return "ass", nil
	case "S_ASS":
		return "ass", nil
	case "S_TEXT/UTF8":
	case "S_TEXT/ASCII":
		return "srt", nil
	}

	return "", errors.New("Format not detected")
}
