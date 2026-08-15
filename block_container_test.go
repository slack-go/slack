package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainerBlock(t *testing.T) {
	section := NewSectionBlock(NewTextBlockObject("mrkdwn", "Content", false, false), nil, nil)
	divider := NewDividerBlock()

	block := NewContainerBlock(section).
		WithBlockID("container-1").
		WithTitle(NewTextBlockObject(PlainTextType, "Title", false, false)).
		WithSubtitle(NewTextBlockObject(MarkdownType, "Subtitle", false, false)).
		WithIcon(NewImageBlockElement("https://example.com/icon.png", "icon")).
		WithWidth(ContainerWidthWide).
		WithCollapsible(true, true)

	assert.Equal(t, MBTContainer, block.BlockType())
	assert.Equal(t, "container", string(block.Type))
	assert.Equal(t, "container-1", block.ID())
	assert.Equal(t, ContainerWidthWide, block.Width)
	assert.True(t, block.IsCollapsible)
	assert.True(t, block.DefaultCollapsed)
	require.Len(t, block.ChildBlocks.BlockSet, 1)
	assert.Equal(t, section, block.ChildBlocks.BlockSet[0])

	block.AddChildBlock(divider)
	require.Len(t, block.ChildBlocks.BlockSet, 2)
	assert.Equal(t, divider, block.ChildBlocks.BlockSet[1])

	assert.NoError(t, block.Validate())
}

func TestContainerBlockValidate(t *testing.T) {
	child := NewSectionBlock(NewTextBlockObject("mrkdwn", "Content", false, false), nil, nil)
	plainTitle := NewTextBlockObject(PlainTextType, "Title", false, false)

	tests := []struct {
		name    string
		block   *ContainerBlock
		wantErr bool
	}{
		{
			name:  "valid",
			block: NewContainerBlock(child).WithTitle(plainTitle),
		},
		{
			name:  "valid with rich_text_title",
			block: NewContainerBlock(child).WithRichTextTitle(NewRichTextBlock("rtt")),
		},
		{
			name:    "missing title",
			block:   NewContainerBlock(child),
			wantErr: true,
		},
		{
			name:    "non plain_text title",
			block:   NewContainerBlock(child).WithTitle(NewTextBlockObject(MarkdownType, "Title", false, false)),
			wantErr: true,
		},
		{
			name:    "no child blocks",
			block:   NewContainerBlock().WithTitle(plainTitle),
			wantErr: true,
		},
		{
			name: "too many child blocks",
			block: NewContainerBlock(
				child, child, child, child, child, child, child, child, child, child, child,
			).WithTitle(plainTitle),
			wantErr: true,
		},
		{
			name:    "invalid width",
			block:   NewContainerBlock(child).WithTitle(plainTitle).WithWidth("gigantic"),
			wantErr: true,
		},
		{
			name:    "header divider on collapsible",
			block:   NewContainerBlock(child).WithTitle(plainTitle).WithCollapsible(true, false).WithHeaderDivider(true),
			wantErr: true,
		},
		{
			name:    "default collapsed without collapsible",
			block:   &ContainerBlock{Type: MBTContainer, Title: plainTitle, DefaultCollapsed: true, ChildBlocks: Blocks{BlockSet: []Block{child}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContainerBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "container",
		"block_id": "container-1",
		"title": {"type": "plain_text", "text": "Deploy status"},
		"subtitle": {"type": "mrkdwn", "text": "*production*"},
		"icon": {"type": "image", "image_url": "https://example.com/icon.png", "alt_text": "icon"},
		"width": "wide",
		"is_collapsible": true,
		"default_collapsed": true,
		"child_blocks": [
			{"type": "section", "text": {"type": "mrkdwn", "text": "All systems go"}},
			{"type": "divider"}
		]
	}`

	var block ContainerBlock
	require.NoError(t, json.Unmarshal([]byte(payload), &block))

	assert.Equal(t, MBTContainer, block.BlockType())
	assert.Equal(t, "container-1", block.ID())
	assert.Equal(t, ContainerWidthWide, block.Width)
	assert.True(t, block.IsCollapsible)
	require.NotNil(t, block.Title)
	assert.Equal(t, "Deploy status", block.Title.Text)
	require.Len(t, block.ChildBlocks.BlockSet, 2)
	assert.Equal(t, MBTSection, block.ChildBlocks.BlockSet[0].BlockType())
	assert.Equal(t, MBTDivider, block.ChildBlocks.BlockSet[1].BlockType())

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var want, got any
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(marshalled, &got))
	assert.Equal(t, want, got)
}
