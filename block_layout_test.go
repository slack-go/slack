package slack

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLayoutBlock(t *testing.T) {
	section := NewSectionBlock(NewTextBlockObject(MarkdownType, "*Status*\nReady", false, false), nil, nil)
	item := NewLayoutItem(section).WithWidth(LayoutWidthWide).WithHorizontalAlignment(LayoutHorizontalAlignmentEnd)
	row := NewLayoutRow(item).AddItem(NewLayoutItem(NewDividerBlock()))
	block := NewLayoutBlock("Release", row).
		WithBlockID("release-layout").
		WithLabelVisible(true).
		WithGridLines(LayoutGridLinesRows).
		WithFallback(section).
		AddRow(NewLayoutRow(NewLayoutItem(section)))

	assert.Equal(t, MBTLayout, block.BlockType())
	assert.Equal(t, "release-layout", block.ID())
	assert.Equal(t, "Release", block.Label)
	assert.True(t, block.IsLabelVisible)
	assert.Equal(t, LayoutGridLinesRows, block.GridLines)
	require.Len(t, block.Rows, 2)
	require.Len(t, block.Rows[0].Items, 2)
	assert.Same(t, section, item.Block)
	assert.Equal(t, LayoutWidthWide, item.Width)
	assert.Equal(t, LayoutHorizontalAlignmentEnd, item.HorizontalAlignment)
	require.NotNil(t, block.Fallback)
	assert.Equal(t, []Block{section}, block.Fallback.BlockSet)
	require.NoError(t, block.Validate())
}

func TestLayoutBlockRoundTrip(t *testing.T) {
	for _, payload := range []string{
		`{"type":"layout","label":"Minimal","rows":[{"items":[{"block":{"type":"divider"}}]}]}`,
		`{"type":"layout","label":"Empty fallback","rows":[{"items":[{"block":{"type":"divider"}}]}],"fallback":[]}`,
		`{"type":"layout","block_id":"summary","label":"Summary","is_label_visible":true,"grid_lines":"outer","rows":[{"items":[{"width":"wide","horizontal_alignment":"start","block":{"type":"section","text":{"type":"mrkdwn","text":"*Ready*"}}},{"width":"narrow","horizontal_alignment":"end","block":{"type":"actions","elements":[{"type":"button","action_id":"view","text":{"type":"plain_text","text":"View"}}]}}]}],"fallback":[{"type":"section","text":{"type":"mrkdwn","text":"Ready — View"}}]}`,
		`{"type":"layout","label":"Container child","rows":[{"items":[{"block":{"type":"container","title":{"type":"plain_text","text":"Details"},"child_blocks":[{"type":"section","text":{"type":"mrkdwn","text":"Nested"}}]}}]}],"fallback":[]}`,
		`{"type":"container","title":{"type":"plain_text","text":"Summary"},"child_blocks":[{"type":"layout","label":"Nested layout","rows":[{"items":[{"block":{"type":"divider"}}]}],"fallback":[]}]}`,
		`{"type":"layout","label":"Future child","rows":[{"items":[{"block":{"type":"future_block","new_field":{"enabled":true}}}]}],"fallback":[]}`,
	} {
		t.Run(payload, func(t *testing.T) {
			var blocks Blocks
			require.NoError(t, json.Unmarshal([]byte("["+payload+"]"), &blocks))
			require.Len(t, blocks.BlockSet, 1)
			var layout *LayoutBlock
			if container, ok := blocks.BlockSet[0].(*ContainerBlock); ok {
				require.IsType(t, &LayoutBlock{}, container.ChildBlocks.BlockSet[0])
				layout = container.ChildBlocks.BlockSet[0].(*LayoutBlock)
			} else {
				require.IsType(t, &LayoutBlock{}, blocks.BlockSet[0])
				layout = blocks.BlockSet[0].(*LayoutBlock)
			}
			require.Len(t, layout.Rows, 1)
			switch layout.Label {
			case "Minimal", "Empty fallback", "Nested layout":
				assert.IsType(t, &DividerBlock{}, layout.Rows[0].Items[0].Block)
			case "Summary":
				assert.IsType(t, &SectionBlock{}, layout.Rows[0].Items[0].Block)
				actions, ok := layout.Rows[0].Items[1].Block.(*ActionBlock)
				require.True(t, ok)
				assert.IsType(t, &ButtonBlockElement{}, actions.Elements.ElementSet[0])
				assert.IsType(t, &SectionBlock{}, layout.Fallback.BlockSet[0])
			case "Container child":
				child, ok := layout.Rows[0].Items[0].Block.(*ContainerBlock)
				require.True(t, ok)
				assert.IsType(t, &SectionBlock{}, child.ChildBlocks.BlockSet[0])
			case "Future child":
				assert.IsType(t, &UnknownBlock{}, layout.Rows[0].Items[0].Block)
			}
			encoded, err := json.Marshal(blocks)
			require.NoError(t, err)
			assert.JSONEq(t, "["+payload+"]", string(encoded))
		})
	}
}

func TestLayoutBlockFallback(t *testing.T) {
	block := NewLayoutBlock("Summary", NewLayoutRow(NewLayoutItem(NewDividerBlock())))
	encoded, err := json.Marshal(block)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), `"fallback"`)
	block.WithFallback()
	encoded, err = json.Marshal(block)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"fallback":[]`)
}

func TestLayoutItemUnmarshalErrors(t *testing.T) {
	for _, payload := range []string{`{}`, `{"block":null}`, `{"block":[]}`, `{"block":42}`, `{"block":{"type":"section","text":42}}`} {
		t.Run(payload, func(t *testing.T) {
			var blocks Blocks
			err := json.Unmarshal([]byte(`[{"type":"layout","label":"Invalid","rows":[{"items":[`+payload+`]}]}]`), &blocks)
			require.Error(t, err)
		})
	}
}

func TestLayoutBlockValidate(t *testing.T) {
	valid := func() *LayoutBlock { return NewLayoutBlock("Summary", NewLayoutRow(NewLayoutItem(NewDividerBlock()))) }
	for _, tt := range []struct {
		name   string
		change func(*LayoutBlock)
	}{
		{"wrong type", func(b *LayoutBlock) { b.Type = MBTContainer }},
		{"missing label", func(b *LayoutBlock) { b.Label = "" }},
		{"long label", func(b *LayoutBlock) { b.Label = strings.Repeat("界", 151) }},
		{"missing rows", func(b *LayoutBlock) { b.Rows = nil }},
		{"too many rows", func(b *LayoutBlock) { b.Rows = make([]*LayoutRow, 6) }},
		{"nil row", func(b *LayoutBlock) { b.Rows[0] = nil }},
		{"missing items", func(b *LayoutBlock) { b.Rows[0].Items = nil }},
		{"too many items", func(b *LayoutBlock) { b.Rows[0].Items = make([]*LayoutItem, 4) }},
		{"nil item", func(b *LayoutBlock) { b.Rows[0].Items[0] = nil }},
		{"missing block", func(b *LayoutBlock) { b.Rows[0].Items[0].Block = nil }},
		{"typed nil block", func(b *LayoutBlock) { b.Rows[0].Items[0].Block = (*SectionBlock)(nil) }},
		{"markdown child", func(b *LayoutBlock) { b.Rows[0].Items[0].Block = NewMarkdownBlock("", "Invalid") }},
		{"layout child", func(b *LayoutBlock) { b.Rows[0].Items[0].Block = valid() }},
		{"invalid width", func(b *LayoutBlock) { b.Rows[0].Items[0].Width = "full" }},
		{"invalid alignment", func(b *LayoutBlock) { b.Rows[0].Items[0].HorizontalAlignment = "left" }},
		{"invalid grid lines", func(b *LayoutBlock) { b.GridLines = "columns" }},
		{"interactive fallback", func(b *LayoutBlock) { b.WithFallback(NewActionBlock("actions")) }},
		{"nil fallback block", func(b *LayoutBlock) { b.WithFallback(nil) }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			block := valid()
			tt.change(block)
			require.Error(t, block.Validate())
		})
	}
	block := valid()
	block.Label = strings.Repeat("界", 150)
	block.WithFallback()
	for range 4 {
		block.AddRow(NewLayoutRow(NewLayoutItem(NewDividerBlock()), NewLayoutItem(NewDividerBlock()), NewLayoutItem(NewDividerBlock())))
	}
	require.NoError(t, block.Validate(), "Unicode label and row/item limits are inclusive")
	for _, width := range []LayoutWidth{"", LayoutWidthNarrow, LayoutWidthStandard, LayoutWidthWide} {
		for _, alignment := range []LayoutHorizontalAlignment{"", LayoutHorizontalAlignmentStart, LayoutHorizontalAlignmentCenter, LayoutHorizontalAlignmentEnd} {
			for _, lines := range []LayoutGridLines{"", LayoutGridLinesNone, LayoutGridLinesRows, LayoutGridLinesOuter} {
				block := valid().WithGridLines(lines)
				block.Rows[0].Items[0].WithWidth(width).WithHorizontalAlignment(alignment)
				require.NoError(t, block.Validate())
			}
		}
	}
}
