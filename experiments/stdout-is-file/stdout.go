package stdout

import (
	"io"
	"os"
	"syscall"
)

// WriteViaFile writes data using the os.File.Write method.
// Both os.Stdout and a regular file use this identical code path.
func WriteViaFile(f *os.File, data []byte) (int, error) {
	return f.Write(data)
}

// WriteViaSyscall writes data using syscall.Write, bypassing
// Go's os.File wrapper (internal mutex, poll descriptor).
// This is the lowest-level write available from user space.
func WriteViaSyscall(fd int, data []byte) (int, error) {
	return syscall.Write(fd, data)
}

// WriteViaInterface writes data through the io.Writer interface.
// os.Stdout and regular files both satisfy io.Writer, making them
// interchangeable at the interface level.
func WriteViaInterface(w io.Writer, data []byte) (int, error) {
	return w.Write(data)
}
