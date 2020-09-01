package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"code.htres.cn/casicloud/alb/pkg/model"
)

// LBPolicyStore help store lbpolicy
type LBPolicyStore interface {
	Open() error
	Close() error

	Get() []model.LBPolicy
	GetByID(id string) (model.LBPolicy, bool)
	Save() error
	Load() error
	Add(model.LBPolicy) error
	Delete(model.LBPolicy) error
	DeleteAll() error
}

const policyPath = "policy.store"

type fileStore struct {
	path  string
	mux   sync.Mutex
	cache map[string]model.LBPolicy
}

// NewLBPolicyStore create LBPolicyStore
func NewLBPolicyStore(config *Config) (LBPolicyStore, error) {
	path := filepath.Join(config.WorkDir, policyPath)
	cache := make(map[string]model.LBPolicy, 0)
	return &fileStore{
		path:  path,
		cache: cache,
	}, nil
}

// Get get policy
func (f *fileStore) Get() []model.LBPolicy {
	arr := []model.LBPolicy{}
	for _, v := range f.cache {
		arr = append(arr, v)
	}

	return arr
}

// Get get policy
func (f *fileStore) GetByID(id string) (model.LBPolicy, bool) {
	val, ok := f.cache[id]
	return val, ok
}

// Open for use
func (f *fileStore) Open() error {
	return nil
}

// Save save all policy
func (f *fileStore) Save() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.saveCache()
}

//Load load all policy from files
func (f *fileStore) Load() error {
	f.mux.Lock()
	defer f.mux.Unlock()

	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		// if not exists, return empty slice
		return nil
	}

	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}
	cache := make(map[string]model.LBPolicy, 0)
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return err
	}
	f.cache = cache
	return nil
}

//Add add new policy
func (f *fileStore) Add(n model.LBPolicy) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	f.cache[n.GetID()] = n
	return f.saveCache()
}

//Delete policy by id
func (f *fileStore) Delete(n model.LBPolicy) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	delete(f.cache, n.GetID())
	return f.saveCache()
}

//Close filestore
func (f *fileStore) Close() error {
	return nil
}

//DeleteAll remove all cache files
func (f *fileStore) DeleteAll() error {
	f.mux.Lock()
	defer f.mux.Unlock()

	for k := range f.cache {
		delete(f.cache, k)
	}

	var err = os.Remove(f.path)
	if err != nil {

		if os.IsNotExist(err) {
			sysLogger.Debugf("delete policy file not exists")
			return nil
		}
		return err
	}

	return nil
}

func (f *fileStore) saveCache() error {
	data, err := json.Marshal(f.cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.path, data, 0644)
	if err != nil {
		return err
	}
	sysLogger.Debugf("save policy\n")
	return nil
}
