usage
```
package main

import (
	"fmt"
	cxlockrw "github.com/cloudxaas/golock/rw"
	"runtime"
)

func main() {
	numShards := runtime.NumCPU() * 64
	lock := cxlockrw.NewShardedRWLock(numShards)
	defer lock.Close()

	// Example usage of the sharded read-write lock
	key := "exampleKey"
	lock.RLock(key)
	fmt.Println("Read operation under RLock")
	lock.RUnlock(key)

	lock.Lock(key)
	fmt.Println("Write operation under Lock")
	lock.Unlock(key)
}
```
