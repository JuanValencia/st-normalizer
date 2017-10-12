# Simple Url Normalizer  

### Usage

# dependencies
- go get github.com/sharethis-github/st-normalizer

# example use in code
```

package main

import (
  "fmt"
  "github.com/sharethis-github/st-normalizer"
)

func main() {
  testUrl := "https://www.example.com?p=3678"
  normalizerData, err := normalizer.NewNormalizer(testUrl)
  if err != nil {
    fmt.Println(err.Error())
    panic(err)
  }

  normalizerData.Normalize()
  fmt.Println(normalizerData)
}

```

# output of above example code
```
raw_url: https://www.example.com?p=3678
protocol: https
raw_query_params: p=3678
canonical_url: https://www.example.com
canonical_url_hash: e149be135a8b6803951f75776d589aaa
url_identifier: www.example.com?p=3678
url_identifier_hash: 600ef3b84b8ebc3db724c3d3e1aff542
```
