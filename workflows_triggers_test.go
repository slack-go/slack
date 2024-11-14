package slack

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

type workflowsHandler struct {
	rawResponse string
}

func (h *workflowsHandler) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.rawResponse))
}

func TestSlack_WorkflowsTriggersPermissionsAdd(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName     string
		input        *WorkflowsTriggersPermissionsAddInput
		rawResp      string
		expectedResp *WorkflowsTriggersPermissionsAddOutput
		expectedErr  error
	}{
		{
			caseName: "success",
			input: &WorkflowsTriggersPermissionsAddInput{
				TriggerId:  "Ft0000000001",
				ChannelIds: []string{"C0000000001"},
				OrgIds:     []string{"E00000001"},
				TeamIds:    []string{"T0000000001"},
				UserIds:    []string{"U0000000001", "U0000000002"},
			},
			rawResp: `{
				"ok": true,
				"permission_type": "named_entities",
				"user_ids": ["U0000000001", "U0000000002"],
				"channel_ids": ["C0000000001"],
				"org_ids": ["E00000001"],
				"team_ids": ["T0000000001"]
			}`,
			expectedResp: &WorkflowsTriggersPermissionsAddOutput{
				PermissionType: "named_entities",
				UserIds:        []string{"U0000000001", "U0000000002"},
				ChannelIds:     []string{"C0000000001"},
				OrgIds:         []string{"E00000001"},
				TeamIds:        []string{"T0000000001"},
			},
			expectedErr: nil,
		},
	}

	h := &workflowsHandler{}
	http.HandleFunc("/workflows.triggers.permissions.add", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			resp, err := api.WorkflowsTriggersPermissionsAdd(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
			if resp == nil || c.expectedResp == nil {
				return
			}
			if !reflect.DeepEqual(resp, c.expectedResp) {
				t.Fatalf("expected:\n\t%v\n but got:\n\t%v\n", c.expectedResp, resp)
			}
		})
	}
}

func TestSlack_WorkflowsTriggersPermissionsList(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName     string
		input        *WorkflowsTriggersPermissionsListInput
		rawResp      string
		expectedResp *WorkflowsTriggersPermissionsListOutput
		expectedErr  error
	}{
		{
			caseName: "success",
			input: &WorkflowsTriggersPermissionsListInput{
				TriggerId: "Ft0000000001",
			},
			rawResp: `{
				"ok": true,
				"permission_type": "named_entities",
				"user_ids": ["U0000000001", "U0000000002"],
				"channel_ids": ["C0000000001"],
				"org_ids": ["E00000001"],
				"team_ids": ["T0000000001"]
			}`,
			expectedResp: &WorkflowsTriggersPermissionsListOutput{
				PermissionType: "named_entities",
				UserIds:        []string{"U0000000001", "U0000000002"},
				ChannelIds:     []string{"C0000000001"},
				OrgIds:         []string{"E00000001"},
				TeamIds:        []string{"T0000000001"},
			},
			expectedErr: nil,
		},
	}

	h := &workflowsHandler{}
	http.HandleFunc("/workflows.triggers.permissions.list", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			resp, err := api.WorkflowsTriggersPermissionsList(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
			if resp == nil || c.expectedResp == nil {
				return
			}
			if !reflect.DeepEqual(resp, c.expectedResp) {
				t.Fatalf("expected:\n\t%v\n but got:\n\t%v\n", c.expectedResp, resp)
			}
		})
	}
}

func TestSlack_WorkflowsTriggersPermissionsRemove(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName     string
		input        *WorkflowsTriggersPermissionsRemoveInput
		rawResp      string
		expectedResp *WorkflowsTriggersPermissionsRemoveOutput
		expectedErr  error
	}{
		{
			caseName: "success",
			input: &WorkflowsTriggersPermissionsRemoveInput{
				TriggerId:  "Ft0000000001",
				ChannelIds: []string{"C0000000001"},
				OrgIds:     []string{"E00000001"},
				TeamIds:    []string{"T0000000001"},
				UserIds:    []string{"U0000000001", "U0000000002"},
			},
			rawResp: `{
				"ok": true,
				"permission_type": "named_entities",
				"user_ids": ["U0000000001", "U0000000002"],
				"channel_ids": ["C0000000001"],
				"org_ids": ["E00000001"],
				"team_ids": ["T0000000001"]
			}`,
			expectedResp: &WorkflowsTriggersPermissionsRemoveOutput{
				PermissionType: "named_entities",
				UserIds:        []string{"U0000000001", "U0000000002"},
				ChannelIds:     []string{"C0000000001"},
				OrgIds:         []string{"E00000001"},
				TeamIds:        []string{"T0000000001"},
			},
			expectedErr: nil,
		},
	}

	h := &workflowsHandler{}
	http.HandleFunc("/workflows.triggers.permissions.remove", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			resp, err := api.WorkflowsTriggersPermissionsRemove(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
			if resp == nil || c.expectedResp == nil {
				return
			}
			if !reflect.DeepEqual(resp, c.expectedResp) {
				t.Fatalf("expected:\n\t%v\n but got:\n\t%v\n", c.expectedResp, resp)
			}
		})
	}
}

func TestSlack_WorkflowsTriggersPermissionsSet(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName     string
		input        *WorkflowsTriggersPermissionsSetInput
		rawResp      string
		expectedResp *WorkflowsTriggersPermissionsSetOutput
		expectedErr  error
	}{
		{
			caseName: "success",
			input: &WorkflowsTriggersPermissionsSetInput{
				PermissionType: "named_entities",
				TriggerId:      "Ft0000000001",
				ChannelIds:     []string{"C0000000001"},
				OrgIds:         []string{"E00000001"},
				TeamIds:        []string{"T0000000001"},
				UserIds:        []string{"U0000000001", "U0000000002"},
			},
			rawResp: `{
				"ok": true,
				"permission_type": "named_entities",
				"user_ids": ["U0000000001", "U0000000002"],
				"channel_ids": ["C0000000001"],
				"org_ids": ["E00000001"],
				"team_ids": ["T0000000001"]
			}`,
			expectedResp: &WorkflowsTriggersPermissionsSetOutput{
				PermissionType: "named_entities",
				UserIds:        []string{"U0000000001", "U0000000002"},
				ChannelIds:     []string{"C0000000001"},
				OrgIds:         []string{"E00000001"},
				TeamIds:        []string{"T0000000001"},
			},
			expectedErr: nil,
		},
	}

	h := &workflowsHandler{}
	http.HandleFunc("/workflows.triggers.permissions.set", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			resp, err := api.WorkflowsTriggersPermissionsSet(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
			if resp == nil || c.expectedResp == nil {
				return
			}
			if !reflect.DeepEqual(resp, c.expectedResp) {
				t.Fatalf("expected:\n\t%v\n but got:\n\t%v\n", c.expectedResp, resp)
			}
		})
	}
}
