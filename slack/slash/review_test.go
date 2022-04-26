package slash

import (
	"testing"

	"github.com/AngelVI13/slack-assistant/device"
)

func TestChooseReviewer(t *testing.T) {
	users := map[string]device.AccessRight{
		"angel.iliev":    device.ADMIN,
		"laima.strigo":   device.STANDARD,
		"aurimas.razmis": device.ADMIN,
	}
	sender := "angel.iliev"

	// Seems a bit pointless but its good to check that sender is never the reviewer
	// Since chooseReviewer is done randomly - run it a few times to make sure this doesn't happen
	for i := 0; i < 100; i += 1 {
		reviewer := chooseReviewer(sender, users)
		if reviewer == sender {
			t.Errorf("Chosen reviewer is the same as sender! reviewer == sender == %s", reviewer)
			break
		}

		if _, ok := users[reviewer]; !ok {
			t.Errorf("Reviewer is not part of users list! reviewer=%s, users=%+v", reviewer, users)
			break
		}
	}
}
