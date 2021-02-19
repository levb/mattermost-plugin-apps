// +build !e2e

package store

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"

	"github.com/mattermost/mattermost-plugin-apps/awsclient/mock_awsclient"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/mock/mock_aws"
)

func TestManifestInit(t *testing.T) {
	for _, tc := range []struct {
		name        string
		data        string
		expectAWS   func(*mock_awsclient.MockClient)
		expectError string
	}{
		{
			name: "simple s3 happy",
			data: ` { "jira" : "v3.0.0" }`,
			expectAWS: func(awsClient *mock_awsclient.MockClient) {
				awsClient.EXPECT().GetS3(
					gomock.Eq("manifestbucket"), gomock.Eq("manifest_jira_v3.0.0")).Return(
					[]byte(`{"app_id":"jira", "type":"http"}`), nil).Times(1)
			},
		},
		{
			name: "s3 prefix happy",
			data: ` { "jira" : "s3:v3.0.0" }`,
			expectAWS: func(awsClient *mock_awsclient.MockClient) {
				awsClient.EXPECT().GetS3(
					gomock.Eq("manifestbucket"), gomock.Eq("manifest_jira_v3.0.0")).Return(
					[]byte(`{"app_id":"jira", "type":"http"}`), nil).Times(1)
			},
		},
		{
			name: "s3 multiple happy",
			data: ` { "jira" : "v3.0.0", "jira301" : "s3:v3.0.1" }`,
			expectAWS: func(awsClient *mock_awsclient.MockClient) {
				awsClient.EXPECT().GetS3(
					gomock.Eq("manifestbucket"), gomock.Eq("manifest_jira_v3.0.0")).Return(
					[]byte(`{"app_id":"jira", "type":"http"}`), nil).Times(1)
				awsClient.EXPECT().GetS3(
					gomock.Eq("manifestbucket"), gomock.Eq("manifest_jira301_v3.0.1")).Return(
					[]byte(`{"app_id":"jira301", "type":"http"}`), nil).Times(1)
			},
		},
		{
			name:        "file happy",
			data:        ` { "jira301" : "file:test-does-not-exist" }`,
			expectError: "failed to load global manifest for jira301: open /testassets/test-does-not-exist: no such file or directory",
		},
		{
			name:        "https happy",
			data:        ` { "jira302" : "https://host.test/file302" }`,
			expectError: `failed to load global manifest for jira302: Get "https://host.test/file302": dial tcp: lookup host.test: no such host`,
		},
		{
			name:        "invalid",
			data:        ` { "jira302" : "invalid:does:not::matter" }`,
			expectError: "failed to load global manifest for jira302: invalid:does:not::matter is invalid",
		},
		{
			name: "invalid app type",
			data: ` { "jira" : "v3.0.0" }`,
			expectAWS: func(awsClient *mock_awsclient.MockClient) {
				awsClient.EXPECT().GetS3(
					gomock.Eq("manifestbucket"), gomock.Eq("manifest_jira_v3.0.0")).Return(
					[]byte(`{"app_id":"jira", "type":"invalid"}`), nil).Times(1)
			},
			expectError: "invalid type: invalid",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pluginAPI := &plugintest.API{}
			mm := pluginapi.NewClient(pluginAPI)
			awsClient := mock_awsclient.NewMockClient(ctrl)
			aws := mock_aws.NewMockService(ctrl)
			conf := config.NewTestConfigurator(&config.Config{
				StoredConfig: &config.StoredConfig{
					AWSManifestBucket: "manifestbucket",
				},
			})
			aws.EXPECT().Client().AnyTimes().Return(awsClient)
			if tc.expectAWS != nil {
				tc.expectAWS(awsClient)
			}

			f := bytes.NewReader([]byte(tc.data))
			s := NewService(mm, conf, aws)
			ms := s.Manifest.(*manifestStore)
			err := ms.init(f, "/testassets")
			if tc.expectError == "" {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
				require.Equal(t, tc.expectError, err.Error())
			}
		})
	}
}
