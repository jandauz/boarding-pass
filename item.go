package bcbp

import "regexp"

// item represents an item in the IATA 729 Bar Coded Boarding Pass specification.
//
// An item can belong to one of 3 main categories:
//
//	Mandatory
//	Conditional
//	Security
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
	format      string
	regex       *regexp.Regexp
	items       []item
}

// validate validates s against item.regex.
func (i item) validate(s string) bool {
	if i.id == beginningOfSecurityData && len(s) == 0 {
		return true
	}
	return i.regex.FindString(s) != ""
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
		format:      `"M"`,
		regex:       formatCodeRegex,
	},
	{
		id:          numberOfLegsEncoded,
		description: "Number of Legs Encoded",
		length:      1,
		format:      "a number between 1 to 4",
		regex:       numberOfLegsEncodedRegex,
	},
	{
		id:          passengerName,
		description: "Passenger Name",
		length:      20,
		format:      `20 characters with trailing whitespaces where the last name must be at most 18 characters followed by "/" and an alpha initial`,
		regex:       passengerNameRegex,
	},
	{
		id:          electronicTicketIndicator,
		description: "Electronic Ticket Indicator",
		length:      1,
		format:      "E or L",
		regex:       electronicTicketRegex,
	},
	{
		id:          operatingCarrierPNRCode,
		description: "Operating Carrier PNR Code",
		length:      7,
		format:      "7 alphanumeric characters with trailing whitespaces",
		regex:       operatingCarrierPNRCodeRegex,
	},
	{
		id:          fromCityAirportCode,
		description: "From City Airport Code",
		length:      3,
		format:      "3 alpha characters",
		regex:       airportCodeRegex,
	},
	{
		id:          toCityAirportCode,
		description: "To City Airport Code",
		length:      3,
		format:      "3 alpha characters",
		regex:       airportCodeRegex,
	},
	{
		id:          operatingCarrierDesignator,
		description: "Operating Carrier Designator",
		length:      3,
		format:      "3 alphanumeric characters with trailing whitespaces",
		regex:       operatingCarrierDesignatorRegex,
	},
	{
		id:          flightNumber,
		description: "Flight Number",
		length:      5,
		format:      "4 digits with leading zeroes followed by an optional alpha suffix or whitespace",
		regex:       flightNumberRegex,
	},
	{
		id:          dateOfFlight,
		description: "Date of Flight (Julian Date)",
		length:      3,
		format:      "3 digits with leading zeroes with maximum value of 365 (366 for leap years)",
		regex:       dateOfFlightRegex,
	},
	{
		id:          compartmentCode,
		description: "Compartment Code",
		length:      1,
		format:      "an alpha character",
		regex:       compartmentCodeRegex,
	},
	{
		id:          seatNumber,
		description: "Seat Number",
		length:      4,
		format:      "3 digits with leading zeroes followed by an alpha",
		regex:       seatNumberRegex,
	},
	{
		id:          checkinSequenceNumber,
		description: "Check-in Sequence Number",
		length:      5,
		format:      "4 digits with leading zeroes followed by an optional alpha or whitespace",
		regex:       checkInSequenceNumberRegex,
	},
	{
		id:          passengerStatus,
		description: "Passenger Status",
		length:      1,
		format:      "an alphanumeric character",
		regex:       passengerStatusRegex,
	},
	{
		id:          fieldSizeOfVariableSizeField,
		description: "Field Size of variable size field",
		length:      2,
		format:      "a hex number with leading zeroes",
		regex:       hexRegex,
		items: []item{
			{
				id:          beginningOfVersionNumber,
				description: "Beginning of version number",
				length:      1,
				format:      `">"`,
				regex:       beginningOfVersionNumberRegex,
			},
			{
				id:          versionNumber,
				description: "Version Number",
				length:      1,
				format:      "a number between 1 and 8",
				regex:       versionNumberRegex,
			},
			{
				id:          fieldSizeOfFollowingStructuredMessageUnique,
				description: "Field Size of following structured message - unique",
				length:      2,
				format:      "a hex number with leading zeroes",
				regex:       hexRegex,
				items: []item{
					{
						id:          passengerDescription,
						description: "Passenger Description",
						length:      1,
						format:      "an alphanumeric character",
						regex:       passengerDescriptionRegex,
					},
					{
						id:          sourceOfCheckin,
						description: "Source of check-in",
						length:      1,
						format:      "W, K, X, R, M, O, T, V, A, or whitespace",
						regex:       sourceOfCheckInRegex,
					},
					{
						id:          sourceOfBoardingPassIssuance,
						description: "Source of Boarding Pass Issuance",
						length:      1,
						format:      "W, K, X, R, M, O, T, V, or whitespace",
						regex:       sourceOfBoardingPassIssuanceRegex,
					},
					{
						id:          dateOfIssueOfBoardingPass,
						description: "Date of Issue of Boarding Pass (Julian Date)",
						length:      4,
						format:      "4 digits with leading zeroes with last 3 digits having maximum value of 365 (366 for leap years)",
						regex:       dateOfIssueOfBoardingPassRegex,
					},
					{
						id:          documentType,
						description: "Document Type",
						length:      1,
						format:      "B, I, or whitespace",
						regex:       documentTypeRegex,
					},
					{
						id:          airlineDesignatorOfBoardingPassIssuer,
						description: "Airline Designator of boarding pass issuer",
						length:      3,
						format:      "left justified 3 alphanumeric characters with trailing whitespaces",
						regex:       airlineDesignatorOfBoardingPassIssuerRegex,
					},
					{
						id:          baggageTagLicensePlateNumber,
						description: "Baggage Tag License Plate Number(s)",
						length:      13,
						// IATA 792 spec states that this field is alphanumeric
						// however the interpretation of the data shows it to be
						// numeric only.
						format: "13 digits",
						regex:  baggageTagLicensePlateNumberRegex,
					},
					{
						id:          firstNonConsecutiveBaggageTagLicensePlateNumber,
						description: "1st Non-Consecutive Baggage Tag License Plate Number",
						length:      13,
						// IATA 792 spec states that this field is alphanumeric
						// however the interpretation of the data shows it to be
						// numeric only.
						format: "13 digits",
						regex:  baggageTagLicensePlateNumberRegex,
					},
					{
						id:          secondNonConsecutiveBaggageTagLicensePlateNumber,
						description: "2nd Non-Consecutive Baggage Tag License Plate Number",
						length:      13,
						// IATA 792 spec states that this field is alphanumeric
						// however the interpretation of the data shows it to be
						// numeric only.
						format: "13 numeric characters",
						regex:  baggageTagLicensePlateNumberRegex,
					},
				},
			},
			{
				id:          fieldSizeOfFollowingStructuredMessageRepeated,
				description: "Field Size of following structured message - repeated",
				length:      2,
				format:      "a hex number with leading zeroes",
				regex:       hexRegex,
				items: []item{
					{
						id:          airlineNumericCode,
						description: "Airline Numeric Code",
						length:      3,
						format:      "3 digits with leading zeroes",
						regex:       airlineNumericCodeRegex,
					},
					{
						id:          documentFormSerialNumber,
						description: "Document Form/Serial Number",
						length:      10,
						format:      "10 alphanumeric characters with leading zeroes",
						regex:       documentFormSerialNumberRegex,
					},
					{
						id:          selecteeIndicator,
						description: "Selectee Indicator",
						length:      1,
						format:      "0, 1, 2, or whitespace",
						regex:       selecteeIndicatorRegex,
					},
					{
						id:          internationalDocumentationVerification,
						description: "International Documentation Verification",
						length:      1,
						format:      "0, 1, 2, or whitespace",
						regex:       internationalDocumentationVerificationRegex,
					},
					{
						id:          marketingCarrierDesignator,
						description: "Marketing Carrier Designator",
						length:      3,
						format:      "3 alphanumeric characters with trailing whitespaces",
						regex:       marketingCarrierDesignatorRegex,
					},
					{
						id:          frequentFlyerAirlineDesignator,
						description: "Frequent Flyer Airline Designator",
						length:      3,
						format:      "3 alphanumeric characters with trailing whitespaces",
						regex:       frequentFlyerAirlineDesignatorRegex,
					},
					{
						id:          frequentFlyerNumber,
						description: "Frequent Flyer Number",
						length:      16,
						format:      "16 alphanumeric characters with trailing whitespaces",
						regex:       frequentFlyerNumberRegex,
					},
					{
						id:          idadIndicator,
						description: "ID/AD Indicator",
						length:      1,
						format:      "an alphanumeric character or whitespace",
						regex:       idadIndicatorRegex,
					},
					{
						id:          freeBaggageAllowance,
						description: "Free Baggage Allowance",
						length:      3,
						format:      "2 digits with leading zeroes followed by K or L; or 1 digit followed by PC",
						regex:       freeBaggageAllowanceRegex,
					},
					{
						id:          fastTrack,
						description: "Fast Track",
						length:      1,
						format:      `Y, N, or " "`,
						regex:       fastTrackRegex,
					},
				},
			},
			{
				id:          forIndividualAirlineUse,
				description: "For individual airline use",
				regex:       dotRegex,
			},
		},
	},
	{
		id:          beginningOfSecurityData,
		description: "Beginning of Security data",
		length:      1,
		format:      `"^"`,
		regex:       beginningOfSecurityDataRegex,
	},
	{
		id:          typeOfSecurityData,
		description: "Type of Security data",
		length:      1,
		format:      "an alphanumeric character",
		regex:       typeOfSecurityDataRegex,
	},
	{
		id:          lengthOfSecurityData,
		description: "Length of Security data",
		length:      2,
		format:      "a hex number",
		regex:       hexRegex,
		items: []item{
			{
				id:          securityData,
				description: "Security data",
				regex:       dotRegex,
			},
		},
	},
}
