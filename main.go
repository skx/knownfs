// This is the implementation of a simple FUSE-based filesystem.
//
// The contents of the filesystem will include one sub-directory
// for every host which is stored in your ~/.ssh/known_hosts file.
//
// Each directory will contain a further entry for the SSH fingerprint.
//
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/skx/knownfs/hostsreader"
)

// Our structure
type KnownFS struct {
	// The filesystem
	pathfs.FileSystem

	// The helper for reading the keys
	helper *hostsreader.HostReader
}

// GetAttr reads the attributes of the given file.
//
// Think of it as the implementation of stat() for the filesystem.
func (me *KnownFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {

	// Get entries.
	known, err := me.helper.Hosts()

	if err != nil {
		fmt.Printf("Error calling GetHosts")
		return nil, fuse.ENOENT
	}

	// Directory entry for a host?
	if known[name] != "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	// Otherwise if the entry is a hosts' fingerprint file then
	// return that it is a file & the correct size.
	for host, key := range known {
		if name == host+"/fingerprint" {
			return &fuse.Attr{
				Mode: fuse.S_IFREG | 0644, Size: uint64(len(key)),
			}, fuse.OK
		}
	}

	// Missing-file
	fmt.Printf("GetAttr(%s)\n", name)
	return nil, fuse.ENOENT
}

// OpenDir is called when a directory is opened, and should return the
// contents of that directory.
//
// We handle two cases: opening our mount-point, which involves creating
// one subdirectory for each host, and opening one of the per-host files
// which just involves creating a single fingerprint file.
func (me *KnownFS) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	var ret []fuse.DirEntry

	// top-level
	if name == "" {
		known, err := me.helper.Hosts()
		if err != nil {
			fmt.Printf("Error calling GetHosts")
			return nil, fuse.ENOENT
		}

		for host := range known {
			ret = append(ret, fuse.DirEntry{Name: host, Mode: fuse.S_IFDIR})
		}
		return ret, fuse.OK
	} else {

		// We assume we've opened a host.
		ret = append(ret, fuse.DirEntry{Name: "fingerprint", Mode: fuse.S_IFREG})
		return ret, fuse.OK
	}
}

// Open opens a file for reading/writing.
func (me *KnownFS) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

	// No writing is permitted
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	known, err := me.helper.Hosts()
	if err != nil {
		fmt.Printf("Error calling GetHosts")
		return nil, fuse.ENOENT
	}

	// Did we find the host?
	for host, key := range known {
		if name == host+"/fingerprint" {
			return nodefs.NewDataFile([]byte(key + "\n")), fuse.OK
		}
	}

	// Otherwise no entry.
	return nil, fuse.ENOENT
}

// Entry point.
func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n knownfs MOUNTPOINT")
	}
	nfs := pathfs.NewPathNodeFs(&KnownFS{FileSystem: pathfs.NewDefaultFileSystem(), helper: hostsreader.New()}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount failed: %v\n", err)
	}
	server.Serve()
}
