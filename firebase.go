// package remoteconfig provides a thin wrapper around the Firebase Remote
// Config admin API to support automated updates of config values.
package remoteconfig

import (
	"context"

	rc "google.golang.org/api/firebaseremoteconfig/v1"
	"google.golang.org/api/option"
)

const scopeFirebaseRemoteConfig = "https://www.googleapis.com/auth/firebase.remoteconfig"

// Client wraps a Google Cloud projectID and a remoteconfig.ProjectsService
type Client struct {
	projectID string
	ps        *rc.ProjectsService
}

// NewClient initializes a new Client with the provided projectID and
// credentials (Google Cloud Service Account key, JSON). Panics on error.
func NewClient(ctx context.Context, projectID string, credentials []byte) *Client {
	svc, err := rc.NewService(
		ctx,
		option.WithCredentialsJSON(credentials),
		option.WithScopes(scopeFirebaseRemoteConfig),
	)
	if err != nil {
		panic(err)
	}

	return &Client{projectID, svc.Projects}
}

// SetDefaultValues updates the published remote config's default key/value
// pairs with the provided params. Parameters that already exist in the
// published remote config will be overwritten.
func (c *Client) SetDefaultValues(ctx context.Context, params map[string]string) error {
	config, err := c.ps.GetRemoteConfig("projects/" + c.projectID).Context(ctx).Do()
	if err != nil {
		return err
	}

	update := &rc.RemoteConfig{
		Conditions: config.Conditions,
		Parameters: config.Parameters,
	}

	for k, v := range params {
		update.Parameters[k] = rc.RemoteConfigParameter{
			DefaultValue: &rc.RemoteConfigParameterValue{
				Value: v,
			},
		}
	}

	call := c.ps.UpdateRemoteConfig("projects/"+c.projectID, update).Context(ctx)
	call.Header().Set("If-Match", config.Header.Get("ETag"))
	if _, err := call.Do(); err != nil {
		return err
	}

	return nil
}
