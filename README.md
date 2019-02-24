[![Go Report Card](https://goreportcard.com/badge/github.com/skx/knownfs)](https://goreportcard.com/report/github.com/skx/knownfs)
[![license](https://img.shields.io/github/license/skx/knownfs.svg)](https://github.com/skx/knownfs/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/skx/knownfs.svg)](https://github.com/skx/knownfs/releases/latest)


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

You should see a single subdirectory for each hostname listed in your `~/.ssh/known_hosts` file, and inside the directory you'll find a file named `fingerprint` with the hosts' fingerprint.

Once you're done you will need to unmount the mount-point.  If you have `fusermount` installed you can do so like this:

      $ fusermount -u ~/knownfs/

If not you'll need root permissions to unmount the end-point:

      $ sudo umount ~/knownfs/


# Options

By default you'll see entries for each host found, whether those entries are hostnames or IP addresses.  For example on my own system I see this:

      frodo ~/knownfs $ ls | head -n 5
      10.0.0.10
      10.10.10.100
      10.10.10.20
      10.10.10.97
      10.10.10.98

I prefer to only view _real_ hosts, so I exclude IPv4/IPv6-based entries like so:

      $ knownfs -hosts-only

That gives me just hostnames:

      $ ls | head -5
      blogspam.blogspam.net
      blogspam.net
      blog.steve.fi
      builder.steve.org.uk
      builder.vpn

You can also specify the path to an alternative `known_hosts` file, with `-config /path/to/file`.


## Github Setup

This repository is configured to run tests upon every commit, and when
pull-requests are created/updated.  The testing is carried out via
[.github/run-tests.sh](.github/run-tests.sh) which is used by the
[github-action-tester](https://github.com/skx/github-action-tester) action.

Releases are automated in a similar fashion via [.github/build](.github/build),
and the [github-action-publish-binaries](https://github.com/skx/github-action-publish-binaries) action.


Steve
--
