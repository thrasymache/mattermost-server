package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWatcherInvalidDirectory(t *testing.T) {
	tempDir := os.TempDir()
	f, err := ioutil.TempFile(tempDir, "TestWatcher")
	require.NoError(t, err)

	callback := func() {}
	_, err = newWatcher(f.Name(), callback)
	require.Error(t, err, "should have failed to watch the entire temp directory")
}

func TestWatcher(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestWatcher")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	f, err := ioutil.TempFile(tempDir, "TestWatcher")
	require.NoError(t, err)
	defer f.Close()
	defer os.Remove(f.Name())

	called := make(chan bool)
	callback := func() {
		called <- true
	}
	watcher, err := newWatcher(f.Name(), callback)
	require.NoError(t, err)
	defer watcher.Close()

	// Write to a different file
	ioutil.WriteFile(filepath.Join(tempDir, "unrelated"), []byte("data"), 0644)
	select {
	case <-called:
		t.Fatal("callback should not have been called for unrelated file")
	case <-time.After(1 * time.Second):
	}

	// Write to the watched file
	ioutil.WriteFile(f.Name(), []byte("data"), 0644)
	select {
	case <-called:
	case <-time.After(5 * time.Second):
		t.Fatal("callback should have been called when file written")
	}

	// Delete the file
	err = os.Remove(f.Name())
	require.NoError(t, err)
	select {
	case <-called:
		t.Fatal("callback should not have been called for deletion")
	case <-time.After(1 * time.Second):
	}

	// Create the file
	ioutil.WriteFile(f.Name(), []byte("data"), 0644)
	select {
	case <-called:
	case <-time.After(5 * time.Second):
		t.Fatal("callback should have been called when file created")
	}
}