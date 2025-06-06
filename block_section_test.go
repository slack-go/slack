package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSectionBlock(t *testing.T) {
	textInfo := NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|The Ritz-Carlton New Orleans>*\n★★★★★\n$340 per night\nRated: 9.1 - Excellent", false, false)
	sectionBlock := NewSectionBlock(textInfo, nil, nil, SectionBlockOptionBlockID("test_block"))

	assert.Equal(t, sectionBlock.BlockType(), MBTSection)
	assert.Equal(t, string(sectionBlock.Type), "section")
	assert.Equal(t, sectionBlock.BlockID, "test_block")
	assert.Equal(t, sectionBlock.ID(), "test_block")
	assert.Equal(t, len(sectionBlock.Fields), 0)
	assert.Nil(t, sectionBlock.Accessory)
	assert.Equal(t, sectionBlock.Text.Type, "mrkdwn")
	assert.Contains(t, sectionBlock.Text.Text, "New Orleans")
}

func TestNewBlockSectionContainsAddedTextBlockAndAccessory(t *testing.T) {
	textBlockObject := NewTextBlockObject("mrkdwn", "You have a new test: *Hi there* :wave:", true, false)
	conflictImage := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/notificationsWarningIcon.png", "notifications warning icon")
	sectionBlock := NewSectionBlock(textBlockObject, nil, NewAccessory(conflictImage))

	assert.Equal(t, sectionBlock.BlockType(), MBTSection)
	assert.Equal(t, len(sectionBlock.BlockID), 0)
	textBlockInSection := sectionBlock.Text
	assert.Equal(t, textBlockInSection.Text, textBlockObject.Text)
	assert.Equal(t, textBlockInSection.Type, textBlockObject.Type)
	assert.Nil(t, textBlockInSection.Emoji)
	assert.False(t, textBlockInSection.Verbatim)
	assert.Equal(t, sectionBlock.Accessory.ImageElement, conflictImage)
}

func TestSectionBlockOptionExpand(t *testing.T) {
	textInfo := NewTextBlockObject("mrkdwn", "This is a long text that should be auto-expanded", false, false)

	// Create a section block with expand option set to true
	sectionBlock := NewSectionBlock(textInfo, nil, nil, SectionBlockOptionExpand(true))

	// Verify that the expand field is set correctly
	assert.True(t, sectionBlock.Expand)

	// Create another section block with expand option set to false
	sectionBlock = NewSectionBlock(textInfo, nil, nil, SectionBlockOptionExpand(false))

	// Verify that the expand field is set correctly
	assert.False(t, sectionBlock.Expand)
}
