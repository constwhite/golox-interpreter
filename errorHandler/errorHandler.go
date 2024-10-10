package errorHandler

import (
	"fmt"
	"io"
)

func ReportError(w io.Writer, msg string, where string, line int) {
	fmt.Fprintf(w, "[line: %v] Error: %s %s\n", line, msg, where)

}
func RuntimeError(w io.Writer, err error, line int) {
	fmt.Fprintf(w, "%v\n[line:%v]", err, line)
}

func CompileError(w io.Writer, err error, line int) {
	fmt.Fprintf(w, "%v\n[line:%v]", err, line)
}
