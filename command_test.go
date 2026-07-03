package command_test

import (
	"bytes"
	"errors"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-capture"
)

// errWrite is the sentinel a failWriter returns once its allowance is spent.
var errWrite = errors.New("write failed")

// failWriter succeeds for okWrites calls, then fails every subsequent Write
// with errWrite. It lets a test pinpoint which io.Writer.Write inside Capture
// aborts the pipeline: okWrites=0 fails the line write, okWrites=1 fails the
// trailing-newline write.
type failWriter struct {
	okWrites int
	calls    int
}

func (w *failWriter) Write(p []byte) (int, error) {
	if w.calls >= w.okWrites {
		w.calls++
		return 0, errWrite
	}
	w.calls++
	return len(p), nil
}

func TestCapture_Passthrough(t *testing.T) {
	var buf bytes.Buffer
	lines, err := testable.TestLines(command.Capture(&buf), "alpha\nbravo\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0] != "alpha" || lines[1] != "bravo" {
		t.Fatalf("passthrough: got %q", lines)
	}
	if buf.String() != "alpha\nbravo\n" {
		t.Fatalf("captured: got %q, want %q", buf.String(), "alpha\nbravo\n")
	}
}

func TestCapture_NoWriters(t *testing.T) {
	lines, err := testable.TestLines(command.Capture(), "x\ny\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %q", lines)
	}
}

func TestCapture_MultipleWriters(t *testing.T) {
	var a, b bytes.Buffer
	if _, err := testable.TestLines(command.Capture(&a, &b), "one\ntwo\n"); err != nil {
		t.Fatal(err)
	}
	if a.String() != "one\ntwo\n" || b.String() != "one\ntwo\n" {
		t.Fatalf("a=%q b=%q", a.String(), b.String())
	}
}

func TestCapture_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	lines, err := testable.TestLines(command.Capture(&buf), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 || buf.Len() != 0 {
		t.Fatalf("got lines=%q buf=%q", lines, buf.String())
	}
}

// TestCapture_LineWriteErrorPropagates: a writer that fails on the line write
// aborts the pipeline and surfaces the exact error.
func TestCapture_LineWriteErrorPropagates(t *testing.T) {
	w := &failWriter{okWrites: 0}
	_, err := testable.TestLines(command.Capture(w), "alpha\n")
	if !errors.Is(err, errWrite) {
		t.Fatalf("got err %v, want %v", err, errWrite)
	}
}

// TestCapture_NewlineWriteErrorPropagates: a writer that succeeds on the line
// write but fails on the trailing-newline write still aborts and propagates.
func TestCapture_NewlineWriteErrorPropagates(t *testing.T) {
	w := &failWriter{okWrites: 1}
	_, err := testable.TestLines(command.Capture(w), "alpha\n")
	if !errors.Is(err, errWrite) {
		t.Fatalf("got err %v, want %v", err, errWrite)
	}
}

// TestCapture_StopsAtFirstFailingWriter: when an earlier writer fails, the
// later writer is never written to (the loop aborts on the first error).
func TestCapture_StopsAtFirstFailingWriter(t *testing.T) {
	bad := &failWriter{okWrites: 0}
	var good bytes.Buffer
	_, err := testable.TestLines(command.Capture(bad, &good), "alpha\n")
	if !errors.Is(err, errWrite) {
		t.Fatalf("got err %v, want %v", err, errWrite)
	}
	if good.Len() != 0 {
		t.Fatalf("downstream writer received %q after upstream failure", good.String())
	}
}

// TestCapture_MidPipeline captures between two stages without disturbing flow.
func TestCapture_MidPipeline(t *testing.T) {
	var mid bytes.Buffer
	pipe := gloo.Pipe(command.Capture(&mid), command.Capture())
	lines, err := testable.TestLines(pipe, "go\nrust\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || mid.String() != "go\nrust\n" {
		t.Fatalf("lines=%q mid=%q", lines, mid.String())
	}
}
