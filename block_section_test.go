package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSectionBlock(t *testing.T) {

	textInfo := NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|The Ritz-Carlton New Orleans>*\n★★★★★\n$340 per night\nRated: 9.1 - Excellent", false, false)

	sectionBlock := NewSectionBlock(textInfo, nil, nil, SectionBlockOptionBlockID("test_block"))
	assert.Equal(t, string(sectionBlock.Type), "section")
	assert.Equal(t, string(sectionBlock.BlockID), "test_block")
	assert.Equal(t, len(sectionBlock.Fields), 0)
	assert.Nil(t, sectionBlock.Accessory)
	assert.Equal(t, sectionBlock.Text.Type, "mrkdwn")
	assert.Contains(t, sectionBlock.Text.Text, "New Orleans")

}
