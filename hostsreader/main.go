package hostsreader

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// HostReader is the structure for our object.
type HostReader struct {
	// time holds the modification time of ~/.ssh/known_hosts
	time time.Time

	// the map of known-hosts and their fingerprints
	entries map[string]string
}

// New is our constructor.
func New() *HostReader {
	self := new(HostReader)
	self.entries = make(map[string]string)
	return self
}

// HasChanged returns true if the given file has changed since our
// last read - using the mtime of the file to decide.
//
// It allows us to avoid reparsing the file if the contents haven't
// changed.
func (self *HostReader) HasChanged() (bool, error) {
	file := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	data, err := os.Stat(file)
	if err != nil {
		return false, err
	}

	if data.ModTime().After(self.time) {
		return true, nil
	}
	return false, nil
}

// Hosts returns the map of known-hosts and their associated keys
//
// It caches accesses to the file via the `HasChanged` method to
// speed things up.
func (self *HostReader) Hosts() (map[string]string, error) {

	// Has the file changed recently?
	changed, err := self.HasChanged()
	if err != nil {
		return self.entries, err
	}

	// If not then we can return the entries - providing
	// we've parsed at least once recently.
	if changed == false && len(self.entries) > 0 {
		return self.entries, nil
	}

	// Here we might have been called because this is
	// our first invocation, or because the file has
	// changed.
	// Clear old entries in case of the latter.
	for k := range self.entries {
		delete(self.entries, k)
	}

	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return self.entries, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}

		// key, comment, hosts, ?? , err
		key, _, hosts, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
		if err != nil {
			return self.entries, err
		} else {
			for _, i := range hosts {
				host, _, _ := net.SplitHostPort(i)
				self.entries[host] = ssh.FingerprintLegacyMD5(key)
			}
		}
	}

	self.time = time.Now()
	return self.entries, nil
}
