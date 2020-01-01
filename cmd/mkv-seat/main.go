package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"

	"github.com/herlon214/mkv-seat/pkg/mkv"
	"github.com/herlon214/mkv-seat/pkg/translator"
)

var (
	apiKey          string
	languageFrom    string
	languageFromTag language.Tag
	languageTo      string
	languageToTag   language.Tag
	outputFormat    string
	skipExisting    bool
	logger          = logrus.New()
)

func main() {
	fmt.Println("          _                               _   ")
	fmt.Println("_ __ ___ | | ____   __     ___  ___  __ _| |_ ")
	fmt.Println("| '_ ` _ \\| |/ /\\ \\ / /    / __|/ _ \\/ _` | __|")
	fmt.Println("| | | | | |   <  \\ V /     \\__ \\  __/ (_| | |_ ")
	fmt.Println("|_| |_| |_|_|\\_\\  \\_/      |___/\\___|\\__,_|\\__|")
	fmt.Printf("\n\n")

	// Create the root command
	rootCmd := &cobra.Command{
		Use:  "mkv-seat file.mkv",
		Run:  Run,
		Args: cobra.MinimumNArgs(1),
	}

	// Add flags to the command
	rootCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "srt", "Output format, e.g: srt")
	rootCmd.Flags().StringVarP(&apiKey, "key", "k", "", "Google Cloud Translation Api Key, e.g: AIvaSyCiLjaWkykUoROHq2lqqbVoUA3ZTyv7xQI")
	rootCmd.Flags().StringVarP(&languageFrom, "lang-from", "f", "", "Original subtitle language (following the BCP 47), e.g: en")
	rootCmd.Flags().StringVarP(&languageTo, "lang-to", "t", "", "Output subtitle language (following the BCP 47), e.g: pt-BR")
	rootCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip files generation if the subtitle already exists")
	rootCmd.Execute()
}

// Run is the body of the root command
func Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrln("You must specify at least one mkv file")
	}

	// Check for valid language tags in order to fail fast
	if languageFrom != "" {
		tag, err := language.Parse(languageFrom)
		if err != nil {
			logger.Panicf("Invalid language from: %s", err.Error())
		}

		languageFromTag = tag
	}
	if languageTo != "" {
		tag, err := language.Parse(languageTo)
		if err != nil {
			logger.Panicf("Invalid language to: %s", err.Error())
		}

		languageToTag = tag
	}

	path := args[0]
	pathPieces := strings.Split(path, "/")
	file := pathPieces[len(pathPieces)-1]
	filePieces := strings.Split(file, ".")
	fileName := strings.Join(filePieces[0:len(filePieces)-1], ".")
	outputFolder := "."
	if len(pathPieces) > 1 {
		outputFolder = strings.Join(pathPieces[:len(pathPieces)-1], "/")
	}

	// Check if the skip is enabled
	if skipExisting && FileExists(fmt.Sprintf("%s/%s.str", outputFolder, fileName)) {
		logger.Infof("Skipping generation for %s", path)
		return
	}

	// Extract subtitle
	subtitle := mkv.ExtractSubtitle(path, logger)
	if subtitle == nil {
		logger.Panic("Failed to extract subtitle")
	}

	// Check if it needs to be translated
	if apiKey != "" && languageFrom != "" && languageTo != "" {
		logger.Infof("Translating subtitle from %s to %s", languageFrom, languageTo)
		subtitle = Translate(subtitle, logger)
	}

	// Save the file
	outputFile := fmt.Sprintf("%s/%s.%s", outputFolder, fileName, outputFormat)
	srt, err := os.Create(outputFile)
	if err != nil {
		logger.Panicf("Error to create output file: %s", err.Error())
	}

	switch outputFormat {
	case "srt":
		err = subtitle.WriteToSRT(srt)
	case "ass":
		err = subtitle.WriteToSSA(srt)
	case "ssa":
		err = subtitle.WriteToSSA(srt)
	}

	if err != nil {
		logger.Panicf("Error to save output %s.%s file: %s", fileName, outputFormat, err.Error())
	}

	logger.Infof("Output subtitle saved to %s", outputFile)

	return
}

// Translate an subtitle file and returns the subtitle
func Translate(subtitle *astisub.Subtitles, logger *logrus.Logger) *astisub.Subtitles {
	ctx := context.Background()

	// Put all texts into an array
	logger.Info("[Translation] Collecting texts...")
	texts := make([]string, 0)
	id := 0
	for _, item := range subtitle.Items {
		for _, line := range item.Lines {
			for _, lineItem := range line.Items {
				// Fix breaklines
				text := strings.Replace(lineItem.Text, "\\N", " \n ", -1)
				texts = append(texts, text)
				id++
			}
		}
	}

	logger.Infof("[Translation] Collected %d texts", len(texts))

	// Create the translator
	translator := translator.NewGoogleTranslator(ctx, apiKey, logger)

	// Translate the given texts
	translations, err := translator.Translate(texts, languageFromTag, languageToTag)
	if err != nil {
		logger.Fatalf("[Translation] Failed to translate texts: %s", err.Error())
	}

	// Update the texts
	id = 0
	for i, item := range subtitle.Items {
		for j, line := range item.Lines {
			for k := range line.Items {
				if id > len(translations)-1 {
					break
				}

				subtitle.Items[i].Lines[j].Items[k].Text = translations[id]
				id++
			}
		}
	}

	return subtitle
}

// FileExists verify if the srt already exists
func FileExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return true
	}

	return false
}
