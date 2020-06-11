# klocker

Golang Key Locker: lock by string key, so as to avoid giant lock.

### Usage

```go
package main

import (
	"github.com/rayzui/klocker"
)

func() main() {
  mu := &klocker.KMutex{}
  mu.Lock("kmutex")
  // Do something
  mu.Unlock("kmutex")
  
  rwmu := &klocker.RWKMutex{}
  rwmu.RLock("rwkmutex")
  // Get something
  rwmu.RUnlock("rwkmutex")
  rwmu.Lock("rwkmutex")
  // Update something
  rwmu.Unlock("rwkmutex")
}
```

