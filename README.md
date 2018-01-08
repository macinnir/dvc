# DVC

# Usage 

```

import (
    "github.com/macinnir/dvc"
)

func main() {
    d := dvc.DVC{}
    d.Run("path/to/changeset/files", "databaseHost:port", "username", "password")
}
```