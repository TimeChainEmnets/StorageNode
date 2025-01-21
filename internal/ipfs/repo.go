package ipfs

import "os"

type RepoManager struct {
	path string
}

func NewRepoManager(path string) *RepoManager {

	return &RepoManager{

		path: path,
	}

}

func (r *RepoManager) CheckRepo() (bool, error) {

	_, err := os.Stat(r.path)

	if os.IsNotExist(err) {

		return false, nil

	}

	if err != nil {

		return false, err

	}

	return true, nil

}

func (r *RepoManager) InitRepo() error {

	return os.MkdirAll(r.path, 0755)

}
