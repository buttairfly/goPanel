package hardware

// TileConfigs is a slice of TileConfig
type TileConfigs [](*TileConfig)

func (tc TileConfigs) Len() int {
	return len(tc)
}

func (tc TileConfigs) Less(i, j int) bool {
	return tc[i].GetConnectionOrder() <= tc[j].GetConnectionOrder()
}

func (tc TileConfigs) Swap(i, j int) {
	tc[i], tc[j] = tc[j], tc[i]
}

// Set implmepents TileConfigs function
func (tc TileConfigs) Set(index int, tileConfig *TileConfig) {
	tc[index] = tileConfig
}
