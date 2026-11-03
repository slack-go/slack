package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestManifestUserConnections(t *testing.T) {
	for _, payload := range []string{
		`{"is_enabled":true,"auth_hosts":["example.com"]}`,
		`{"is_enabled":false,"auth_hosts":["example.com","other.example"]}`,
		`{"is_enabled":false,"auth_hosts":[]}`,
	} {
		t.Run(payload, func(t *testing.T) {
			var manifest Manifest
			require.NoError(t, json.Unmarshal([]byte(`{"features":{"user_connections":`+payload+`}}`), &manifest))
			require.NotNil(t, manifest.Features.UserConnections)
			encoded, err := json.Marshal(manifest.Features.UserConnections)
			require.NoError(t, err)
			assert.JSONEq(t, payload, string(encoded))
			yamlBytes, err := yaml.Marshal(manifest)
			require.NoError(t, err)
			var fromYAML Manifest
			require.NoError(t, yaml.Unmarshal(yamlBytes, &fromYAML))
			assert.Equal(t, manifest.Features.UserConnections, fromYAML.Features.UserConnections)
		})
	}
	encoded, err := json.Marshal(Features{})
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "user_connections")
}
