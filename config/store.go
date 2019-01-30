package config

import (
	"github.com/mattermost/mattermost-server/model"
)

// Store abstracts the act of getting and setting the configuration.
type Store interface {
	// Get fetches the current configuration.
	Get() (*model.Config, error)

	// GetEnvironmentOverrides fetches the configuration fields overridden by environment variables.
	GetEnvironmentOverrides() map[string]interface{}

	// Set replaces the current configuration in its entirety.
	Set(*model.Config) error

	// Patch merges the given config with the current configuration.
	Patch(*model.Config) error

	// AddListener adds a callback function to invoke when the configuration is modified.
	AddListener(listener func(oldConfig *model.Config, newConfig *model.Config)) string

	// RemoveListener removes a callback function using an id returned from AddListener.
	RemoveListener(id string)

	// Close cleans up resources associated with the store.
	Close() error
}
