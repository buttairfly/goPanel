package spots

import (
	"image"
	"io/ioutil"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/pkg/filereadwriter"
	"github.com/buttairfly/goPanel/pkg/testhelper"
)

var _ filereadwriter.Yaml = (*InputPictureConfig)(nil)

func TestNewSpotsFromConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	const testFolder = "testdata/"
	cases := []struct {
		desc           string
		config         InputPictureConfig
		resultFile     string
		actualFile     string
		configFileName string
		err            error
	}{
		{
			desc: "20x10_testconfig",
			config: InputPictureConfig{
				Offset:     image.Point{X: 10, Y: 10},
				TileWidth:  10,
				TileHeight: 10,
				TileSpots:  []image.Point{{X: 5, Y: 5}},
				Height:     10,
				Width:      20,
			},
			resultFile:     testFolder + "expected.spots",
			actualFile:     testFolder + "actual.spots",
			configFileName: testFolder + "20x10_spots.yaml",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {

			testhelper.FileExistsOrSkip(t, c.configFileName)

			if testhelper.RecordCall() {
				t.Logf("Write Config to file %v", c.configFileName)
				require.NoError(t, c.config.WriteToYamlFile(c.configFileName))
			}
			spots, err := NewSpotsFromConfig(c.configFileName, logger)
			data := []byte(spotsToStr(spots))
			if testhelper.RecordCall() {
				t.Logf("Write result to file %v", c.resultFile)
				require.NoError(t, ioutil.WriteFile(c.resultFile, data, 0622))
			}
			assert.Equal(t, c.err, err, "error occurred")
			assert.NoError(t, err, "error reading file")
			assert.NoError(t, ioutil.WriteFile(c.actualFile, data, 0622))
			defer os.Remove(c.actualFile)
			testhelper.Diff(t, c.resultFile, c.actualFile)
		})
	}
}

func spotsToStr(spots Spots) string {
	spew.Config.SortKeys = true
	spew.Config.DisablePointerAddresses = true
	return spew.Sdump(spots)
}
