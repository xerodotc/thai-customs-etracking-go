# thai-customs-etracking-go

Unofficial Go API for checking custom clearance status and data on
Thai Customs e-Tracking System (based on web scraping)

## Installation
1. Use command below to install an API
```
$ go get -u github.com/xerodotc/thai-customs-etracking-go
```
2. Import into your code
```go
import "github.com/xerodotc/thai-customs-etracking-go"
```

## Usage / Quickstart

Use `etracking.Lookup("EAXXXXXXXXXJP")` to retrieve item's custom clearance status and data
```go
data, err := etracking.Lookup("EAXXXXXXXXXJP")
```

## License / Copyright
This project is licensed under MIT license.
See [LICENSE](LICENSE) file for details.

The author of this project is not related to Thai Customs in anyway.

Thai Customs e-Tracking is a trademark of Thai Customs Department
