
# KnownFS

A simple FUSE-based filesystem which exports the contents of `~/.ssh/known_hosts`.

For every hostname listed in your known_hosts file this filesystem will create a directory, and that directory will contain a single file with the servers' fingerprint.


# Installation

     $ go get github.com/hanwen/go-fuse/...
     $ go get -u github.com/skx/knownfs
     $ go install github.com/skx/knownfs

# Usage:

Mount it:

     $ mkdir ~/knownfs/
     $ knownfs ~/knownfs/

In another window:

     $ ls -1 ~/knownfs/
