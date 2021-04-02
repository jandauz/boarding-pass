# boarding-pass
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/jandauz/boarding-pass/Go)
![Codecov](https://img.shields.io/codecov/c/github/jandauz/boarding-pass)
[![Go Report Card](https://goreportcard.com/badge/github.com/jandauz/boarding-pass)](https://goreportcard.com/report/github.com/jandauz/boarding-pass)
[![Go Reference](https://pkg.go.dev/badge/github.com/jandauz/boarding-pass/bcbp.svg)](https://pkg.go.dev/github.com/jandauz/boarding-pass/bcbp)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjandauz%2Fboarding-pass.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjandauz%2Fboarding-pass?ref=badge_shield)

boarding-pass is a partial [Go](https://golang.org) port of
[georgesmith64/bcbp](https://github.com/georgesmith46/bcbp) with some
inspiration from [martinmroz/iata_bcbp](https://github.com/martinmroz/iata_bcbp).

The main difference between the libraries is that
[boarding-pass](https://github.com/jandauz/boarding-pass) only offers users the
ability to decode the data in a Bar Coded Boarding Pass into structured data.

From an implementation perspective, [georgesmith64/bcbp](https://github.com/georgesmith46/bcbp)
is written in JavaScript and thus has the flexibility to create the structured
data on the fly. [martinmroz/iata_bcbp](https://github.com/martinmroz/iata_bcbp)
is written in Rust and uses `impl` to expose the data through method calls.
[boarding-pass](https://github.com/jandauz/boarding-pass) aims to build the
[BCBP struct](https://pkg.go.dev/github.com/jandauz/boarding-pass/bcbp#BCBP)
on the fly as efficiently as possible and return that to the caller.

[boarding-pass](https://github.com/jandauz/boarding-pass) supports up to
version 8 of the IATA Resolution 792 spec. Although, both
[georgesmith64/bcbp](https://github.com/georgesmith46/bcbp) and
[martinmroz/iata_bcbp](https://github.com/martinmroz/iata_bcbp) support
up to version 6, information on version 7 and 8 were gleaned from this
[issue](https://github.com/georgesmith46/bcbp/issues/3).

## Installation
```bash
go get github.com/jandauz/boarding-pass
```

## Getting started

```go
package main

import (
	"fmt"

	"github.com/jandauz/boarding-pass"
)

func main() {
	const s = "M1DESMARAIS/LUC       EABC123 YULFRAAC 0834 326J001A0025 100"
	b, err := bcbp.FromStr(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Format Code:", b.FormatCode)
	fmt.Println("NumberOfLegsEncoded:", b.NumberOfLegsEncoded)
	fmt.Println("PassengerName:", b.PassengerName)
	fmt.Println("ElectronicTicketIndicator:", b.ElectronicTicketIndicator)
	fmt.Println("OperatingCarrierPNRCode:", b.Legs[0].OperatingCarrierPNRCode)
	fmt.Println("FromCityAirportCode:", b.Legs[0].FromCityAirportCode)
	fmt.Println("ToCityAirportCode:", b.Legs[0].ToCityAirportCode)
	fmt.Println("OperatingCarrierDesignator:", b.Legs[0].OperatingCarrierDesignator)
	fmt.Println("FlightNumber:", b.Legs[0].FlightNumber)
	fmt.Println("DateOfFlight:", b.Legs[0].DateOfFlight)
	fmt.Println("CompartmentCode:", b.Legs[0].CompartmentCode)
	fmt.Println("SeatNumber:", b.Legs[0].SeatNumber)
	fmt.Println("CheckInSequenceNumber:", b.Legs[0].CheckInSequenceNumber)
	fmt.Println("PassengerStatus:", b.Legs[0].PassengerStatus)
	// Output:
	// Format Code: M
	// NumberOfLegsEncoded: 1
	// PassengerName: DESMARAIS/LUC
	// ElectronicTicketIndicator: E
	// OperatingCarrierPNRCode: ABC123
	// FromCityAirportCode: YUL
	// ToCityAirportCode: FRA
	// OperatingCarrierDesignator: AC
	// FlightNumber: 0834
	// DateOfFlight: 2021-11-22
	// CompartmentCode: J
	// SeatNumber: 001A
	// CheckInSequenceNumber: 0025
	// PassengerStatus: 1
}
```

## Notes
[boarding-pass](https://github.com/jandauz/boarding-pass) currently does not
attempt to interpret the data except for`NumberOfLegsEncoded`, `DateOfFlight`,
and `DateOfBoardingPassIssuance`.

`NumberOfLegsEncoded` is a `uint`. This is necessary for determing how many
`Legs` to process.

Both `DateOfFlight` and `DateOfBoardingPassIssuance` are
strings formatted using
[RFC 3339 full-date format](https://tools.ietf.org/html/rfc3339#section-5.6).
There is currently no attempt to determine if the values are realistic dates
e.g. an unrealistic date would be on that is far ahead in the future.

## Benchmark
```bash
goos: windows
goarch: amd64
pkg: github.com/jandauz/boarding-pass
cpu: AMD Ryzen 5 3600 6-Core Processor
BenchmarkFromStr_Mandatory_No_Security_Single
BenchmarkFromStr_Mandatory_No_Security_Single-12    	 1398615	       851.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkFromStr_Mandatory_Single
BenchmarkFromStr_Mandatory_Single-12                	 1231972	       973.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkFromStr_Full_No_Security_Single
BenchmarkFromStr_Full_No_Security_Single-12         	  722799	      1714 ns/op	      16 B/op	       1 allocs/op
BenchmarkFromStr_Full_Single
BenchmarkFromStr_Full_Single-12                     	  663570	      1863 ns/op	      16 B/op	       1 allocs/op
BenchmarkFromStr_Full_No_Security_Multi
BenchmarkFromStr_Full_No_Security_Multi-12          	  455092	      2559 ns/op	      16 B/op	       1 allocs/op
BenchmarkFromStr_Full_Multi
BenchmarkFromStr_Full_Multi-12                      	  455092	      2682 ns/op	      16 B/op	       1 allocs/op
PASS
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjandauz%2Fboarding-pass.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjandauz%2Fboarding-pass?ref=badge_large)