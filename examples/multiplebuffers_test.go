package capture_test

import (
	"bytes"
	"fmt"

	"github.com/gloo-foo/testable"

	capture "github.com/gloo-foo/cmd-capture"
)

func ExampleCapture_multipleBuffers() {
	// One call can fan the stream out to several writers.
	var a, b bytes.Buffer
	_, _ = testable.TestLines(capture.Capture(&a, &b), "data\n")
	fmt.Printf("a=%q b=%q", a.String(), b.String())
	// Output:
	// a="data\n" b="data\n"
}
