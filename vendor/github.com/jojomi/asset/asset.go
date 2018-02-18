package asset

import (
	"fmt"
	"io/ioutil"
	"path"
)

type Handler struct {
	Sources []AssetSource
}

type AssetSource interface {
	Get(name string) ([]byte, error)
}

func (h *Handler) Get(name string) ([]byte, error) {
	if h.Sources == nil || len(h.Sources) == 0 {
		return nil, fmt.Errorf("no sources defined searching for %s", name)
	}

	for _, s := range h.Sources {
		data, err := s.Get(name)
		if err != nil {
			continue
		}

		return data, nil
	}
	return nil, fmt.Errorf("Asset not found: %s", name)
}

// helpers

type FSAssetSource struct {
	path string
}

func NewFSAssetSource(path string) *FSAssetSource {
	return &FSAssetSource{
		path: path,
	}
}

func (f *FSAssetSource) Get(name string) ([]byte, error) {
	data, err := ioutil.ReadFile(path.Join(f.path, name))
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BinDataAssetSource struct {
	assetFunc func(name string) ([]byte, error)
}

func NewBinDataAssetSource(assetFunc func(name string) ([]byte, error)) *BinDataAssetSource {
	return &BinDataAssetSource{
		assetFunc: assetFunc,
	}
}

func (f *BinDataAssetSource) Get(name string) ([]byte, error) {
	return f.assetFunc(name)
}
