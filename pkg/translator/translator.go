package translator

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// Translator is the interface for all the translators
type Translator interface {
	Translate(texts []string, from, to language.Tag) ([]string, error)
}

// NewGoogleTranslator returns a new Google Translator instance
func NewGoogleTranslator(ctx context.Context, apiKey string, logger *logrus.Logger) Translator {
	return &GoogleTranslator{
		logger:          logger,
		limitPerRequest: 100,
		apiKey:          apiKey,
		ctx:             ctx,
	}
}
