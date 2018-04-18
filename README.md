
# KnownFS

A simple FUSE-based filesystem which exports the contents of `~/.ssh/known_hosts`.

For every hostname listed in your known_hosts file this filesystem will create a directory, and that directory will contain a single file with the servers' fingerprint.


# Installation

     $ go get -u github.com/skx/knownfs
     $ go install github.com/skx/knownfs

Now you should discover you have a binary installed at `$GOPATH/bin/knownfs`.


# Usage:

Make a directory for the filesystem, and mount it like so:

     $ mkdir ~/knownfs/
     $ knownfs ~/knownfs/

In another window:

     $ ls -1 ~/knownfs/
