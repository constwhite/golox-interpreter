package errorHandler

import (
	"fmt"
	"io"
)

func ReportError(w io.Writer, msg string, line int) {
	fmt.Fprintf(w, "[line: %v] Error: %s\n", line, msg)
}
