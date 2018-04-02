// Package slackutilsx is a utility package that doesn't promise API stability.
// its for experimental functionality and utilities.
package slackutilsx

import "unicode/utf8"

// ChannelType the type of channel based on the channelID
type ChannelType int

func (t ChannelType) String() string {
	switch t {
	case CTypeDM:
		return "Direct"
	case CTypeGroup:
		return "Group"
	case CTypeChannel:
		return "Channel"
	default:
		return "Unknown"
	}
}

const (
	// Unknown represents channels we cannot properly detect.
	CTypeUnknown ChannelType = iota
	// DM is a private channel between two slack users.
	CTypeDM
	// Group is a group channel.
	CTypeGroup
	// Channel is a public channel.
	CTypeChannel
)

// DetectChannelType converts a channelID to a ChannelType.
// channelID must not be empty. However, if it is not empty, the channel type will default to Unknown.
func DetectChannelType(channelID string) ChannelType {
	// intentionally ignore the error and just default to CTypeUnknown
	switch r, _ := utf8.DecodeRuneInString(channelID); r {
	case 'C':
		return CTypeChannel
	case 'G':
		return CTypeGroup
	case 'D':
		return CTypeDM
	default:
		return CTypeUnknown
	}
}
