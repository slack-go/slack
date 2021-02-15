package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCallBlock(t *testing.T) {

	call := &ZoomCall{
		MediaBackendType: MBETPlatformCall,
		Info: ZoomCallInfo{
			ActiveParticipants: []*CallParticipant{},
			AllParticipants: []*CallParticipant{
				&CallParticipant{
					ID:          "KArkyRIMNvAuHA-hxSsoPg",
					AvatarUrl:   "",
					DisplayName: "Adam Savage",
				},
			},
			AppId:             "A5GE9BMQC",
			Channels:          []string{"FU1NC8HAN"},
			CreatedBy:         "U0106V8C3DM",
			DateEnd:           0,
			DateStart:         1613383918,
			DesktopAppJoinUrl: "zoommtg://zoom.us/join?action=join\u0026confno=9097887221\u0026pwd=Y1MzaVl3QlArYnhvWnlqdEtqUlE1UT09\u0026confid=dXNzPWdfTm1CNkhWd1BDaXVYX1lmUG51RTk0eVFBd3NMNFVSUDBmMjBBMlJyLWdrVGRJNXlMQnFCeGJCV0V1eWlIRHZYY1FUbEF4N2p4NzM0YldsV2VuWFpNNVp3S1UueWtsQWRCMmdfSWtpQzhUcA%3D%3D\u0026t=1613383917995",
			DisplayId:         "909-7887-221",
			HasEnded:          false,
			Id:                "M01MA5CBOOK",
			IsDmCall:          false,
			JoinUrl:           "https://zoom.us/j/9097887221?pwd=Y1MzaVl3QlArYnhvWnlqdEtqUlE1UT09",
			Name:              "Zoom meeting started by Adam Savage",
			WasAccepted:       false,
			WasMissed:         false,
			WasRejected:       false,
		},
	}
	callBlock := NewCallBlock("8JP", "R01MX5MCBQE", call)
	assert.Equal(t, string(callBlock.Type), "call")

}
