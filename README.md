It defaults the logger to the stdlib log implementation. 

if LOG_TO_FILE is true ,it supoorts file log.

## Config

Set the logger level , flags...

```go
// import micro/go-log
import "github.com/XXX/go-log"

// SetLogger expects github.com/go-log/log.Logger interface
log.SetLogger(mylogger)
log.SetLevel(...)
log.SetFlags(...)
// Set log file path
log.SetPath(...)
```
