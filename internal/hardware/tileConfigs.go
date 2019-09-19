package hardware

import (
	"log"
	"sort"
)

// TileConfigSlice implements TileConfigs
type TileConfigSlice []TileConfig

// TileConfigs is a slice of TileConfig
type TileConfigs interface {
	sort.Interface
	GetSlice() []TileConfig
	Set(index int, tileConfig TileConfig)
}

func (tc TileConfigSlice) Len() int {
	return len(tc)
}

func (tc TileConfigSlice) Less(i, j int) bool {
	if tc[i].GetConnectionOrder() == tc[j].GetConnectionOrder() {
		log.Fatalf("ConnectionOrder of two modules (%d,%d) must not be equal: %d, %d",
			i, j, tc[i].GetConnectionOrder(), tc[j].GetConnectionOrder())
	}
	return tc[i].GetConnectionOrder() < tc[j].GetConnectionOrder()
}

func (tc TileConfigSlice) Swap(i, j int) {
	tc[i], tc[j] = tc[j], tc[i]
}

// GetSlice implmepents TileConfigs function
func (tc TileConfigSlice) GetSlice() []TileConfig {
	return []TileConfig(tc)
}

// Set implmepents TileConfigs function
func (tc TileConfigSlice) Set(index int, tileConfig TileConfig) {
	tc.GetSlice()[index] = tileConfig
}
