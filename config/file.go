package config

import (
	"encoding/json"
	"fmt"
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
	config               *model.Config
	environmentOverrides map[string]interface{}
	path                 string
	disableWatcher       bool
	watcher              *ConfigWatcher
}

func NewFileStore(path string) (*fileStore, error) {
	var err error
	path, err = resolveConfigFilePath(path)
	if err != nil {
		return nil, err
	}

	fs := &fileStore{
		path:           path,
		disableWatcher: false,
	}
	if err = fs.Load(); err != nil {
		return nil, err
	}

	return fs, nil
}

// resolveConfigFilePath attempts to resolve the given path to a configuration file to an absolute path.
//
// Considerations include backwards compatibility when searching for configuration files from the
// myriad of supported inputs in various versions to date.
func resolveConfigFilePath(path string) (string, error) {
	// Absolute paths are explicit and require no resolution.
	if !filepath.IsAbs(path) {
		return path, nil
	}

	// Search for the given relative path or filename in various directories, resolving to the
	// absolute path if found.
	if configFile := fileutils.FindConfigFile(path); configFile != "" {
		return configFile, nil
	}

	// Otherwise, search for the config/ folder and build an absolute path anchored there.
	if configFolder, found := fileutils.FindDir("config"); found {
		return filepath.Join(configFolder, path), nil
	}

	// Fail altogether if we can't even find the config/ folder. This should only happen if
	// the executable is relocated away from the supporting files.
	return "", fmt.Errorf("failed to resolve config file path from %s", path)
}

func (s *fileStore) Get() *model.Config {
	return s.config
}

func (s *fileStore) startWatcher() error {
	if s.watcher != nil {
		return nil
	}

	/*
		if s.disableWatcher {
			return nil
		}
	*/

	watcher, err := NewConfigWatcher(s.path, func() {
		// s.ReloadConfig()
	})
	if err != nil {
		return err
	}

	s.watcher = watcher
	return nil
}

func (s *fileStore) stopWatcher() {
	if s.watcher != nil {
		s.watcher.Close()
		s.watcher = nil
	}
}

// TODO: EnableWatcher?
func (s *fileStore) DisableWatcher() {
	s.disableWatcher = false
	s.stopWatcher()
}

func (s *fileStore) Load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s for reading", s.path)
	}
	defer f.Close()

	allowEnvironmentOverrides := true
	loadedCfg, environmentOverrides, err := readConfig(f, allowEnvironmentOverrides)
	if err != nil {
		return errors.Wrapf(err, "failed to load config from %s", s.path)
	}

	// TODO: Move this out?
	*loadedCfg.ServiceSettings.SiteURL = strings.TrimRight(*loadedCfg.ServiceSettings.SiteURL, "/")

	oldCfg := s.config
	s.config = loadedCfg
	s.environmentOverrides = environmentOverrides
	s.invokeConfigListeners(oldCfg, loadedCfg)

	return nil
}

func (s *fileStore) invokeConfigListeners(oldCfg, newCfg *model.Config) error {
	return nil
}

func (s *fileStore) Set(newCfg *model.Config) error {
	oldCfg := s.Get()

	newCfg.SetDefaults()

	// Sometimes the config is received with "fake" data in sensitive fields. Apply the real
	// data from the existing config as necessary.
	desanitize(oldCfg, newCfg)

	if err := newCfg.IsValid(); err != nil {
		return errors.Wrap(err, "new configuration is invalid")
	}

	if *oldCfg.ClusterSettings.Enable && *oldCfg.ClusterSettings.ReadOnlyConfig {
		return authorizationError("configuration is read-only")
	}

	if err := s.persist(newCfg); err != nil {
		return errors.Wrap(err, "failed to persist")
	}

	return nil
}

func (s *fileStore) persist(newCfg *model.Config) error {
	s.stopWatcher()

	b, err := json.MarshalIndent(newCfg, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	err = ioutil.WriteFile(s.path, b, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	if !s.disableWatcher {
		s.startWatcher()
	}

	return nil
}

func (s *fileStore) Close() error {
	return nil
}
