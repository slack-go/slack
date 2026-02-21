package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestEntityPresentDetailsParameters(t *testing.T) {
	// Test basic parameter structure
	params := EntityPresentDetailsParameters{
		TriggerID: "1234567890123.1234567890123.abcdef01234567890abcdef012345689",
		Metadata: &EntityDetailsMetadata{
			EntityType: "slack#/entities/file",
			URL:        "https://example.com/document/123",
			ExternalRef: WorkObjectExternalRef{
				ID:   "123",
				Type: "document",
			},
			EntityPayload: map[string]interface{}{
				"title":       "Test Document",
				"description": "A test document for Work Objects",
				"status":      "active",
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Errorf("Failed to marshal EntityPresentDetailsParameters: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled EntityPresentDetailsParameters
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal EntityPresentDetailsParameters: %v", err)
	}

	// Verify the data
	if unmarshaled.TriggerID != params.TriggerID {
		t.Errorf("Expected trigger_id '%s', got '%s'", params.TriggerID, unmarshaled.TriggerID)
	}

	if unmarshaled.Metadata.EntityType != params.Metadata.EntityType {
		t.Errorf("Expected entity_type '%s', got '%s'", params.Metadata.EntityType, unmarshaled.Metadata.EntityType)
	}

	if unmarshaled.Metadata.ExternalRef.ID != params.Metadata.ExternalRef.ID {
		t.Errorf("Expected external_ref.id '%s', got '%s'", params.Metadata.ExternalRef.ID, unmarshaled.Metadata.ExternalRef.ID)
	}
}

func TestEntityDetailsError(t *testing.T) {
	// Test error structure
	errorObj := EntityDetailsError{
		Status:        "restricted",
		CustomTitle:   "Access Denied",
		CustomMessage: "You do not have permission to view this entity.",
		MessageFormat: "markdown",
		Actions: []EntityDetailsAction{
			{
				Text:     "Request Access",
				ActionID: "request_access",
				Value:    "entity_123",
				Style:    "primary",
				URL:      "https://example.com/request-access",
				ProcessingState: &EntityDetailsProcessingState{
					Enabled: true,
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(errorObj)
	if err != nil {
		t.Errorf("Failed to marshal EntityDetailsError: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled EntityDetailsError
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal EntityDetailsError: %v", err)
	}

	// Verify the data
	if unmarshaled.Status != errorObj.Status {
		t.Errorf("Expected status '%s', got '%s'", errorObj.Status, unmarshaled.Status)
	}

	if unmarshaled.CustomTitle != errorObj.CustomTitle {
		t.Errorf("Expected custom_title '%s', got '%s'", errorObj.CustomTitle, unmarshaled.CustomTitle)
	}

	if len(unmarshaled.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(unmarshaled.Actions))
	}

	if len(unmarshaled.Actions) > 0 {
		action := unmarshaled.Actions[0]
		if action.Text != "Request Access" {
			t.Errorf("Expected action text 'Request Access', got '%s'", action.Text)
		}
		if action.ProcessingState == nil || !action.ProcessingState.Enabled {
			t.Error("Expected processing state to be enabled")
		}
	}
}

func TestEntityPresentDetailsWithMetadata(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/entity.presentDetails", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		// Verify the request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("Failed to parse form: %v", err)
			return
		}

		// Check required fields
		triggerID := r.FormValue("trigger_id")
		if triggerID != "1234567890123.1234567890123.abcdef01234567890abcdef012345689" {
			t.Errorf("Expected trigger_id '1234567890123.1234567890123.abcdef01234567890abcdef012345689', got '%s'", triggerID)
		}

		// Check metadata
		metadataStr := r.FormValue("metadata")
		if metadataStr == "" {
			t.Error("Expected metadata to be present")
		} else {
			var metadata EntityDetailsMetadata
			err := json.Unmarshal([]byte(metadataStr), &metadata)
			if err != nil {
				t.Errorf("Failed to unmarshal metadata: %v", err)
			}
			if metadata.EntityType != "slack#/entities/file" {
				t.Errorf("Expected entity_type 'slack#/entities/file', got '%s'", metadata.EntityType)
			}
		}

		// Return success response
		response := EntityPresentDetailsResponse{
			SlackResponse: SlackResponse{Ok: true},
		}
		jsonResponse, _ := json.Marshal(response)
		rw.Write(jsonResponse)
	})

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	metadata := EntityDetailsMetadata{
		EntityType: "slack#/entities/file",
		URL:        "https://example.com/document/123",
		ExternalRef: WorkObjectExternalRef{
			ID:   "123",
			Type: "document",
		},
		EntityPayload: map[string]interface{}{
			"title":       "Test Document",
			"description": "A test document for Work Objects",
		},
	}

	err := api.EntityPresentDetailsWithMetadata("1234567890123.1234567890123.abcdef01234567890abcdef012345689", metadata)
	if err != nil {
		t.Errorf("EntityPresentDetailsWithMetadata returned error: %v", err)
	}
}

func TestEntityPresentDetailsWithError(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/entity.presentDetails", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("Failed to parse form: %v", err)
			return
		}

		// Check error field
		errorStr := r.FormValue("error")
		if errorStr == "" {
			t.Error("Expected error to be present")
		} else {
			var errorObj EntityDetailsError
			err := json.Unmarshal([]byte(errorStr), &errorObj)
			if err != nil {
				t.Errorf("Failed to unmarshal error: %v", err)
			}
			if errorObj.Status != "restricted" {
				t.Errorf("Expected error status 'restricted', got '%s'", errorObj.Status)
			}
		}

		// Return success response
		response := EntityPresentDetailsResponse{
			SlackResponse: SlackResponse{Ok: true},
		}
		jsonResponse, _ := json.Marshal(response)
		rw.Write(jsonResponse)
	})

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	errorObj := EntityDetailsError{
		Status:        "restricted",
		CustomTitle:   "Access Denied",
		CustomMessage: "You do not have permission to view this entity.",
		Actions: []EntityDetailsAction{
			{
				Text:     "Request Access",
				ActionID: "request_access",
				Style:    "primary",
			},
		},
	}

	err := api.EntityPresentDetailsWithError("1234567890123.1234567890123.abcdef01234567890abcdef012345689", errorObj)
	if err != nil {
		t.Errorf("EntityPresentDetailsWithError returned error: %v", err)
	}
}

func TestEntityPresentDetailsWithAuth(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/entity.presentDetails", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("Failed to parse form: %v", err)
			return
		}

		// Check auth fields
		userAuthRequired := r.FormValue("user_auth_required")
		if userAuthRequired != "true" {
			t.Errorf("Expected user_auth_required 'true', got '%s'", userAuthRequired)
		}

		userAuthURL := r.FormValue("user_auth_url")
		if userAuthURL != "https://example.com/auth" {
			t.Errorf("Expected user_auth_url 'https://example.com/auth', got '%s'", userAuthURL)
		}

		userAuthMessage := r.FormValue("user_auth_message")
		if userAuthMessage != "Please authenticate to view this entity." {
			t.Errorf("Expected user_auth_message 'Please authenticate to view this entity.', got '%s'", userAuthMessage)
		}

		// Return success response
		response := EntityPresentDetailsResponse{
			SlackResponse: SlackResponse{Ok: true},
		}
		jsonResponse, _ := json.Marshal(response)
		rw.Write(jsonResponse)
	})

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.EntityPresentDetailsWithAuth(
		"1234567890123.1234567890123.abcdef01234567890abcdef012345689",
		"https://example.com/auth",
		"Please authenticate to view this entity.",
	)
	if err != nil {
		t.Errorf("EntityPresentDetailsWithAuth returned error: %v", err)
	}
}

func TestEntityPresentDetailsContext(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/entity.presentDetails", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		// Return success response
		response := EntityPresentDetailsResponse{
			SlackResponse: SlackResponse{Ok: true},
		}
		jsonResponse, _ := json.Marshal(response)
		rw.Write(jsonResponse)
	})

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := EntityPresentDetailsParameters{
		TriggerID: "1234567890123.1234567890123.abcdef01234567890abcdef012345689",
		Metadata: &EntityDetailsMetadata{
			EntityType: "slack#/entities/task",
			URL:        "https://example.com/task/456",
			ExternalRef: WorkObjectExternalRef{
				ID: "456",
			},
			EntityPayload: map[string]interface{}{
				"title":  "Test Task",
				"status": "in_progress",
			},
		},
	}

	err := api.EntityPresentDetails(params)
	if err != nil {
		t.Errorf("EntityPresentDetails returned error: %v", err)
	}
}
