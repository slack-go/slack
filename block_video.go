package slack

// VideoBlock defines data required to display a video as a block element
//
// More Information: https://api.slack.com/reference/block-kit/blocks#video
type VideoBlock struct {
	Type            MessageBlockType `json:"type"`
	VideoURL        string           `json:"video_url"`
	ThumbnailURL    string           `json:"thumbnail_url"`
	AltText         string           `json:"alt_text"`
	Title           *TextBlockObject `json:"title"`
	BlockID         string           `json:"block_id,omitempty"`
	TitleURL        string           `json:"title_url,omitempty"`
	AuthorName      string           `json:"author_name,omitempty"`
	ProviderName    string           `json:"provider_name,omitempty"`
	ProviderIconURL string           `json:"provider_icon_url,omitempty"`
	Description     *TextBlockObject `json:"description,omitempty"`
}

// BlockType returns the type of the block
func (s VideoBlock) BlockType() MessageBlockType {
	return s.Type
}

// NewVideoBlock returns an instance of a new Video Block type
func NewVideoBlock(videoURL, thumbnailURL, altText, blockID string, title *TextBlockObject) *VideoBlock {
	return &VideoBlock{
		Type:         MBTVideo,
		VideoURL:     videoURL,
		ThumbnailURL: thumbnailURL,
		AltText:      altText,
		BlockID:      blockID,
		Title:        title,
	}
}
