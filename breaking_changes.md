# Breaking Changes

This document provides a list of breaking changes since the V1 release
of the stomp client library.

## 1. No longer using gopkg.in

Version 1 of the library used Gustavo Niemeyer's gopkg.in facility for versioning Go libraries.
For a number of reasons, the stomp library no longer uses this facility. For this reason the
import path has changed.

Version 1:
```go
import (
    "gopkg.in/stomp.v1"
)
```

Version 2:
```go
    "github.com/go-stomp/stomp"
```

