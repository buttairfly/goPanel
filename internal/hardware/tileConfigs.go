package hardware

// TileConfigs is a slice of TileConfig
type TileConfigs [](*TileConfig)

// Set implmepents TileConfigs function
func (tcs TileConfigs) Set(index int, tileConfig *TileConfig) {
	tcs[index] = tileConfig
}

// MarshalTileConfigs is a slice of MarsshalTileConfig
type MarshalTileConfigs [](*MarshalTileConfig)

func (mtcs MarshalTileConfigs) Len() int {
	return len(mtcs)
}

func (mtcs MarshalTileConfigs) Less(i, j int) bool {
	return mtcs[i].ToTileConfig().GetConnectionOrder() <= mtcs[j].ToTileConfig().GetConnectionOrder()
}

func (mtcs MarshalTileConfigs) Swap(i, j int) {
	mtcs[i], mtcs[j] = mtcs[j], mtcs[i]
}

// ToTileConfigs retruns the TileConfig
func (mtcs MarshalTileConfigs) ToTileConfigs() TileConfigs {
	tcs := make(TileConfigs, len(mtcs))
	for i, mtc := range mtcs {
		tcs[i] = mtc.ToTileConfig()
	}
	return tcs
}

// ToMarshalTileConfigs retruns the marshalable tile config
func (tcs TileConfigs) ToMarshalTileConfigs() MarshalTileConfigs {
	mtcs := make(MarshalTileConfigs, len(tcs))
	for i, tc := range tcs {
		mtcs[i] = tc.ToMarshalTileConfig()
	}
	return mtcs
}
