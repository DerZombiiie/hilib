// taken from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
package hilib

import (
	"fmt"
)

// BytesCountSI converts a bytevalue to kB, MB, GB, TB, PB, EB (1000 steps)
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// BytesCountSI converts a bytevalue to kiB, MiB, GiB, TiB, PiB, EiB (1024 steps)
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
