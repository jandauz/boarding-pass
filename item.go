package bcbp

// item represents an item in the IATA 729 Bar Coded Boarding Pass specification.
//
// An item can belong to one of 3 main categories:
//   Mandatory
//   Conditional
//   Security
//
// An item also is either unique or repeated. A unique item only appears once,
// whereas a repeated item appears once for each flight segment.
//
// The length of an item dictates how many characters belong to that field.
// For example, PassengerName is the second field in a Bar Coded Boarding Pass
// and is 20 characters long.
type item struct {
    id          itemID
    description string
    length      int
    items       []item
}

type itemID uint

const (
    formatCode itemID = iota
    numberOfLegsEncoded
    passengerName
    electronicTicketIndicator
    operatingCarrierPNRCode
    fromCityAirportCode
    toCityAirportCode
    operatingCarrierDesignator
    flightNumber
    dateOfFlight
    compartmentCode
    seatNumber
    checkinSequenceNumber
    passengerStatus
    fieldSizeOfVariableSizeField
    beginningOfVersionNumber
    versionNumber
    fieldSizeOfFollowingStructuredMessageUnique
    passengerDescription
    sourceOfCheckin
    sourceOfBoardingPassIssuance
    dateOfIssueOfBoardingPass
    documentType
    airlineDesignatorOfBoardingPassIssuer
    baggageTagLicensePlateNumber
    firstNonConsecutiveBaggageTagLicensePlateNumber
    secondNonConsecutiveBaggageTagLicensePlateNumber
    fieldSizeOfFollowingStructuredMessageRepeated
    airlineNumericCode
    documentFormSerialNumber
    selecteeIndicator
    internationalDocumentationVerification
    marketingCarrierDesignator
    frequentFlyerAirlineDesignator
    frequentFlyerNumber
    idadIndicator
    freeBaggageAllowance
    fastTrack
    forIndividualAirlineUse
    beginningOfSecurityData
    typeOfSecurityData
    lengthOfSecurityData
    securityData
)

// spec is a graph of items that dictates how a Bar Coded Boarding Pass is
// processed.
var spec = []item{
    {
        id:          formatCode,
        description: "Format Code",
        length:      1,
    },
    {
        id:          numberOfLegsEncoded,
        description: "Number of Legs Encoded",
        length:      1,
    },
    {
        id:          passengerName,
        description: "Passenger Name",
        length:      20,
    },
    {
        id:          electronicTicketIndicator,
        description: "Electronic Ticket Indicator",
        length:      1,
    },
    {
        id:          operatingCarrierPNRCode,
        description: "Operating Carrier PNR Code",
        length:      7,
    },
    {
        id:          fromCityAirportCode,
        description: "From City Airport Code",
        length:      3,
    },
    {
        id:          toCityAirportCode,
        description: "To City Airport Code",
        length:      3,
    },
    {
        id:          operatingCarrierDesignator,
        description: "Operating Carrier Designator",
        length:      3,
    },
    {
        id:          flightNumber,
        description: "Flight Number",
        length:      5,
    },
    {
        id:          dateOfFlight,
        description: "Date of Flight (Julian Date)",
        length:      3,
    },
    {
        id:          compartmentCode,
        description: "Compartment Code",
        length:      1,
    },
    {
        id:          seatNumber,
        description: "Seat Number",
        length:      4,
    },
    {
        id:          checkinSequenceNumber,
        description: "Check-in Sequence Number",
        length:      5,
    },
    {
        id:          passengerStatus,
        description: "Passenger Status",
        length:      1,
    },
    {
        id:          fieldSizeOfVariableSizeField,
        description: "Field Size of variable size field",
        length:      2,
        items: []item{
            {
                id:          beginningOfVersionNumber,
                description: "Beginning of version number",
                length:      1,
            },
            {
                id:          versionNumber,
                description: "Version Number",
                length:      1,
            },
            {
                id:          fieldSizeOfFollowingStructuredMessageUnique,
                description: "Field Size of following structured message - unique",
                length:      2,
                items: []item{
                    {
                        id:          passengerDescription,
                        description: "Passenger Description",
                        length:      1,
                    },
                    {
                        id:          sourceOfCheckin,
                        description: "Source of check-in",
                        length:      1,
                    },
                    {
                        id:          sourceOfBoardingPassIssuance,
                        description: "Source of Boarding Pass Issuance",
                        length:      1,
                    },
                    {
                        id:          dateOfIssueOfBoardingPass,
                        description: "Date of Issue of Boarding Pass (Julian Date)",
                        length:      4,
                    },
                    {
                        id:          documentType,
                        description: "Document Type",
                        length:      1,
                    },
                    {
                        id:          airlineDesignatorOfBoardingPassIssuer,
                        description: "Airline Designator of boarding pass issuer",
                        length:      3,
                    },
                    {
                        id:          baggageTagLicensePlateNumber,
                        description: "Baggage Tag License Plate Number(s)",
                        length:      13,
                    },
                    {
                        id:          firstNonConsecutiveBaggageTagLicensePlateNumber,
                        description: "1st Non-Consecutive Baggage Tag License Plate Number",
                        length:      13,
                    },
                    {
                        id:          secondNonConsecutiveBaggageTagLicensePlateNumber,
                        description: "2nd Non-Consecutive Baggage Tag License Plate Number",
                        length:      13,
                    },
                },
            },
            {
                id:          fieldSizeOfFollowingStructuredMessageRepeated,
                description: "Field Size of following structured message - repeated",
                length:      2,
                items: []item{
                    {
                        id:          airlineNumericCode,
                        description: "Airline Numeric Code",
                        length:      3,
                    },
                    {
                        id:          documentFormSerialNumber,
                        description: "Document Form/Serial Number",
                        length:      10,
                    },
                    {
                        id:          selecteeIndicator,
                        description: "Selectee Indicator",
                        length:      1,
                    },
                    {
                        id:          internationalDocumentationVerification,
                        description: "International Documentation Verification",
                        length:      1,
                    },
                    {
                        id:          marketingCarrierDesignator,
                        description: "Marketing Carrier Designator",
                        length:      3,
                    },
                    {
                        id:          frequentFlyerAirlineDesignator,
                        description: "Frequent Flyer Airline Designator",
                        length:      3,
                    },
                    {
                        id:          frequentFlyerNumber,
                        description: "Frequent Flyer Number",
                        length:      16,
                    },
                    {
                        id:          idadIndicator,
                        description: "ID/AD Indicator",
                        length:      1,
                    },
                    {
                        id:          freeBaggageAllowance,
                        description: "Free Baggage Allowance",
                        length:      3,
                    },
                    {
                        id:          fastTrack,
                        description: "Fast Track",
                        length:      1,
                    },
                },
            },
            {
                id:          forIndividualAirlineUse,
                description: "For individual airline use",
            },
        },
    },
    {
        id:          beginningOfSecurityData,
        description: "Beginning of Security Data",
        length:      1,
    },
    {
        id:          typeOfSecurityData,
        description: "Type of Security Data",
        length:      1,
    },
    {
        id:          lengthOfSecurityData,
        description: "Length of Security Data",
        length:      2,
    },
    {
        id:          securityData,
        description: "Security Data",
        length:      100,
    },
}
