package hostsreader

import (
	"bufio"
	"net"
	"os"

	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// HostReader is the structure for our object.
type HostReader struct {
	// The file we're created with
	filename string

	// time holds the modification time of ~/.ssh/known_hosts
	time time.Time

	// the map of known-hosts and their fingerprints
	entries map[string]string
}

// New is our constructor.
func New(filename string) *HostReader {
	self := new(HostReader)
	self.filename = filename
	self.entries = make(map[string]string)
	return self
}

// HasChanged returns true if the given file has changed since our
// last read - using the mtime of the file to decide.
//
// It allows us to avoid reparsing the file if the contents haven't
// changed.
func (me *HostReader) HasChanged() (bool, error) {
	data, err := os.Stat(me.filename)
	if err != nil {
		return false, err
	}

	if data.ModTime().After(me.time) {
		return true, nil
	}
	return false, nil
}

// Hosts returns the map of known-hosts and their associated keys
//
// It caches accesses to the file via the `HasChanged` method to
// speed things up.
func (me *HostReader) Hosts() (map[string]string, error) {

	// Has the file changed recently?
	changed, err := me.HasChanged()
	if err != nil {
		return me.entries, err
	}

	// If not then we can return the entries - providing
	// we've parsed at least once recently.
	if !changed && len(me.entries) > 0 {
		return me.entries, nil
	}

	// Here we might have been called because this is
	// our first invocation, or because the file has
	// changed.
	// Clear old entries in case of the latter.
	for k := range me.entries {
		delete(me.entries, k)
	}

	file, err := os.Open(me.filename)
	if err != nil {
		return me.entries, err
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
			return me.entries, err
		}

		// For each host record the key against it
		for _, i := range hosts {

			// Split off the port if we need to
			if strings.Contains(i, ":") {
				host, _, err := net.SplitHostPort(i)
				if err == nil {
					me.entries[host] = ssh.FingerprintLegacyMD5(key)
				}
			}
		}
	}

	me.time = time.Now()
	return me.entries, nil
}
