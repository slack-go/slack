package slack

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type listMethodTest struct {
	method      string
	form        url.Values
	response    string
	call        func(*Client) (any, error)
	callContext func(*Client, context.Context) (any, error)
}

func listMethodTests() []listMethodTest {
	return []listMethodTest{
		{
			method: "slackLists.create", form: url.Values{"name": {"My List"}},
			response:    `{"ok":true,"list_id":"F1234ABCD","list_metadata":{"schema":[{"key":"name","name":"Name","is_primary_column":true,"type":"text","id":"Col012A3BCDE4"}]}}`,
			call:        func(api *Client) (any, error) { return api.CreateList("My List", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) { return api.CreateListContext(ctx, "My List", nil) },
		},
		{
			method: "slackLists.update", form: url.Values{"id": {"F1234ABCD"}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) { return nil, api.UpdateList("F1234ABCD", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.UpdateListContext(ctx, "F1234ABCD", nil)
			},
		},
		{
			method: "slackLists.access.set", form: url.Values{"list_id": {"F1234ABCD"}, "access_level": {"read"}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) { return nil, api.SetListAccess("F1234ABCD", ListAccessRead, nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.SetListAccessContext(ctx, "F1234ABCD", ListAccessRead, nil)
			},
		},
		{
			method: "slackLists.access.delete", form: url.Values{"list_id": {"F1234ABCD"}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) { return nil, api.DeleteListAccess("F1234ABCD", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.DeleteListAccessContext(ctx, "F1234ABCD", nil)
			},
		},
		{
			method: "slackLists.download.start", form: url.Values{"list_id": {"F1234ABCD"}}, response: `{"ok":true,"job_id":"LeF1234567"}`,
			call: func(api *Client) (any, error) { return api.StartListDownload("F1234ABCD", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return api.StartListDownloadContext(ctx, "F1234ABCD", nil)
			},
		},
		{
			method: "slackLists.download.get", form: url.Values{"list_id": {"F1234ABCD"}, "job_id": {"LeF1234567"}},
			response: `{"ok":true,"status":"COMPLETED","download_url":"https://files.slack.com/export.csv"}`,
			call:     func(api *Client) (any, error) { return api.GetListDownload("F1234ABCD", "LeF1234567", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return api.GetListDownloadContext(ctx, "F1234ABCD", "LeF1234567", nil)
			},
		},
		{
			method: "slackLists.items.create", form: url.Values{"list_id": {"F1234ABCD"}},
			response: `{"ok":true,"item":{"id":"Rec018ALE9718","list_id":"F1234ABCD","date_created":1758744345,"created_by":"W0AB1CDE2","updated_by":"W0AB1CDE2","fields":[],"updated_timestamp":"1758744345"}}`,
			call:     func(api *Client) (any, error) { return api.CreateListItem("F1234ABCD", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return api.CreateListItemContext(ctx, "F1234ABCD", nil)
			},
		},
		{
			method: "slackLists.items.delete", form: url.Values{"list_id": {"F1234ABCD"}, "id": {"Rec1234567"}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) { return nil, api.DeleteListItem("F1234ABCD", "Rec1234567") },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.DeleteListItemContext(ctx, "F1234ABCD", "Rec1234567")
			},
		},
		{
			method: "slackLists.items.deleteMultiple", form: url.Values{"list_id": {"F1234ABCD"}, "ids": {`["Rec1234567","Rec2345678"]`}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) {
				return nil, api.DeleteListItems("F1234ABCD", []string{"Rec1234567", "Rec2345678"})
			},
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.DeleteListItemsContext(ctx, "F1234ABCD", []string{"Rec1234567", "Rec2345678"})
			},
		},
		{
			method: "slackLists.items.info", form: url.Values{"list_id": {"F1234ABCD"}, "id": {"Rec1234567"}},
			response: `{"ok":true,"list":{"id":"F1234ABCD","title":"Sprint Board"},"record":{"id":"Rec1234567","fields":[]},"subtasks":[{"id":"Rec2345678","parent_record_id":"Rec1234567","is_subscribed":false,"fields":[]}]}`,
			call:     func(api *Client) (any, error) { return api.GetListItemInfo("F1234ABCD", "Rec1234567", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return api.GetListItemInfoContext(ctx, "F1234ABCD", "Rec1234567", nil)
			},
		},
		{
			method: "slackLists.items.list", form: url.Values{"list_id": {"F1234ABCD"}},
			response: `{"ok":true,"items":[{"id":"Rec1234567","fields":[]}],"response_metadata":{"next_cursor":"cursor-2"}}`,
			call:     func(api *Client) (any, error) { return api.GetListItems("F1234ABCD", nil) },
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return api.GetListItemsContext(ctx, "F1234ABCD", nil)
			},
		},
		{
			method: "slackLists.items.update", form: url.Values{"list_id": {"F1234ABCD"}, "cells": {`[{"column_id":"Col123","user":["U123"],"row_id":"Rec1234567"}]`}}, response: `{"ok":true}`,
			call: func(api *Client) (any, error) {
				return nil, api.UpdateListItems("F1234ABCD", []ListItemCell{{RowID: "Rec1234567", ListItemField: ListItemField{ColumnID: "Col123", User: []string{"U123"}}}})
			},
			callContext: func(api *Client, ctx context.Context) (any, error) {
				return nil, api.UpdateListItemsContext(ctx, "F1234ABCD", []ListItemCell{{RowID: "Rec1234567", ListItemField: ListItemField{ColumnID: "Col123", User: []string{"U123"}}}})
			},
		},
	}
}

func newListTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return New("testing-token", OptionAPIURL(server.URL+"/"))
}

func TestListMethods(t *testing.T) {
	for _, tt := range listMethodTests() {
		t.Run(tt.method, func(t *testing.T) {
			api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/"+tt.method, r.URL.Path)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				r.PostForm.Del("token")
				assert.Equal(t, tt.form, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.response))
			})
			for _, call := range []func() (any, error){func() (any, error) { return tt.call(api) }, func() (any, error) { return tt.callContext(api, context.Background()) }} {
				result, err := call()
				require.NoError(t, err)
				switch response := result.(type) {
				case *CreateListResponse:
					assert.Equal(t, "F1234ABCD", response.ListID)
					require.Len(t, response.ListMetadata.Schema, 1)
					assert.Equal(t, "Col012A3BCDE4", response.ListMetadata.Schema[0].ID)
					assert.Equal(t, ListColumnText, response.ListMetadata.Schema[0].Type)
				case *StartListDownloadResponse:
					assert.Equal(t, "LeF1234567", response.JobID)
				case *GetListDownloadResponse:
					assert.Equal(t, "COMPLETED", response.Status)
					assert.Equal(t, "https://files.slack.com/export.csv", response.DownloadURL)
				case *CreateListItemResponse:
					assert.Equal(t, "Rec018ALE9718", response.Item.ID)
					assert.Equal(t, int64(1758744345), response.Item.DateCreated)
					assert.Equal(t, "1758744345", response.Item.UpdatedTimestamp)
				case *GetListItemInfoResponse:
					assert.Equal(t, "F1234ABCD", response.List.ID)
					assert.Equal(t, "Sprint Board", response.List.Title)
					assert.Equal(t, "Rec1234567", response.Record.ID)
					require.Len(t, response.Subtasks, 1)
					assert.Equal(t, "Rec1234567", response.Subtasks[0].ParentRecordID)
					require.NotNil(t, response.Subtasks[0].IsSubscribed)
					assert.False(t, *response.Subtasks[0].IsSubscribed)
				case *GetListItemsResponse:
					require.Len(t, response.Items, 1)
					assert.Equal(t, "Rec1234567", response.Items[0].ID)
					assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
					assert.Nil(t, response.List)
				}
			}
		})
	}
}

func TestListMethodErrors(t *testing.T) {
	for _, method := range listMethodTests() {
		t.Run(method.method, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				status int
				body   string
			}{
				{"slack error", http.StatusOK, `{"ok":false,"error":"list_not_found","response_metadata":{"messages":["List is unavailable"]}}`},
				{"invalid JSON", http.StatusOK, `{"ok":`},
				{"HTTP error", http.StatusInternalServerError, `server error`},
				{"rate limit", http.StatusTooManyRequests, `rate limited`},
			} {
				t.Run(tt.name, func(t *testing.T) {
					api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.Header().Set("Retry-After", "3")
						w.WriteHeader(tt.status)
						_, _ = w.Write([]byte(tt.body))
					})
					result, err := method.callContext(api, context.Background())
					require.Error(t, err)
					switch tt.name {
					case "slack error":
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, "list_not_found", slackErr.Err)
						assert.Equal(t, []string{"List is unavailable"}, slackErr.ResponseMetadata.Messages)
						if response, ok := result.(interface{ Err() error }); ok {
							require.NotNil(t, response)
							assert.EqualError(t, response.Err(), "list_not_found")
						}
					case "rate limit":
						var rateErr *RateLimitedError
						require.ErrorAs(t, err, &rateErr)
						assert.Equal(t, 3*time.Second, rateErr.RetryAfter)
					default:
						assert.Nil(t, result)
					}
				})
			}
			t.Run("canceled context", func(t *testing.T) {
				api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("canceled request reached server") })
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := method.callContext(api, ctx)
				require.ErrorIs(t, err, context.Canceled)
			})
		})
	}
}

func TestListOptions(t *testing.T) {
	falseValue, trueValue, zero := false, true, 0
	description := []RichTextBlock{*NewRichTextBlock("", NewRichTextSection(NewRichTextSectionTextElement("Tasks", nil)))}
	for _, tt := range []struct {
		name       string
		method     string
		form       url.Values
		jsonFields map[string]string
		call       func(*Client) error
	}{
		{
			name: "create", method: "slackLists.create",
			form: url.Values{"name": {"Tasks"}, "copy_from_list_id": {"Fsource"}, "include_copied_list_records": {"false"}, "todo_mode": {"true"}},
			jsonFields: map[string]string{
				"description_blocks": `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Tasks"}]}]}]`,
				"schema":             `[{"key":"estimate","name":"Estimate","type":"number","is_primary_column":false,"options":{"precision":0,"show_member_name":false,"notify_users":false}},{"key":"status","name":"Status","type":"select","options":{"choices":[{"value":"done","label":"Done","color":"green"}],"format":"single_select","default_value_typed":{"select":["done"]}}}]`,
			},
			call: func(api *Client) error {
				_, err := api.CreateList("Tasks", CreateListOptionDescriptionBlocks(description), CreateListOptionCopyFromListID("Fsource"), CreateListOptionIncludeCopiedListRecords(false), CreateListOptionTodoMode(true), CreateListOptionSchema([]ListColumn{
					{Key: "estimate", Name: "Estimate", Type: ListColumnNumber, IsPrimaryColumn: &falseValue, Options: &ListColumnOptions{Precision: &zero, ShowMemberName: &falseValue, NotifyUsers: &falseValue}},
					{Key: "status", Name: "Status", Type: ListColumnSelect, Options: &ListColumnOptions{Choices: []ListColumnChoice{{Value: "done", Label: "Done", Color: ListChoiceGreen}}, Format: ListColumnFormatSingleSelect, DefaultValueTyped: &ListColumnDefaultValue{Select: []string{"done"}}}},
				}))
				return err
			},
		},
		{
			name: "update clears description", method: "slackLists.update",
			form: url.Values{"id": {"F123"}, "name": {"New name"}, "todo_mode": {"false"}}, jsonFields: map[string]string{"description_blocks": `[]`},
			call: func(api *Client) error {
				return api.UpdateList("F123", UpdateListOptionName("New name"), UpdateListOptionDescriptionBlocks([]RichTextBlock{}), UpdateListOptionTodoMode(false))
			},
		},
		{
			name: "set channel access", method: "slackLists.access.set", form: url.Values{"list_id": {"F123"}, "access_level": {"write"}}, jsonFields: map[string]string{"channel_ids": `["C123","C456"]`},
			call: func(api *Client) error {
				return api.SetListAccess("F123", ListAccessWrite, SetListAccessOptionChannelIDs([]string{"C123", "C456"}))
			},
		},
		{
			name: "set owner", method: "slackLists.access.set", form: url.Values{"list_id": {"F123"}, "access_level": {"owner"}}, jsonFields: map[string]string{"user_ids": `["U123"]`},
			call: func(api *Client) error {
				return api.SetListAccess("F123", ListAccessOwner, SetListAccessOptionUserIDs([]string{"U123"}))
			},
		},
		{
			name: "delete channel access", method: "slackLists.access.delete", form: url.Values{"list_id": {"F123"}}, jsonFields: map[string]string{"channel_ids": `["C123"]`},
			call: func(api *Client) error {
				return api.DeleteListAccess("F123", DeleteListAccessOptionChannelIDs([]string{"C123"}))
			},
		},
		{
			name: "delete user access", method: "slackLists.access.delete", form: url.Values{"list_id": {"F123"}}, jsonFields: map[string]string{"user_ids": `["U123"]`},
			call: func(api *Client) error {
				return api.DeleteListAccess("F123", DeleteListAccessOptionUserIDs([]string{"U123"}))
			},
		},
		{
			name: "start JSON export", method: "slackLists.download.start", form: url.Values{"list_id": {"F123"}, "format": {"json"}, "include_archived": {"false"}, "include_threads": {"true"}, "include_attachments": {"false"}},
			call: func(api *Client) error {
				_, err := api.StartListDownload("F123", StartListDownloadOptionFormat(ListDownloadJSON), StartListDownloadOptionIncludeArchived(false), StartListDownloadOptionIncludeThreads(true), StartListDownloadOptionIncludeAttachments(false))
				return err
			},
		},
		{
			name: "get JSON export", method: "slackLists.download.get", form: url.Values{"list_id": {"F123"}, "job_id": {"Le123"}, "format": {"json"}, "include_threads": {"true"}, "include_attachments": {"false"}},
			call: func(api *Client) error {
				_, err := api.GetListDownload("F123", "Le123", GetListDownloadOptionFormat(ListDownloadJSON), GetListDownloadOptionIncludeThreads(true), GetListDownloadOptionIncludeAttachments(false))
				return err
			},
		},
		{
			name: "create subtask", method: "slackLists.items.create", form: url.Values{"list_id": {"F123"}, "parent_item_id": {"RecParent"}}, jsonFields: map[string]string{"initial_fields": `[{"column_id":"Col123","checkbox":false}]`},
			call: func(api *Client) error {
				_, err := api.CreateListItem("F123", CreateListItemOptionParentItemID("RecParent"), CreateListItemOptionInitialFields([]ListItemField{{ColumnID: "Col123", Checkbox: &falseValue}}))
				return err
			},
		},
		{
			name: "duplicate item", method: "slackLists.items.create", form: url.Values{"list_id": {"F123"}, "duplicated_item_id": {"RecOriginal"}},
			call: func(api *Client) error {
				_, err := api.CreateListItem("F123", CreateListItemOptionDuplicatedItemID("RecOriginal"))
				return err
			},
		},
		{
			name: "subscription", method: "slackLists.items.info", form: url.Values{"list_id": {"F123"}, "id": {"Rec123"}, "include_is_subscribed": {"true"}},
			call: func(api *Client) error {
				_, err := api.GetListItemInfo("F123", "Rec123", GetListItemInfoOptionIncludeIsSubscribed(true))
				return err
			},
		},
		{
			name: "pagination", method: "slackLists.items.list", form: url.Values{"list_id": {"F123"}, "limit": {"100"}, "cursor": {"cursor-2"}, "archived": {"false"}, "include_list": {"true"}},
			call: func(api *Client) error {
				_, err := api.GetListItems("F123", GetListItemsOptionLimit(100), GetListItemsOptionCursor("cursor-2"), GetListItemsOptionArchived(false), GetListItemsOptionIncludeList(true))
				return err
			},
		},
		{
			name: "update and clear cells", method: "slackLists.items.update", form: url.Values{"list_id": {"F123"}},
			jsonFields: map[string]string{"cells": `[{"row_id":"Rec123","column_id":"ColUser","user":[]},{"row_id":"Rec123","column_id":"ColDone","checkbox":false},{"row_id":"RecNew","row_id_to_create":true,"column_id":"ColDone","checkbox":true}]`},
			call: func(api *Client) error {
				return api.UpdateListItems("F123", []ListItemCell{
					{RowID: "Rec123", ListItemField: ListItemField{ColumnID: "ColUser", User: []string{}}},
					{RowID: "Rec123", ListItemField: ListItemField{ColumnID: "ColDone", Checkbox: &falseValue}},
					{RowID: "RecNew", RowIDToCreate: true, ListItemField: ListItemField{ColumnID: "ColDone", Checkbox: &trueValue}},
				})
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/"+tt.method, r.URL.Path)
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				r.PostForm.Del("token")
				for field, want := range tt.jsonFields {
					assert.JSONEq(t, want, r.PostForm.Get(field), field)
					r.PostForm.Del(field)
				}
				assert.Equal(t, tt.form, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			})
			require.NoError(t, tt.call(api))
		})
	}
}

func TestListMarshalErrors(t *testing.T) {
	for _, tt := range []struct {
		name string
		call func(*Client) error
	}{
		{"schema", func(api *Client) error {
			_, err := api.CreateList("Tasks", CreateListOptionSchema([]ListColumn{{Options: &ListColumnOptions{DefaultValue: json.RawMessage(`invalid`)}}}))
			return err
		}},
		{"initial fields", func(api *Client) error {
			_, err := api.CreateListItem("F123", CreateListItemOptionInitialFields([]ListItemField{{ColumnID: "Col123", Number: []float64{math.NaN()}}}))
			return err
		}},
		{"update cells", func(api *Client) error {
			return api.UpdateListItems("F123", []ListItemCell{{RowID: "Rec123", ListItemField: ListItemField{ColumnID: "Col123", Number: []float64{math.Inf(1)}}}})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			api := newListTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("invalid JSON reached server") })
			err := tt.call(api)
			require.Error(t, err)
		})
	}
}
