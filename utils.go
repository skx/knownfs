package main

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

// GetHosts returns the hostnames/IPs that are present in our users
// ~/.ssh/known_hosts file.
//
// Should return a map of:
//
//   map[hostname] -> key
//
// That way a) hosts are unique, and b) we have the keys.
func GetHosts() (map[string]string, error) {

	// The map we return as a result
	result := make(map[string]string)

	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return result, err
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
			return result, err
		} else {
			for _, i := range hosts {
				host, _, _ := net.SplitHostPort(i)
				result[host] = ssh.FingerprintLegacyMD5(key)
			}
		}
	}

	return result, err
}
