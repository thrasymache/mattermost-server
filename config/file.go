package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils/fileutils"
)

type authorizationError string

func (e authorizationError) IsAuthorizationError() bool { return true }
func (e authorizationError) Error() string              { return string(e) }

var (
	ReadOnlyConfigurationError = errors.New("configuration is read-only")
)

// TODO: config watch, enable, disable, callbacks
type fileStore struct {
	emitter

	config               *model.Config
	environmentOverrides map[string]interface{}
	path                 string
	watch                bool
	watcher              *watcher
	needsSave            bool
}

// NewFileStore creates a new instance of a config store backed by the given file path.
func NewFileStore(path string, watch bool) (*fileStore, error) {
	resolvedPath, err := resolveConfigFilePath(path)
	if err != nil {
		return nil, err
	}

	fs := &fileStore{
		path:  resolvedPath,
		watch: watch,
	}
	if err = fs.Load(); err != nil {
		return nil, err
	}

	if fs.watch {
		fs.startWatcher()
	}

	return fs, nil
}

// resolveConfigFilePath attempts to resolve an absolute path to the given configuration file path.
//
// Considerations include backwards compatibility when searching for configuration files from the
// myriad of supported input styles in various releases to date.
func resolveConfigFilePath(path string) (string, error) {
	// Absolute paths are explicit and require no resolution.
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Search for the given relative path or filename in various directories, resolving to the
	// corresponding absolute path if found.
	if configFile := fileutils.FindConfigFile(path); configFile != "" {
		return configFile, nil
	}

	// Otherwise, search for the config/ folder and build an absolute path anchored there and
	// joining the given input path.
	if configFolder, found := fileutils.FindDir("config"); found {
		return filepath.Join(configFolder, path), nil
	}

	// Fail altogether if we can't even find the config/ folder. This should only happen if
	// the executable is relocated away from the supporting files.
	return "", fmt.Errorf("failed to resolve config file path from %s", path)
}

// Get fetches the current configuration.
func (fs *fileStore) Get() *model.Config {
	return fs.config
}

// GetEnvironmentOverrides fetches the configuration fields overridden by environment variables.
func (fs *fileStore) GetEnvironmentOverrides() map[string]interface{} {
	return fs.environmentOverrides
}

// Set replaces the current configuration in its entirety.
func (fs *fileStore) Set(newCfg *model.Config) (*model.Config, error) {
	oldCfg := fs.Get()

	newCfg.SetDefaults()

	// Sometimes the config is received with "fake" data in sensitive fields. Apply the real
	// data from the existing config as necessary.
	desanitize(oldCfg, newCfg)

	if err := newCfg.IsValid(); err != nil {
		return nil, errors.Wrap(err, "new configuration is invalid")
	}

	if *oldCfg.ClusterSettings.Enable && *oldCfg.ClusterSettings.ReadOnlyConfig {
		return nil, authorizationError("configuration is read-only")
	}

	if err := fs.persist(newCfg); err != nil {
		return nil, errors.Wrap(err, "failed to persist")
	}

	go func() {
		fs.invokeConfigListeners(oldCfg, newCfg)
	}()

	return oldCfg, nil
}

// Patch merges the given config with the current configuration.
func (fs *fileStore) Patch(*model.Config) (*model.Config, error) {
	// TODO
	return fs.config, nil
}

// serialize converts the given configuration into JSON bytes for persistence.
func (fs *fileStore) serialize(cfg *model.Config) ([]byte, error) {
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal")
	}

	return b, nil
}

// persist writes the configuration to the configured file.
func (fs *fileStore) persist(cfg *model.Config) error {
	fs.stopWatcher()

	b, err := fs.serialize(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to serialize")
	}

	err = ioutil.WriteFile(fs.path, b, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	if fs.watch {
		fs.startWatcher()
	}

	return nil
}

// Load updates the current configuration from the backing store.
func (fs *fileStore) Load() error {
	var f io.ReadCloser
	var err error

	f, err = os.Open(fs.path)
	if os.IsNotExist(err) {
		fs.needsSave = true
		defaultCfg := model.Config{}
		defaultCfg.SetDefaults()

		defaultCfgBytes, err := fs.serialize(&defaultCfg)
		if err != nil {
			return errors.Wrap(err, "failed to serialize default config")
		}

		f = ioutil.NopCloser(bytes.NewReader(defaultCfgBytes))

	} else if err != nil {
		return errors.Wrapf(err, "failed to open %s for reading", fs.path)
	}
	defer f.Close()

	allowEnvironmentOverrides := true
	loadedCfg, environmentOverrides, err := readConfig(f, allowEnvironmentOverrides)
	if err != nil {
		return errors.Wrapf(err, "failed to load config from %s", fs.path)
	}

	loadedCfg.SetDefaults()

	// TODO: Move this out?
	*loadedCfg.ServiceSettings.SiteURL = strings.TrimRight(*loadedCfg.ServiceSettings.SiteURL, "/")

	oldCfg := fs.config
	fs.config = loadedCfg
	fs.environmentOverrides = environmentOverrides
	fs.invokeConfigListeners(oldCfg, loadedCfg)

	return nil
}

// startWatcher starts a watcher to monitor for external config file changes.
func (fs *fileStore) startWatcher() error {
	if fs.watcher != nil {
		return nil
	}

	watcher, err := newWatcher(fs.path, func() {
		fs.Load()
	})
	if err != nil {
		return err
	}

	fs.watcher = watcher

	return nil
}

// stopWatcher stops any previously started watcher.
func (fs *fileStore) stopWatcher() {
	if fs.watcher == nil {
		return
	}

	fs.watcher.Close()
	fs.watcher = nil
}

// String returns the path to the file backing the config.
func (fs *fileStore) String() string {
	return "file://" + fs.path
}

// Close cleans up resources associated with the store.
func (fs *fileStore) Close() error {
	fs.stopWatcher()

	return nil
}
