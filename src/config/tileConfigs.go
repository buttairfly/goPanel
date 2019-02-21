package config

import (
	"log"
	"sort"
)

type tileConfigs []TileConfig

// TileConfigs is a slice of TileConfig
type TileConfigs interface {
	sort.Interface
	GetSlice() []TileConfig
	Set(index int, tileConfig *tileConfig)
}

func (tc tileConfigs) Len() int {
	return len(tc)
}

func (tc tileConfigs) Less(i, j int) bool {
	if tc[i].GetConnectionOrder() == tc[j].GetConnectionOrder() {
		log.Printf("ConnectionOrder of two modules (%d,%d) must not be equal: %d, %d",
			i, j, tc[i].GetConnectionOrder(), tc[j].GetConnectionOrder())
	}
	return tc[i].GetConnectionOrder() < tc[j].GetConnectionOrder()
}

func (tc tileConfigs) Swap(i, j int) {
	tc[i], tc[j] = tc[j], tc[i]
}

func (tc tileConfigs) GetSlice() []TileConfig {
	return []TileConfig(tc)
}

func (tc tileConfigs) Set(index int, tileConfig *tileConfig) {
	tc.GetSlice()[index] = TileConfig(tileConfig)
}
