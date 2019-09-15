package mkv

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestExtractMKVSubtitle(t *testing.T) {
	filePath := "../../testdata/SampleVideo.mkv"
	logger := logrus.New()

	expectedFile, err := ioutil.ReadFile("../../testdata/example_subtitle.ass")
	if err != nil {
		t.Errorf("Failed to open example subtitle: %s", err.Error())
	}

	finalFile, err := os.Create("/tmp/mkv-seat/example-subtitle.ass")
	if err != nil {
		t.Errorf("Error to create output file: %s", err.Error())
	}

	subtitle := ExtractSubtitle(filePath, logger)
	subtitle.WriteToSSA(finalFile)

	extractedSub, err := ioutil.ReadFile("/tmp/mkv-seat/example-subtitle.ass")
	if err != nil {
		t.Errorf("Error to read the extracted file: %s", err.Error())
	}

	expected := strings.TrimSpace(string(expectedFile))
	actual := strings.TrimSpace(string(extractedSub))

	assert.Equal(t, expected, actual, "Extracted subtitle not equal example subtitle")
}
