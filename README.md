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
RawUrl: https://www.example.com?p=3678
Protocol: https
RawQueryParams: p=3678
CanonicalUrl: https://www.example.com
CanonicalUrlHash: e149be135a8b6803951f75776d589aaa
UrlIdentifier: www.example.com?p=3678
UrlIdentifierHash: 600ef3b84b8ebc3db724c3d3e1aff54
```
test access
