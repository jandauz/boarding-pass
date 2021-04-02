package bcbp_test

import (
    "fmt"

    "github.com/jandauz/boarding-pass"
)

func ExampleFromStr() {
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
