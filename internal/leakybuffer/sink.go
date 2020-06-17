package leakybuffer

import "github.com/buttairfly/goPanel/internal/hardware"

// DumpFrame produces a frame channel with a recycled frame or new one
func DumpFrame(f hardware.Frame) {
	// Reuse buffer if there's room.
	select {
	case freeList <- f:
		// Buffer on free list; nothing more to do.
	default:
		// Free list full, just carry on.
	}
}
