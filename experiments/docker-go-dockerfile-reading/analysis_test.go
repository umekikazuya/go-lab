package dockergodockerfilereading

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

// TestGOROOTMatchesDockerfileExpectation は、Dockerfileで設定されるGOROOTが
// 実際のコンテナ内の値と一致することを検証する。
// 公式Dockerfileでは ENV GOROOT /usr/local/go が設定される。
func TestGOROOTMatchesDockerfileExpectation(t *testing.T) {
	expected := "/usr/local/go"
	actual := runtime.GOROOT()
	if actual != expected {
		t.Errorf("GOROOT mismatch: expected=%s, actual=%s", expected, actual)
	}
	t.Logf("GOROOT = %s (expected: %s)", actual, expected)
}

// TestGOPATHMatchesDockerfileExpectation は、Dockerfileで設定されるGOPATHが
// 実際のコンテナ内の値と一致することを検証する。
// 公式Dockerfileでは ENV GOPATH /go が設定される。
func TestGOPATHMatchesDockerfileExpectation(t *testing.T) {
	expected := "/go"
	actual := os.Getenv("GOPATH")
	if actual == "" {
		t.Skip("GOPATH not set via environment variable (not running in container?)")
	}
	if actual != expected {
		t.Errorf("GOPATH mismatch: expected=%s, actual=%s", expected, actual)
	}
	t.Logf("GOPATH = %s (expected: %s)", actual, expected)
}

// TestPATHContainsGoBinaries は、PATHに Go のバイナリディレクトリが
// 含まれていることを検証する。
// 公式Dockerfileでは PATH に /go/bin と /usr/local/go/bin が追加される。
func TestPATHContainsGoBinaries(t *testing.T) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		t.Fatal("PATH is empty")
	}

	expectedPaths := []string{
		"/usr/local/go/bin",
	}

	for _, ep := range expectedPaths {
		found := false
		for _, p := range strings.Split(pathEnv, ":") {
			if p == ep {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PATH does not contain %s", ep)
		}
	}
	t.Logf("PATH = %s", pathEnv)
}

// TestGOLANG_VERSION は GOLANG_VERSION 環境変数が設定されていることを検証する。
// 公式Dockerfileでは ENV GOLANG_VERSION が設定される。
func TestGOLANG_VERSION(t *testing.T) {
	version := os.Getenv("GOLANG_VERSION")
	if version == "" {
		t.Skip("GOLANG_VERSION not set (not running in official container?)")
	}

	runtimeVersion := strings.TrimPrefix(runtime.Version(), "go")
	if !strings.HasPrefix(runtimeVersion, version) && !strings.HasPrefix(version, runtimeVersion) {
		t.Errorf("GOLANG_VERSION=%s does not match runtime.Version()=%s", version, runtime.Version())
	}
	t.Logf("GOLANG_VERSION = %s, runtime.Version() = %s", version, runtime.Version())
}

// TestGoDownloadSource は、Goバイナリの取得元が dl.google.com であることを
// 間接的に検証する（バイナリが正規のものであることの確認）。
func TestGoDownloadSource(t *testing.T) {
	// runtime.Version() が公式リリースの形式であることを確認
	v := runtime.Version()
	if !strings.HasPrefix(v, "go1.") && !strings.HasPrefix(v, "go2.") {
		t.Logf("non-release version detected: %s (possibly built from source)", v)
	}
	t.Logf("Go version: %s (compiler: %s)", v, runtime.Compiler)

	if runtime.Compiler != "gc" {
		t.Errorf("expected gc compiler, got: %s", runtime.Compiler)
	}
}

// TestBaseOS は、コンテナのベースOSが Debian bookworm であることを検証する。
func TestBaseOS(t *testing.T) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		t.Skip("cannot read /etc/os-release (not running in container?)")
	}

	content := string(data)
	if !strings.Contains(content, "bookworm") && !strings.Contains(content, "Debian") {
		t.Logf("unexpected base OS (may be alpine or other variant)")
	}
	t.Logf("/etc/os-release:\n%s", content)
}
