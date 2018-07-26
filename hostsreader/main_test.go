package hostsreader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {

	// Create a temporary directory
	p, err := ioutil.TempDir(os.TempDir(), "prefix")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	// Create a temporary file
	file := filepath.Join(p, "ssh_file")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	// With some content
	err = ioutil.WriteFile(file, []byte("\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write to temporary file")
	}

	// Create a helper
	obj := New(file)

	// Parse the file
	out, err := obj.Hosts()
	if err != nil {
		t.Errorf("Error parsing (empty) file :%s", err.Error())
	}
	if len(out) != 0 {
		t.Errorf("Empty file should have no hosts!")
	}
	// Cleanup
	os.RemoveAll(p)
}

func TestContent(t *testing.T) {

	// Create a temporary directory
	p, err := ioutil.TempDir(os.TempDir(), "prefix")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	// Create a temporary file
	file := filepath.Join(p, "ssh_file")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	lines := `
[deagol.vpn]:2222,[10.10.10.100]:2222 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH+761batAEA5KM7JQUrKeNyKftdnRd49E03snPA/j8nP6u7vJlIzf9S2MZlbZyHeh5Hr2wIVwpJF1n5ycg1rG4=
[mail.steve.org.uk]:2222 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH+761batAEA5KM7JQUrKeNyKftdnRd49E03snPA/j8nP6u7vJlIzf9S2MZlbZyHeh5Hr2wIVwpJF1n5ycg1rG4=
`

	// With some content
	err = ioutil.WriteFile(file, []byte(lines), 0644)
	if err != nil {
		t.Errorf("Failed to write to temporary file")
	}

	// Create a helper
	obj := New(file)

	// Parsed results
	var out map[string]string

	// Parse the file
	out, err = obj.Hosts()
	if err != nil {
		t.Errorf("Error parsing file :%s", err.Error())
	}

	// Rewrite some content - to bump the timestamp
	err = ioutil.WriteFile(file, []byte(lines), 0644)
	if err != nil {
		t.Errorf("Failed to write to temporary file")
	}

	out, err = obj.Hosts()
	if err != nil {
		t.Errorf("Error parsing file :%s", err.Error())
	}

	// Ensure we got the expected values
	if len(out) != 3 {
		t.Errorf("File should have two hosts, got %d!", len(out))
	}

	// Now look for the fingerprint matches
	fing := "c8:ca:c2:44:66:98:31:2a:c1:c2:91:e0:fc:b3:91:b2"

	if out["deagol.vpn"] != fing {
		t.Errorf("Key mismatch %s != %s", out["deagol.vpn"], fing)
	}
	if out["mail.steve.org.uk"] != fing {
		t.Errorf("Key mismatch %s != %s", out["mail.steve.org.uk"], fing)
	}
	if out["10.10.10.100"] != fing {
		t.Errorf("Key mismatch %s != %s", out["10.10.10.100"], fing)
	}

	// Cleanup
	os.RemoveAll(p)
}

func TestBogus(t *testing.T) {

	// Create a temporary directory
	p, err := ioutil.TempDir(os.TempDir(), "prefix")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	// Create a temporary file
	file := filepath.Join(p, "ssh_file")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	lines := `
moi kissa tes
`

	// With some content
	err = ioutil.WriteFile(file, []byte(lines), 0644)
	if err != nil {
		t.Errorf("Failed to write to temporary file")
	}

	// Create a helper
	obj := New(file)

	// Parse the file
	_, err = obj.Hosts()
	if err == nil {
		t.Errorf("Expected an error parsing a bogus file, didn't get one!")
	}

	// Cleanup
	os.RemoveAll(p)
}

// Parse a file, remove it, and attempt re-parsing
func TestRemovingFile(t *testing.T) {

	// Create a temporary directory
	p, err := ioutil.TempDir(os.TempDir(), "prefix")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	// Create a temporary file
	file := filepath.Join(p, "ssh_file")
	if err != nil {
		t.Errorf("Error creating temporary directory:%s", err.Error())
	}

	lines := `
[deagol.vpn]:2222,[10.10.10.100]:2222 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH+761batAEA5KM7JQUrKeNyKftdnRd49E03snPA/j8nP6u7vJlIzf9S2MZlbZyHeh5Hr2wIVwpJF1n5ycg1rG4=
`

	// With some content
	err = ioutil.WriteFile(file, []byte(lines), 0644)
	if err != nil {
		t.Errorf("Failed to write to temporary file")
	}

	// Create a helper
	obj := New(file)

	// Parse the file
	_, err = obj.Hosts()
	if err != nil {
		t.Errorf("Expected no error parsing a bogus file, but get one : %s", err.Error())
	}

	// Cleanup
	os.RemoveAll(p)

	// Parse the file again - now it is deleted.
	_, err = obj.Hosts()
	if err == nil {
		t.Errorf("Expected an error parsing a removed file - got none")
	}
}
