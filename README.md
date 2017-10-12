# Simple Url Normalizer  

### Usage

# dependencies
- go get github.com/sharethis-github/st-normalizer

# use in code
```
testUrl := "https://www.example.com?p=3678"
normalizerData, err := st-normalizer.NewNormalizer(testUrl)
if err != nil {
  fmt.Println(err.Error())
  panic(err)
}

normalizerData.Normalize()
fmt.Println(normalizerData)
```
