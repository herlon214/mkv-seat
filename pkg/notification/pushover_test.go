package notification_test

import (
	"os"
	"testing"

	"github.com/herlon214/mkv-seat/pkg/notification"
	"github.com/stretchr/testify/assert"
)

func TestPushover(t *testing.T) {
	pushoverUserKey := os.Getenv("PUSHOVER_USER_KEY")
	pushoverApiKey := os.Getenv("PUSHOVER_API_KEY")

	pushover := notification.NewPushover(pushoverUserKey, pushoverApiKey)

	assert.Nil(t, pushover.Notify("New subtitle extract", "The subtitle for the file '[HorribleSubs] Kimetsu no Yaiba - 23 [1080p].mkv' was extracted successfully!"))
}
