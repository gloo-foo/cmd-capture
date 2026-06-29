package capture_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gloo-foo/testable"

	capture "github.com/gloo-foo/cmd-capture"
)

func ExampleCapture_bufferProcessing() {
	// Capture output, then post-process the captured buffer.
	var buf bytes.Buffer
	_, _ = testable.TestLines(capture.Capture(&buf), "alpha\nbravo\n")
	fmt.Printf("captured %d lines\n", strings.Count(buf.String(), "\n"))
	// Output:
	// captured 2 lines
}
