package libreoffice

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// patchedLibreOfficeVersion is the first trixie build with the fixes for
// CVE-2026-63260/63266/63267/63268/63269/63277 (DSA-6543-1). Keep it in sync
// with LIBREOFFICE_MIN_VERSION in build/Dockerfile.
const patchedLibreOfficeVersion = "4:25.2.3-2+deb13u8"

// These tests run inside the Gotenberg image and check what it ships, so an
// image rebuild cannot bring back the vulnerable LibreOffice or a JVM.

func TestImage_LibreOfficeIsPatched(t *testing.T) {
	out, err := exec.Command("dpkg-query", "-W", "-f=${Version}", "libreoffice-core").Output()
	if err != nil {
		t.Fatalf("query libreoffice-core version: %v", err)
	}

	installed := strings.TrimSpace(string(out))
	err = exec.Command("dpkg", "--compare-versions", installed, "ge", patchedLibreOfficeVersion).Run()
	if err != nil {
		t.Errorf("expected libreoffice-core >= '%s' but got '%s'", patchedLibreOfficeVersion, installed)
	}
}

// LibreOffice runs JDBC drivers named by a document inside a JVM, so the image
// must not contain one.
func TestImage_HasNoJavaRuntime(t *testing.T) {
	if path, err := exec.LookPath("java"); err == nil {
		t.Errorf("expected no java binary but found '%s'", path)
	}

	if _, err := os.Stat("/usr/lib/jvm"); err == nil {
		t.Error("expected no '/usr/lib/jvm' directory but it exists")
	}
}
