package sumfs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func NewFS() *FS {
	return &FS{
		Sums: map[string]string{},
	}
}

type FS struct {
	Sums map[string]string
	lock sync.Mutex
}

func (f *FS) AddSum(key, value string) {
	f.lock.Lock()
	f.Sums[key] = value
	f.lock.Unlock()
}

func (f *FS) Read(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		wg := sync.WaitGroup{}
		errs := make(chan error)
		for _, entry := range entries {
			wg.Add(1)
			subpath := filepath.Join(path, entry.Name())
			go func() {
				err := f.Read(subpath)
				if err != nil {
					errs <- err
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(errs)
		var er []error
		for e := range errs {
			er = append(er, e)
		}
		if len(er) > 0 {
			s := ""
			for i, e := range er {
				if i > 0 {
					s += ", "
				}
				s += e.Error()
			}
			return errors.New(s)
		}
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return err
	}
	sum := hex.EncodeToString(h.Sum(nil))
	f.AddSum(path, sum)
	return nil
}
