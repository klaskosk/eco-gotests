package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"github.com/golang/glog"
)

const (
	cacheVersion = "1"
	cacheDir     = "eco-gotests"
	cacheFile    = "report-cache.json"
)

var (
	errCacheMiss        = fmt.Errorf("cache miss")
	buildInfo, useCache = debug.ReadBuildInfo()
)

// IsMiss returns true if the given error is a cache miss error and false otherwise.
func IsMiss(err error) bool {
	return errors.Is(err, errCacheMiss)
}

// Cache represents the format of the cache file. It will be saved as JSON according to the XDG base directory
// specification.
type Cache struct {
	Version string
	Trees   map[string]*SuiteTree
}

// NewCache creates a new cache instance. It will attempt to load the cache from the cache file. If the file does not
// exist, a new cache will be created but not saved until Save is called.
func NewCache() (*Cache, error) {
	cachePath := getCachePath()

	glog.V(100).Infof("Attempting to load cache from %s", cachePath)

	if _, err := os.Stat(cachePath); err == nil {
		return loadCache(cachePath)
	}

	return &Cache{
		Version: cacheVersion,
		Trees:   map[string]*SuiteTree{},
	}, nil
}

// Save saves the cache to the cache file.
func (cache *Cache) Save() error {
	glog.V(100).Infof("Saving cache with %d trees to %s", len(cache.Trees), getCachePath())

	cachePath := getCachePath()
	cacheDir := path.Dir(cachePath)

	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(cache)
	if err != nil {
		return err
	}

	return nil
}

// GetRemote checks to see if the repo with the given remote branch is in the cache and returns the suite tree without
// fetching it from the remote. It returns a cache miss error if the repo is not in the cache.
func (cache *Cache) GetRemote(repo, branch string) (*SuiteTree, error) {
	glog.V(100).Infof("Checking if repo %s with branch %s is in cache", repo, branch)

	key, err := GetRemoteRevision(repo, branch)
	if err != nil {
		glog.V(100).Infof("Failed to get remote revision for repo %s and branch %s: %v", repo, branch, err)

		return nil, err
	}

	tree, ok := cache.Trees[key]
	if !ok {
		glog.V(100).Infof("Repo %s with branch %s not in cache", repo, branch)

		return nil, errCacheMiss
	}

	return tree, nil
}

// Get returns the suite tree for the given repo path from the cache. It returns a cache miss error if the repo has
// uncommitted changes or if the cache does not contain the repo.
func (cache *Cache) Get(repoPath string) (*SuiteTree, error) {
	glog.V(100).Infof("Getting cache for repo %s", repoPath)

	key, err := getRepoKey(repoPath)
	if err != nil {
		return nil, err
	}

	if tree, ok := cache.Trees[key]; ok {
		return tree, nil
	}

	return nil, errCacheMiss
}

// GetOrCreate returns the suite tree for the given repo path from the cache. It first calls Get and if there is a cache
// miss, it calls the given create function and adds the result to the cache. Note that if the repo has local changes,
// the create function will always be called, but the result will not be added to the cache.
func (cache *Cache) GetOrCreate(repoPath string, create func() (*SuiteTree, error)) (*SuiteTree, error) {
	glog.V(100).Infof("Getting or creating cache for repo %s", repoPath)

	tree, err := cache.Get(repoPath)
	if err == nil {
		return tree, nil
	}

	if !IsMiss(err) {
		return nil, err
	}

	glog.V(100).Infof("Cache miss for repo %s, calling create function", repoPath)

	tree, err = create()
	if err != nil {
		return nil, err
	}

	key, err := getRepoKey(repoPath)
	if err == nil {
		cache.Trees[key] = tree
	} else {
		glog.V(100).Infof("Failed to get cache key for repo %s, not adding to cache", repoPath)
	}

	return tree, nil
}

// getRepoKey returns the key to use for the cache. It returns a cache miss error if the repo has uncommitted changes.
func getRepoKey(repoPath string) (string, error) {
	glog.V(100).Infof("Getting cache key for repo %s", repoPath)

	changes, err := HasLocalChanges(repoPath)
	if err != nil {
		return "", err
	}

	if changes {
		glog.V(100).Infof("Repo %s has uncommitted changes, cache will always miss", repoPath)

		return "", errCacheMiss
	}

	revision, err := GetRepoRevision(repoPath)
	if err != nil {
		return "", err
	}

	return revision, nil
}

// getCachePath returns the path to the cache file. It is not guaranteed to exist.
func getCachePath() string {
	glog.V(100).Infof(
		"Getting cache path in XDG_CACHE_HOME %s with cacheDir %s and cacheFile %s",
		os.Getenv("XDG_CACHE_HOME"), cacheDir, cacheFile)

	cacheHome, ok := os.LookupEnv("XDG_CACHE_HOME")
	if !ok {
		cacheHome = "."
	}

	return path.Join(cacheHome, cacheDir, cacheFile)
}

// loadCache loads the cache from the given path. It returns an error if the cache file could not be loaded or if the
// cache file has an incompatible version.
func loadCache(path string) (*Cache, error) {
	glog.V(100).Infof("Loading cache from %s", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	cache := &Cache{}
	err = json.NewDecoder(file).Decode(cache)

	if err != nil {
		return nil, err
	}

	if cache.Version != cacheVersion {
		return nil, fmt.Errorf("cache file %s has incompatible version %s; expected %s", path, cache.Version, cacheVersion)
	}

	return cache, nil
}
