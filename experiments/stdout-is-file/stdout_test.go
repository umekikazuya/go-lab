package stdout_test

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"syscall"
	"testing"

	stdout "go-lab/experiments/stdout-is-file"
)

// ============================================================
// 前提検証（testing.T）: 構造的等価性
// ============================================================

// TestTypeIdentity verifies that os.Stdout is the exact same Go type
// as a file returned by os.Create — both are *os.File.
func TestTypeIdentity(t *testing.T) {
	// Compile-time verification: os.Stdout is assignable to *os.File.
	var _ *os.File = os.Stdout

	f, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stdoutType := reflect.TypeOf(os.Stdout)
	fileType := reflect.TypeOf(f)
	if stdoutType != fileType {
		t.Fatalf("type mismatch: Stdout=%v, file=%v", stdoutType, fileType)
	}
	t.Logf("os.Stdout type = %v (same as os.Create result)", stdoutType)
}

// TestFileDescriptor verifies that os.Stdout wraps POSIX fd 1,
// and regular files receive fd >= 3 (after stdin=0, stdout=1, stderr=2).
func TestFileDescriptor(t *testing.T) {
	if fd := os.Stdout.Fd(); fd != 1 {
		t.Fatalf("os.Stdout.Fd() = %d, want 1", fd)
	}
	t.Logf("os.Stdout.Fd() = %d", os.Stdout.Fd())

	f, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if fd := f.Fd(); fd < 3 {
		t.Fatalf("regular file fd = %d, want >= 3", fd)
	}
	t.Logf("regular file fd = %d", f.Fd())
}

// TestWriterInterface verifies that both os.Stdout and a regular file
// satisfy io.Writer, and WriteViaInterface works identically with either.
func TestWriterInterface(t *testing.T) {
	// Compile-time verification.
	var _ io.Writer = os.Stdout

	f, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var _ io.Writer = f

	data := []byte("interface test\n")

	// Write to /dev/null via interface (avoids polluting test output).
	devNull, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()

	n1, err := stdout.WriteViaInterface(devNull, data)
	if err != nil {
		t.Fatalf("WriteViaInterface(devNull): %v", err)
	}
	n2, err := stdout.WriteViaInterface(f, data)
	if err != nil {
		t.Fatalf("WriteViaInterface(file): %v", err)
	}
	if n1 != n2 {
		t.Fatalf("byte count mismatch: devNull=%d, file=%d", n1, n2)
	}
	t.Logf("WriteViaInterface: both wrote %d bytes", n1)
}

// TestRedirection proves that os.Stdout is a mutable variable.
// Replacing it causes fmt.Println to write to the new target,
// demonstrating that fmt.Println has no special "terminal" logic.
func TestRedirection(t *testing.T) {
	orig := os.Stdout
	defer func() { os.Stdout = orig }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w
	fmt.Println("redirected")
	w.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	got := string(buf)
	want := "redirected\n"
	if got != want {
		t.Fatalf("after redirection: got %q, want %q", got, want)
	}
	t.Logf("fmt.Println wrote to pipe: %q", got)
}

// TestRedirectionToFile proves that os.Stdout can be replaced with
// a regular disk file. fmt.Println writes to that file, and the
// content can be read back — the most direct proof that stdout
// is "just a writable file".
func TestRedirectionToFile(t *testing.T) {
	orig := os.Stdout
	defer func() { os.Stdout = orig }()

	path := t.TempDir() + "/stdout.txt"
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = f
	fmt.Println("written to file")
	f.Close()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(content)
	want := "written to file\n"
	if got != want {
		t.Fatalf("file content: got %q, want %q", got, want)
	}
	t.Logf("fmt.Println wrote to disk file: %q", got)
}

// ============================================================
// 主実験（testing.B）: 書き込みパスの性能等価性
//
// 書き込み先を /dev/null に統一し、カーネル側の差異を排除する。
// ============================================================

var benchData = []byte("hello, world\n") // 13 bytes, fixed

// sink prevents the compiler from eliminating Write calls.
var sink int

// BenchmarkWriteStdout writes through os.Stdout (redirected to /dev/null).
// Pattern A: os.File.Write via the os.Stdout variable.
func BenchmarkWriteStdout(b *testing.B) {
	devNull, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer devNull.Close()

	orig := os.Stdout
	os.Stdout = devNull
	defer func() { os.Stdout = orig }()

	b.ResetTimer()
	for b.Loop() {
		n, _ := os.Stdout.Write(benchData)
		sink = n
	}
}

// BenchmarkWriteFile writes through a regular *os.File (→ /dev/null).
// Pattern B: os.File.Write via a local variable.
func BenchmarkWriteFile(b *testing.B) {
	devNull, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer devNull.Close()

	b.ResetTimer()
	for b.Loop() {
		n, _ := devNull.Write(benchData)
		sink = n
	}
}

// BenchmarkSyscallWriteFd1 writes using syscall.Write with fd 1
// (kernel-level stdout, redirected to /dev/null via dup2).
// Pattern C: raw syscall through the stdout file descriptor.
func BenchmarkSyscallWriteFd1(b *testing.B) {
	devNull, err := syscall.Open("/dev/null", syscall.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer syscall.Close(devNull)

	origFd, err := syscall.Dup(1)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		syscall.Dup2(origFd, 1)
		syscall.Close(origFd)
	}()

	if err := syscall.Dup2(devNull, 1); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		n, _ := syscall.Write(1, benchData)
		sink = n
	}
}

// BenchmarkSyscallWriteFileFd writes using syscall.Write with a
// regular file descriptor (→ /dev/null).
// Pattern D: raw syscall through a non-stdout file descriptor.
func BenchmarkSyscallWriteFileFd(b *testing.B) {
	fd, err := syscall.Open("/dev/null", syscall.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer syscall.Close(fd)

	b.ResetTimer()
	for b.Loop() {
		n, _ := syscall.Write(fd, benchData)
		sink = n
	}
}
