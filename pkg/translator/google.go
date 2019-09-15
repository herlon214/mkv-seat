package translator

import (
	"context"
	"log"
	"math"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// GoogleTranslator implementation
type GoogleTranslator struct {
	apiKey          string
	limitPerRequest int

	logger *logrus.Logger
	ctx    context.Context
}

// Translate the given texts
func (gt *GoogleTranslator) Translate(texts []string, from, to language.Tag) ([]string, error) {
	translatedTexts := make([]string, 0)
	client, err := translate.NewClient(gt.ctx, option.WithAPIKey(gt.apiKey))
	if err != nil {
		log.Fatal(err)
	}

	// Divide the texts by chunks and round it up
	chunks := int(math.Ceil(float64(len(texts)) / float64(gt.limitPerRequest)))
	for i := 0; i < chunks; i++ {
		gt.logger.Infof("[Translation] Requesting translations %d/%d...", i, chunks)
		start := i * gt.limitPerRequest
		end := math.Min(float64(i*gt.limitPerRequest+gt.limitPerRequest), float64(len(texts)))
		actualChunk := texts[start:int(end)]

		// Request translations
		opts := translate.Options{
			Source: from,
		}
		resp, err := client.Translate(gt.ctx, actualChunk, to, &opts)
		if err != nil {
			gt.logger.Errorf("[Translation] Failed to request translations for chunk %d: %s", i, err.Error())
			return nil, err
		}

		// Append results
		for _, translation := range resp {
			translatedTexts = append(translatedTexts, filterText(translation.Text))
		}
	}

	return translatedTexts, nil
}

func filterText(text string) string {
	text = strings.Replace(text, "&quot;", `"`, -1)

	return text
}
