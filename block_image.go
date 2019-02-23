package slack

import "strings"

// ImageBlock defines data required to display an image as a block element
//
// More Information: https://api.slack.com/reference/messaging/blocks#image
type ImageBlock struct {
	Type     string           `json:"type"`
	ImageURL string           `json:"image_url"`
	AltText  string           `json:"alt_text"`
	BlockID  string           `json:"block_id,omitempty"`
	Title    *TextBlockObject `json:"title"`
}

// ValidateBlock ensures that the type set to the block is found in the list of
// valid slack block.
func (s *ImageBlock) ValidateBlock() bool {
	return isStringInSlice(validBlockList, strings.ToLower(s.Type))

}

// NewImageBlock returns an instance of a new Image Block type
func NewImageBlock(imageURL, altText, blockID string, title *TextBlockObject) *ImageBlock {
	return &ImageBlock{
		Type:     "image",
		ImageURL: imageURL,
		AltText:  altText,
		BlockID:  blockID,
		Title:    title,
	}
}
