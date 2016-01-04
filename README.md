

# How to build

Enable vendor experiment in Go:

```export GO15VENDOREXPERIMENT=1```

To build:

```go build```

# Adding Dependencies

First get the godep package

```go get github.com/tools/godep```

Use go get as usual to get the dependency. And then from the root project folder, run:

```export GO15VENDOREXPERIMENT=1```

And

```godep save```

# TODO
- Auto detection of platform
- Check in common if platform command not found
- Check /usr/share/ location by default and switch back to ~/.tldr if not found
- RPM, DEB packages
- Windows support

- better error handling
- Proxy support
