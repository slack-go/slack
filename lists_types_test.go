package slack

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListItemFieldEncoding(t *testing.T) {
	for _, tt := range []struct {
		name  string
		field ListItemField
		empty ListItemField
		value string
	}{
		{"rich_text", ListItemField{RichText: []RichTextBlock{*NewRichTextBlock("", NewRichTextSection(NewRichTextSectionTextElement("Tasks", nil)))}}, ListItemField{RichText: []RichTextBlock{}}, `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Tasks"}]}]}]`},
		{"user", ListItemField{User: []string{"U123"}}, ListItemField{User: []string{}}, `["U123"]`},
		{"date", ListItemField{Date: []string{"2025-10-10"}}, ListItemField{Date: []string{}}, `["2025-10-10"]`},
		{"select", ListItemField{Select: []string{"done"}}, ListItemField{Select: []string{}}, `["done"]`},
		{"number", ListItemField{Number: []float64{0, 1.5}}, ListItemField{Number: []float64{}}, `[0,1.5]`},
		{"email", ListItemField{Email: []string{"test@example.com"}}, ListItemField{Email: []string{}}, `["test@example.com"]`},
		{"phone", ListItemField{Phone: []string{"+1-555-123-4567"}}, ListItemField{Phone: []string{}}, `["+1-555-123-4567"]`},
		{"attachment", ListItemField{Attachment: []string{"F123"}}, ListItemField{Attachment: []string{}}, `["F123"]`},
		{"link", ListItemField{Link: []ListItemLink{{OriginalURL: "https://docs.slack.dev/", DisplayName: "Slack Developer Docs"}}}, ListItemField{Link: []ListItemLink{}}, `[{"original_url":"https://docs.slack.dev/","display_as_url":false,"display_name":"Slack Developer Docs"}]`},
		{"message", ListItemField{Message: []string{"https://example.slack.com/archives/C123/p1234567890123456"}}, ListItemField{Message: []string{}}, `["https://example.slack.com/archives/C123/p1234567890123456"]`},
		{"rating", ListItemField{Rating: []int{0}}, ListItemField{Rating: []int{}}, `[0]`},
		{"timestamp", ListItemField{Timestamp: []int64{1704067200}}, ListItemField{Timestamp: []int64{}}, `[1704067200]`},
		{"channel", ListItemField{Channel: []string{"C123"}}, ListItemField{Channel: []string{}}, `["C123"]`},
		{"reference", ListItemField{Reference: []ListItemReference{{File: &ListItemFileReference{FileID: "F123"}}}}, ListItemField{Reference: []ListItemReference{}}, `[{"file":{"file_id":"F123"}}]`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.field.ColumnID = "Col123"
			data, err := json.Marshal(tt.field)
			require.NoError(t, err)
			assert.JSONEq(t, `{"column_id":"Col123","`+tt.name+`":`+tt.value+`}`, string(data))
			tt.empty.ColumnID = "Col123"
			data, err = json.Marshal(tt.empty)
			require.NoError(t, err)
			assert.JSONEq(t, `{"column_id":"Col123","`+tt.name+`":[]}`, string(data))
		})
	}
	for _, value := range []bool{true, false} {
		data, err := json.Marshal(ListItemField{ColumnID: "Col123", Checkbox: &value})
		require.NoError(t, err)
		var fields map[string]any
		require.NoError(t, json.Unmarshal(data, &fields))
		assert.Equal(t, map[string]any{"column_id": "Col123", "checkbox": value}, fields)
	}
	data, err := json.Marshal(ListItemField{ColumnID: "Col123"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"column_id":"Col123"}`, string(data))
}

func TestListItemResponseFields(t *testing.T) {
	// Field shapes follow slackLists.items.list's response template. Checkbox is
	// a scalar in the official Node SDK, despite brackets in that template.
	const payload = `{
  "ok":true,
  "items":[{
   "id":"Rec123", "list_id":"F123", "date_created":1758744345,
   "created_by":"U123", "updated_by":"U456", "updated_timestamp":"1758744346",
   "parent_record_id":"RecParent", "archived":true, "is_subscribed":false,
   "saved":{"is_archived":false,"date_due":1758744347,"state":"in_progress"},
   "saved_fields":{"ColDate":{"is_archived":true,"date_completed":1758744346,"state":"completed"}},
   "fields":[
    {"column_id":"ColText","key":"title","value":"Tasks","text":"Tasks","rich_text":[{"type":"rich_text","block_id":"block1","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Tasks","style":{"bold":true}}]}]}]},
    {"column_id":"ColNumber","value":1.5,"number":[0,1.5]},
    {"column_id":"ColCheckbox","value":false,"checkbox":false},
    {"column_id":"ColEmpty","value":null},
    {"column_id":"ColDate","date":["2025-10-10"],"timestamp":[-1]},
    {"column_id":"ColUser","user":["U123"]},
    {"column_id":"ColSelect","select":["done"]},
    {"column_id":"ColAttachment","attachment":["F456"]},
    {"column_id":"ColChannel","channel":["C123"]},
    {"column_id":"ColEmail","email":["test@example.com"]},
    {"column_id":"ColPhone","phone":["+1-555-123-4567"]},
    {"column_id":"ColRating","rating":[4]},
    {"column_id":"ColMessage","message":[{"value":"https://example.slack.com/archives/C123/p1234567890123456","channel_id":"C123","ts":"1234567890.123456","thread_ts":"1234567890.000001"}]},
    {"column_id":"ColLink","link":[{"originalUrl":"https://docs.slack.dev/","displayAsUrl":false,"displayName":"Docs","attachment":{"title":"Docs"}}]},
    {"column_id":"ColReference","reference":[
     {"file":{"file_id":"F456"}},
     {"message":{"channel_id":"C123","ts":"1234567890.123456","thread_ts":"1234567890.000001"}},
     {"list_record":{"list_id":"F123","row_id":"Rec456"}},
     {"canvas_section":{"file_id":"F789","section_id":"temp:C:ABC"}}
    ]}
   ]
  }],
  "list":{"id":"F123","title":"Sprint Board","list_limits":{"row_count":7}},
  "response_metadata":{"next_cursor":"next-page"}
 }`
	api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.FormValue("include_list"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	})
	response, err := api.GetListItems("F123", GetListItemsOptionIncludeList(true))
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	item := response.Items[0]
	assert.Equal(t, "U123", item.CreatedBy)
	assert.Equal(t, "U456", item.UpdatedBy)
	assert.Equal(t, "RecParent", item.ParentRecordID)
	assert.True(t, item.Archived)
	require.NotNil(t, item.IsSubscribed)
	assert.False(t, *item.IsSubscribed)
	require.NotNil(t, item.Saved)
	assert.Equal(t, ListItemReminderInProgress, item.Saved.State)
	assert.Equal(t, int64(1758744347), item.Saved.DateDue)
	assert.Equal(t, ListItemReminderCompleted, item.SavedFields["ColDate"].State)
	assert.Equal(t, int64(1758744346), item.SavedFields["ColDate"].DateCompleted)
	assert.True(t, item.SavedFields["ColDate"].IsArchived)
	require.Len(t, item.Fields, 15)
	fields := item.Fields
	assert.Equal(t, "title", fields[0].Key)
	assert.Equal(t, "Tasks", fields[0].Text)
	require.Len(t, fields[0].RichText, 1)
	require.Len(t, fields[0].RichText[0].Elements, 1)
	section, ok := fields[0].RichText[0].Elements[0].(*RichTextSection)
	require.True(t, ok)
	require.Len(t, section.Elements, 1)
	text, ok := section.Elements[0].(*RichTextSectionTextElement)
	require.True(t, ok)
	assert.Equal(t, "Tasks", text.Text)
	require.NotNil(t, text.Style)
	assert.True(t, text.Style.Bold)
	for i, want := range []string{`"Tasks"`, `1.5`, `false`, `null`} {
		assert.JSONEq(t, want, string(fields[i].Value))
	}
	assert.Equal(t, []float64{0, 1.5}, fields[1].Number)
	require.NotNil(t, fields[2].Checkbox)
	assert.False(t, *fields[2].Checkbox)
	assert.Equal(t, []string{"2025-10-10"}, fields[4].Date)
	assert.Equal(t, []int64{-1}, fields[4].Timestamp)
	assert.Equal(t, []string{"U123"}, fields[5].User)
	assert.Equal(t, []string{"done"}, fields[6].Select)
	assert.Equal(t, []string{"F456"}, fields[7].Attachment)
	assert.Equal(t, []string{"C123"}, fields[8].Channel)
	assert.Equal(t, []string{"test@example.com"}, fields[9].Email)
	assert.Equal(t, []string{"+1-555-123-4567"}, fields[10].Phone)
	assert.Equal(t, []int{4}, fields[11].Rating)
	require.Len(t, fields[12].Message, 1)
	assert.Equal(t, "C123", fields[12].Message[0].ChannelID)
	assert.Equal(t, "1234567890.123456", fields[12].Message[0].Timestamp)
	assert.Equal(t, "1234567890.000001", fields[12].Message[0].ThreadTimestamp)
	require.Len(t, fields[13].Link, 1)
	link := fields[13].Link[0]
	assert.Equal(t, "https://docs.slack.dev/", link.OriginalURL)
	assert.Equal(t, "Docs", link.DisplayName)
	require.NotNil(t, link.DisplayAsURL)
	assert.False(t, *link.DisplayAsURL)
	assert.JSONEq(t, `{"title":"Docs"}`, string(link.Attachment))
	require.Len(t, fields[14].Reference, 4)
	assert.Equal(t, &ListItemFileReference{FileID: "F456"}, fields[14].Reference[0].File)
	assert.Equal(t, &ListItemMessage{ChannelID: "C123", Timestamp: "1234567890.123456", ThreadTimestamp: "1234567890.000001"}, fields[14].Reference[1].Message)
	assert.Equal(t, &ListItemRecordReference{ListID: "F123", RowID: "Rec456"}, fields[14].Reference[2].ListRecord)
	assert.Equal(t, &ListItemCanvasReference{FileID: "F789", SectionID: "temp:C:ABC"}, fields[14].Reference[3].CanvasSection)
	require.NotNil(t, response.List)
	assert.Equal(t, "Sprint Board", response.List.Title)
	assert.Equal(t, 7, response.List.ListLimits.RowCount)
	assert.Equal(t, "next-page", response.ResponseMetadata.Cursor)
}

func TestListMetadata(t *testing.T) {
	// From the slackLists.items.info response, with one column and view.
	const payload = `{
  "ok":true,
  "list":{
   "id":"F1234567","created":1758744341,"timestamp":1758744341,"name":"list","title":"Sprint Board",
   "mimetype":"application/vnd.slack-list","filetype":"list","pretty_type":"List","user":"W0AB1CDE2",
   "user_team":"E0AB1CD2E","editable":true,"size":0,"mode":"list","is_external":false,
   "external_type":"","is_public":false,"public_url_shared":false,"display_as_bot":false,"username":"",
   "list_metadata":{
    "schema":[{"id":"Col123","name":"Owner","key":"owner","type":"user","is_primary_column":false,"options":{"format":"multi_entity","default_value":null,"default_value_typed":{"user":["U123"],"channel":["C123"],"select":["done"]},"show_member_name":false,"notify_users":false,"precision":0,"date_format":"YYYY/MM/DD","emoji":":star:","emoji_team_id":"T123","max":5}}],
    "subtask_schema":[{"id":"ColTitle","name":"Title","key":"title","type":"text","is_primary_column":true}],
    "views":[{"id":"View123","name":"All items","type":"table","is_locked":false,"position":"1758744342","columns":[{"visible":true,"key":"owner","id":"Col123","position":"5000000000"}],"date_created":1758744342,"created_by":"U123","stick_column_left":false,"is_all_items_view":true,"default_view_key":"all_items","show_completed_items":false}],
    "integrations":[],"icon":"","description":"Sprint work",
    "description_blocks":[{"type":"rich_text","block_id":"XYcTZ","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Sprint work"}]}]}],
    "is_trial":false,"creation_source":{"type":"copy_from_list","reference_id":"Fsource"},"todo_mode":false,"default_view":"View123"
   },
   "list_limits":{"over_row_maximum":false,"row_count_limit":1000,"row_count":7,"archived_row_count":0,"over_column_maximum":false,"column_count":20,"column_count_limit":30,"over_view_maximum":false,"view_count":2,"view_count_limit":50,"max_attachments_per_cell":10},
   "url_private":"https://files.slack.com/list","url_private_download":"https://files.slack.com/download/list",
   "permalink":"https://example.slack.com/lists/T123/F1234567","permalink_public":"https://files.slack.com/public",
   "last_editor":"U123","list_csv_download_url":"https://files.slack.com/csv/list","updated":1758744347,
   "is_starred":false,"skipped_shares":true,"teams_shared_with":["E0AB1CD2E"],"is_restricted_sharing_enabled":false,
   "has_rich_preview":false,"file_access":"visible","access":"owner","org_or_workspace_access":"none","is_ai_suggested":false,"is_modified_by_ai":true
  },
  "record":{"id":"Rec123","fields":[]},"subtasks":[]
 }`
	api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	})
	response, err := api.GetListItemInfo("F1234567", "Rec123")
	require.NoError(t, err)
	list := response.List
	assert.Equal(t, "F1234567", list.ID)
	assert.Equal(t, JSONTime(1758744341), list.Created)
	assert.Equal(t, JSONTime(1758744341), list.Timestamp)
	assert.Equal(t, "list", list.Filetype)
	assert.Equal(t, "E0AB1CD2E", list.UserTeam)
	assert.Equal(t, "U123", list.LastEditor)
	assert.Equal(t, JSONTime(1758744347), list.Updated)
	assert.Equal(t, ListAccessOwner, list.Access)
	assert.Equal(t, "visible", list.FileAccess)
	assert.Equal(t, "none", list.OrgOrWorkspaceAccess)
	assert.True(t, list.SkippedShares)
	assert.True(t, list.IsModifiedByAI)
	assert.Equal(t, []string{"E0AB1CD2E"}, list.TeamsSharedWith)
	assert.Equal(t, "https://files.slack.com/csv/list", list.ListCSVDownloadURL)
	assert.Equal(t, ListLimits{RowCountLimit: 1000, RowCount: 7, ColumnCount: 20, ColumnCountLimit: 30, ViewCount: 2, ViewCountLimit: 50, MaxAttachmentsPerCell: 10}, list.ListLimits)
	metadata := list.ListMetadata
	assert.Equal(t, "Sprint work", metadata.Description)
	require.Len(t, metadata.DescriptionBlocks, 1)
	assert.Equal(t, "XYcTZ", metadata.DescriptionBlocks[0].BlockID)
	assert.Equal(t, &ListCreationSource{Type: "copy_from_list", ReferenceID: "Fsource"}, metadata.CreationSource)
	assert.Equal(t, "View123", metadata.DefaultView)
	require.Len(t, metadata.SubtaskSchema, 1)
	require.NotNil(t, metadata.SubtaskSchema[0].IsPrimaryColumn)
	assert.True(t, *metadata.SubtaskSchema[0].IsPrimaryColumn)
	require.Len(t, metadata.Schema, 1)
	options := metadata.Schema[0].Options
	require.NotNil(t, options)
	assert.Equal(t, ListColumnFormatMultiEntity, options.Format)
	assert.JSONEq(t, `null`, string(options.DefaultValue))
	assert.Equal(t, &ListColumnDefaultValue{User: []string{"U123"}, Channel: []string{"C123"}, Select: []string{"done"}}, options.DefaultValueTyped)
	require.NotNil(t, options.ShowMemberName)
	assert.False(t, *options.ShowMemberName)
	require.NotNil(t, options.NotifyUsers)
	assert.False(t, *options.NotifyUsers)
	require.NotNil(t, options.Precision)
	assert.Zero(t, *options.Precision)
	assert.Equal(t, ListDateFormatYearMonthDay, options.DateFormat)
	assert.Equal(t, ":star:", options.Emoji)
	assert.Equal(t, "T123", options.EmojiTeamID)
	require.NotNil(t, options.Max)
	assert.Equal(t, 5, *options.Max)
	require.Len(t, metadata.Views, 1)
	view := metadata.Views[0]
	assert.Equal(t, "1758744342", view.Position)
	assert.Equal(t, int64(1758744342), view.DateCreated)
	assert.Equal(t, "U123", view.CreatedBy)
	assert.True(t, view.IsAllItemsView)
	assert.Equal(t, "all_items", view.DefaultViewKey)
	require.NotNil(t, view.ShowCompletedItems)
	assert.False(t, *view.ShowCompletedItems)
	assert.Equal(t, []ListViewColumn{{Visible: true, Key: "owner", ID: "Col123", Position: "5000000000"}}, view.Columns)
}
