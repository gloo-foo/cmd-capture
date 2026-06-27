package capture_test

import (
	"bytes"
	"fmt"

	capture "github.com/gloo-foo/cmd-capture"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"
)

func ExampleCapture_pipeline() {
	// Capture mid-pipeline: the buffer sees the stream, which still flows on.
	var mid bytes.Buffer
	pipe := gloo.Pipe(capture.Capture(&mid), capture.Capture())
	lines, _ := testable.TestLines(pipe, "x\ny\n")
	fmt.Printf("passed through %d lines; captured %q", len(lines), mid.String())
	// Output:
	// passed through 2 lines; captured "x\ny\n"
}
