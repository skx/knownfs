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

type KnownFS struct {
	pathfs.FileSystem
	helper *hostsreader.HostReader
}

// Get attributes of the named file/directory
func (me *KnownFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {

	// Get entries.
	known, err := me.helper.Hosts()

	if err != nil {
		fmt.Printf("Error calling GetHosts")
		return nil, fuse.ENOENT
	}

	// Directory entry for a host.
	if known[name] != "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	for host, key := range known {
		if name == host+"/fingerprint" {
			return &fuse.Attr{
				Mode: fuse.S_IFREG | 0644, Size: uint64(len(key)),
			}, fuse.OK
		}
	}

	fmt.Printf("GetAttr(%s)\n", name)
	return nil, fuse.ENOENT
}

// Open a directory (i.e. read the contents).
func (me *KnownFS) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	var ret []fuse.DirEntry

	// top-level
	if name == "" {
		known, err := me.helper.Hosts()
		if err != nil {
			fmt.Printf("Error calling GetHosts")
			return nil, fuse.ENOENT
		}

		for host, _ := range known {
			ret = append(ret, fuse.DirEntry{Name: host, Mode: fuse.S_IFDIR})
		}
		return ret, fuse.OK
	} else {

		// We assume we've opened a host.
		ret = append(ret, fuse.DirEntry{Name: "fingerprint", Mode: fuse.S_IFREG})
		return ret, fuse.OK
	}
	return nil, fuse.ENOENT
}

// Open a file.
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
