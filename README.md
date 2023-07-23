# watchpkgsite

Run this while inside your repo and it'll run pkgsite and keep fetching updates from the remote and rerunning pkgsite when there's something new.

## Installation

```
$ go install github.com/xrash/watchpkgsite/cmd/watchpkgsite@latest
```

## Running

```
$ watchpkgsite --addr :8765
```

Run point to another repo:

```
$  watchpkgsite --workdir /path/to/repo --addr :8765
```

Help:

```
$  watchpkgsite --help
```

