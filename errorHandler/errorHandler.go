package errorHandler

import (
	"fmt"
	"io"
)

func ReportError(w io.Writer, msg string, where string, line int) {
	fmt.Fprintf(w, "[line: %v] Error: %s %s\n", line, msg, where)

}
