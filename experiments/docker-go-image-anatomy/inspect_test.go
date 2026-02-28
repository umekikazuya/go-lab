package dockergoimageanatomy

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestGoVersion は runtime.Version() が期待するプレフィクスを持つことを検証する。
func TestGoVersion(t *testing.T) {
	v := runtime.Version()
	if !strings.HasPrefix(v, "go") {
		t.Errorf("unexpected version format: %s", v)
	}
	t.Logf("runtime.Version() = %s", v)
}

// TestGOROOT は runtime.GOROOT() が空でなく、実際にディレクトリが存在することを検証する。
func TestGOROOT(t *testing.T) {
	goroot := runtime.GOROOT()
	if goroot == "" {
		t.Fatal("runtime.GOROOT() returned empty string")
	}

	info, err := os.Stat(goroot)
	if err != nil {
		t.Fatalf("GOROOT directory does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("GOROOT is not a directory: %s", goroot)
	}
	t.Logf("runtime.GOROOT() = %s", goroot)
}

// TestGoBinaryExists は GOROOT/bin/go と GOROOT/bin/gofmt の存在を検証する。
func TestGoBinaryExists(t *testing.T) {
	goroot := runtime.GOROOT()
	binaries := []string{"go", "gofmt"}

	for _, name := range binaries {
		path := filepath.Join(goroot, "bin", name)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("binary not found: %s: %v", path, err)
			continue
		}
		// 実行可能ビットが立っているか
		if info.Mode()&0111 == 0 {
			t.Errorf("binary is not executable: %s (mode=%s)", path, info.Mode())
		}
		t.Logf("found: %s (size=%d, mode=%s)", path, info.Size(), info.Mode())
	}
}

// TestGOPATH は GOPATH 環境変数が設定されていることを検証する。
func TestGOPATH(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Log("GOPATH is not set via env; using default")
		return
	}
	t.Logf("GOPATH = %s", gopath)

	info, err := os.Stat(gopath)
	if err != nil {
		t.Logf("GOPATH directory does not exist yet: %v", err)
		return
	}
	if !info.IsDir() {
		t.Errorf("GOPATH is not a directory: %s", gopath)
	}
}

// TestGOROOTLayout は GOROOT 配下の主要ディレクトリが存在することを検証する。
func TestGOROOTLayout(t *testing.T) {
	goroot := runtime.GOROOT()
	expectedDirs := []string{
		"bin",
		"src",
		"pkg",
		"lib",
	}

	for _, dir := range expectedDirs {
		path := filepath.Join(goroot, dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected directory not found: %s: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected directory but got file: %s", path)
		}
		t.Logf("found: %s/", path)
	}
}

// TestRuntimeInfo は runtime パッケージから取得可能な基本情報を記録する。
func TestRuntimeInfo(t *testing.T) {
	t.Logf("GOOS      = %s", runtime.GOOS)
	t.Logf("GOARCH    = %s", runtime.GOARCH)
	t.Logf("Compiler  = %s", runtime.Compiler)
	t.Logf("NumCPU    = %d", runtime.NumCPU())
	t.Logf("Version   = %s", runtime.Version())
	t.Logf("GOROOT    = %s", runtime.GOROOT())
}
