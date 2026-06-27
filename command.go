package command

import (
	"io"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// Capture returns a Command that passes every line through unchanged while
// writing it (newline-terminated) to each provided writer. Place it anywhere in
// a pipeline to capture the stream into one or more buffers without altering it;
// with no writers it is a pass-through.
//
// Capturing is a side-effect, so Capture is a Tap: the line continues
// downstream after being written.
func Capture(writers ...io.Writer) gloo.Command[[]byte, []byte] {
	return patterns.Tap(func(line []byte) error {
		for _, w := range writers {
			if err := writeLine(w, line); err != nil {
				return err
			}
		}
		return nil
	})
}

// writeLine writes line to w followed by a newline, reporting the first error.
func writeLine(w io.Writer, line []byte) error {
	if _, err := w.Write(line); err != nil {
		return err
	}
	_, err := w.Write(newline)
	return err
}

var newline = []byte{'\n'}
