[![Build Status](https://travis-ci.org/anoopengineer/tldr.svg?branch=master)](https://travis-ci.org/anoopengineer/tldr)

# tldr client implemented in golang #
Native binary client that works in Linux, Mac OSX and Windows.

![Screenshot of TLDR](screenshots/tldr-screenshot.png)

# Installing #
Native binaries are available for Linux, Mac OSX and Windows in the [Release tab](https://github.com/anoopengineer/tldr/releases).

# How to build

```
make build
```

# How to create a release

```
make dist
```

# Saving new dependencies

First get the godep package

```
make godep-save
```


# TODO
- Better readme
- Auto detection of platform
- Check in common if platform command not found
- RPM, DEB packages
- Windows support

- better error handling
- Proxy support??
