package exit

import (
	"fmt"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"

	"go.uber.org/zap"
)

// Goroutine is a struct to show goroutines
type Goroutine struct {
	Number int64
	Action string
	Path   string
	Origin string
	Line   int64
}

// GetGoroutine returns a slice of all currently running goroutines
func GetGoroutine(logger *zap.Logger) []Goroutine {
	debug.SetTraceback("all")
	buf := make([]byte, 1<<16)
	bufLen := runtime.Stack(buf, true)
	debug.SetTraceback("default")
	//logger.Info("stack", zap.ByteString("stack", buf[:bufLen]))

	stackRegex := regexp.MustCompile(`goroutine (\d+) \[([\w \-]+)\]:\n([\w\-\_/\\,.:+()*$&§"=! \n\t]+)\n`)
	pathRegex := regexp.MustCompile(`([\w\-\_/\\,.:+()*$&§"=!]+)\n\t([\w./\\]+):(\d+) \+0x[\dA-Fa-f]+\n?$`)
	find := stackRegex.FindAllSubmatch(buf[:bufLen], -1)
	// fmt.Printf("find %q\n", find)
	if find != nil {
		gophers := make([]Goroutine, len(find))
		for i, current := range find {
			num, err := strconv.ParseInt(string(current[1]), 10, 64)
			if err != nil {
				logger.Panic("convert goroutine number to int", zap.Error(err))
			}
			routine := Goroutine{
				Number: num,
				Action: string(current[2]),
			}
			pathFind := pathRegex.FindSubmatch(current[3])
			//fmt.Printf("pathFind %d %q\n", num, pathFind)
			if pathFind != nil {
				routine.Origin = string(pathFind[1])
				routine.Path = string(pathFind[2])

				num, err = strconv.ParseInt(string(pathFind[3]), 10, 64)
				if err != nil {
					logger.Panic("convert path line to int", zap.Error(err))
				}
				routine.Line = num
			} else {
				routine.Path = string(current[3])
			}
			gophers[i] = routine
		}
		// fmt.Printf("gophers %+v\n", gophers)
		return gophers
	}
	return []Goroutine{}
}

// ShortGoroutine returns a short printable
func (g *Goroutine) ShortGoroutine() string {
	return fmt.Sprintf("%d [%s]", g.Number, g.Action)
}
