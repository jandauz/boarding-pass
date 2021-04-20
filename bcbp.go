package bcbp

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// BCBP is a structured representation of an IATA 792 Bar Coded Boarding Pass.
type BCBP struct {
	// FormatCode is the format of the BCBP. M for multiple.
	FormatCode string `json:"format_code"`

	// NumberOfLegsEncoded is the number of flight segments encoded on
	// the barcode. It is at minimum 1; up to a maximum of 4.
	NumberOfLegsEncoded uint `json:"number_of_legs_encoded"`

	// PassengerName is the name of the passenger. It is encoded in the
	// following format:
	//   SURNAME/GIVEN_NAME
	// Up to 20 characters are encoded. If there is not enough space for the
	// given name then the surname is truncated at the 18th character
	// followed by a "/" and one alpha initial.
	//
	// The formatting is left justified with trailing whitespaces.
	PassengerName string `json:"passenger_name"`

	// ElectronicTicketIndicator is a flag that indicates whether or not
	// the boarding pass is issued against an electronic ticket. E or L.
	ElectronicTicketIndicator string `json:"electronic_ticket_indicator"`

	// Version number is the version of IATA 792 spec that is used to encode
	// the barcode. The latest version is 8.
	VersionNumber uint `json:"version_number,omitempty"`

	// PassengerDescription is the description of the passenger. It can be one
	// of the following values:
	//   0 - Adult
	//   1 - Male
	//   2 - Female
	//   3 - Child
	//   4 - Infant
	//   5 - No passenger (cabin baggage)
	//   6 - Adult traveling with infant
	//   7 - Unaccompanied minor
	//   X - Unspecified
	//   U - Undisclosed
	//
	//   Values 8-9 and A-T, V, W, Y, and Z are reserved for future industry use.
	PassengerDescription string `json:"passenger_description,omitempty"`

	// SourceOfCheckin is where the check-in was initiated. It can be one of
	// the following values:
	//   W - Web
	//   K - Airport kiosk
	//   R - Remote or off site kiosk
	//   M - Mobile device
	//   O - Airport agent
	//   T - Town agent
	//   V - Third party vendor
	//   A - Automated check-in
	SourceOfCheckIn string `json:"source_of_check_in,omitempty"`

	// SourceOfBoardingPassIssuance is where the boarding pass was issued.
	// It can be one of the following values:
	//   W - Web printed
	//   K - Airport kiosk printed
	//   X - Transfer kiosk printed
	//   R - Remote or off site kiosk printed
	//   M - Mobile device printed
	//   O - Airport agent printed
	//   T - Town agent printed
	//   V - Third party vendor printed
	SourceOfBoardingPassIssuance string `json:"source_of_boarding_pass_issuance,omitempty"`

	// DateOfIssueOfBoarding pass is the date the boarding pass was issued
	// include the last digit of the year in Julian Date.
	// For example, if the current date is January 1, 2021 the equivalent
	// in Julian Date would be 1001.
	//
	// See https://en.wikipedia.org/wiki/Julian_day for more information.
	DateOfIssueOfBoardingPass string `json:"date_of_issue_of_boarding_pass,omitempty"`

	// DocumentType is the type of travel document provided.
	// B for boarding pass; I for itinerary receipt.
	DocumentType string `json:"document_type,omitempty"`

	// AirlineDesignatorOfBoardingPassIssuer is the airline code of the airline
	// that issued the boarding pass.
	//
	// The formatting is left justified with trailing whitespaces.
	AirlineDesignatorOfBoardingPassIssuer string `json:"airline_designator_of_boarding_pass_issuer,omitempty"`

	// BaggageTagLicensePlateNumber represents the first consecutive series of
	// bag tag license plate number(s). This is a 13 character field encoded
	// in the following format:
	//       0: "0" for interline tag, "1" for fall-back tag, "2" for interline rush tag
	//     1-3: Carrier numeric code
	//     4-9: Carrier initial tag number with leading zeroes
	//   10-12: Number of consecutive bags (up to 999)
	BaggageTagLicensePlateNumber string `json:"baggage_tag_license_plate_number,omitempty"`

	// FirstNonConsecutiveBaggageTagLicensePlateNumber represents additional
	// bag tag license plate number(s) that are not consecutive with the
	// first series. The format is the same as BaggageTagLicensePlateNumber.
	FirstNonConsecutiveBaggageTagLicensePlateNumber string `json:"first_non_consecutive_baggage_tag_license_plate_number,omitempty"`

	// SecondNonConsecutiveBaggageTagLicensePlateNumber represents additional
	// bag tag license plate number(s) that are not consecutive with the
	// second series. The format is the same as
	// SecondNonConsecutiveBaggageTagLicensePlateNumber.
	SecondNonConsecutiveBaggageTagLicensePlateNumber string `json:"second_non_consecutive_baggage_tag_license_plate_number,omitempty"`

	// Legs represent individual flight segments. The number of legs is
	// defined by NumberOfLegsEncoded.
	Legs Legs `json:"legs,omitempty"`

	// TypeOfSecurityData is the type of security used on the barcode.
	TypeOfSecurityData string `json:"type_of_security_data,omitempty"`

	// SecurityData is used to verify that the boarding pass was not tampered.
	SecurityData string `json:"security_data,omitempty"`

	// data is the data encoded on a Bar Coded Boarding Pass.
	data string

	// dateBuf is used as a buffer when using time.AppendFormat() to convert
	// Julian dates into RFC3339 full-date formats (2006-01-02).
	//
	// See https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
	// for more information.
	dateBuf []byte

	// pos is the starting index of the character being processed in data.
	// This is used by whitespace() for pretty printing error reports.
	pos int
}

// Legs is an array of 4 Leg.
type Legs [4]Leg

// MarshalJSON implements the encoding.Marshaler interface.
// The output is an array of Leg that omits empty elements.
func (l Legs) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	for i := 0; i < len(l); i++ {
		if l[i] == (Leg{}) {
			continue
		}

		// Realistically should not return an error. According to the
		// documentation, json.Marshal will return an UnsupportedTypeError if
		// v is a channel, complex, or function value. Since Leg is a data
		// model for a flight segment it will only contain built-in types.
		b, _ := json.Marshal(l[i])
		sb.Write(b)
		sb.WriteString(",")
	}
	return []byte("[" + strings.TrimSuffix(sb.String(), ",") + "]"), nil
}

// Leg is a flight segment. The repeated fields of a Bar Coded Boarding Pass
// are captured in each Leg.
type Leg struct {
	// OperatingCarrierPNRCode is the Passenger Name Record used to identify
	// the booking in the reservation system of the operating carrier.
	//
	// The formatting is left justified with trailing whitespaces and up to
	// 7 alphanumeric characters.
	OperatingCarrierPNRCode string `json:"operating_carrier_pnr_code"`

	// FromCityAirportCode is the IATA code of the origin airport.
	//
	// The formatting is 3 alpha characters.
	FromCityAirportCode string `json:"from_city_airport_code"`

	// ToCityAirportCode is the IATA code of the destination airport.
	//
	// The formatting is 3 alpha characters.
	ToCityAirportCode string `json:"to_city_airport_code"`

	// OperatingCarrierDesignator is the airline code of the operating carrier.
	// It can be the same as MarketingCarrierDesignator. Two or three character
	// designators may be used.
	//
	// The formatting is left justified with trailing whitespaces.
	OperatingCarrierDesignator string `json:"operating_carrier_designator"`

	// FlightNumber is the number of the flight.
	//
	// The formatting is up to 4 digits with leading zeroes followed by an
	// optional alpha suffix or whitespace.
	FlightNumber string `json:"flight_number"`

	// DateOfFlight is the scheduled flight date in Julian Date. The date
	// is expressed in the number of days (inclusive) from January 1.
	// For example, if the current date is January 1, the Julian Date is 1.
	//
	// See https://en.wikipedia.org/wiki/Julian_day for more information.
	//
	// The formatting is numerical with leading zeroes.
	DateOfFlight string `json:"date_of_flight"`

	// CompartmentCode is the code of the compartment also know as the
	// Cabin Type.
	//
	// It can be one of the following values:
	//   First Class Category
	//   R - Supersonic
	//   P - First Class Premium
	//   F - First Class
	//   A - First Class Discounted
	//
	//   Business Class Category
	//   J - Business Class Premium
	//   C - Business Class
	//   D - Business Class Discounted
	//   I - Business Class Discounted
	//   Z - Business Class Discounted
	//
	//   Economy/Coach Class Category
	//   W - Economy/Coach Premium
	//   S - Economy/Coach
	//   Y - Economy/Coach
	//   B - Economy/Coach Discounted
	//   H - Economy/Coach Discounted
	//   K - Economy/Coach Discounted
	//   L - Economy/Coach Discounted
	//   M - Economy/Coach Discounted
	//   N - Economy/Coach Discounted
	//   Q - Economy/Coach Discounted
	//   T - Economy/Coach Discounted
	//   V - Economy/Coach Discounted
	//   X - Economy/Coach Discounted
	CompartmentCode string `json:"compartment_code"`

	// SeatNumber is the seat assigned to the passenger.
	//
	// The formatting is 3 digits with leading zeroes followed by an alpha.
	// There are exceptions where the 4 characters can be:
	//   INF
	//   GATE
	//   STBY
	SeatNumber string `json:"seat_number"`

	// CheckInSequenceNumber is the order in which the passenger has checked-in
	// for the flight.
	//
	// The formatting is 4 digits with leading zeroes followed by an optional
	// alpha or whitespace. Infant passengers are an exception in which case
	// 5 alphanumeric characters may be used.
	CheckInSequenceNumber string `json:"check_in_sequence_number"`

	// PassengerStatus is the status of the passenger.
	//
	// It can be one of the following values:
	//   0 - Ticket issuance/Passenger not checked in
	//   1 - Ticket issuance/Passenger checked in
	//   2 - Baggage checked/Passenger not checked in
	//   3 - Baggage checked/Passenger checked in
	//   4 - Passenger passed security check
	//   5 - Passenger passed gate exit (coupon used)
	//   6 - Transit
	//   7 - Standby
	//       Seat number not printed on boarding pass at time of check-in
	//       Seat number to be printed at time of seat assignment
	//   8 - Boarding data revalidation done
	//       Gate, Boarding time, and Seat on Revalidation Field already used
	//   9 - Original boarding line used at time of ticket issuance
	//   A - Up or down-grading required at close out
	//       e.g. when passenger waitlisted in C class and OK in Y class
	//
	// Values B-Z are reserved for future industry use
	PassengerStatus string `json:"passenger_status"`

	// AirlineNumericCode is the numeric code of the airline.
	//
	// See here - https://www.iata.org/en/publications/directories/code-search/.
	//
	// The formatting is right justified 3 digits with leading zeroes.
	AirlineNumericCode string `json:"airline_numeric_code,omitempty"`

	// DocumentFormSerialNumber is the document number which is comprised of:
	//   Airline code
	//   Form code
	//   Serial number
	//
	// The formatting is right justified 10 alphanumeric characters with leading zeroes.
	DocumentFormSerialNumber string `json:"document_form_serial_number,omitempty"`

	// SelecteeIndicator is a flag that is used by some agencies for additional
	// screening and it assists airlines to classify customers that require
	// inspection at airports in certain countries. It is a mandatory field
	// when US travel is involved.
	//
	// It can be one of the following values:
	//   0 - Not selectee
	//   1 - Selectee
	//   2 - Known passenger
	SelecteeIndicator string `json:"selectee_indicator,omitempty"`

	// InternationalDocumentationVerification is a flag that is used carriers
	// to identify passengers requiring their travel documentation to be
	// verified.
	//
	// It can be one of the following values:
	//   0 - Travel documentation verification required
	//   1 - Travel documentation verification not required
	//   2 - Travel documentation verification performed
	InternationalDocumentationVerification string `json:"international_documentation_verification,omitempty"`

	// MarketingCarrierDesignator is the airline code of the marketing carrier.
	// It can be the same as OperatingCarrierDesignator. Two or three character
	// designators may be used.
	//
	// The formatting is left justified with trailing spaces.
	MarketingCarrierDesignator string `json:"marketing_carrier_designator,omitempty"`

	// FrequentFlyerAirlineDesignator is the airline code of the airline's
	// frequent flyer program. Two or three character designators may be used.
	//
	// The formatting is left justified with trailing spaces.
	FrequentFlyerAirlineDesignator string `json:"frequent_flyer_airline_designator,omitempty"`

	// FrequentFlyerNumber is the passenger's number in the airline's frequent
	// flyer program.
	//
	// The formatting depends on individual carriers and alliances.
	// It is left justified with trailing whitespaces. It can be between 13-16
	// characters and either numeric or alphanumeric characters.
	FrequentFlyerNumber string `json:"frequent_flyer_number,omitempty"`

	// IDADIndicator is a flag that specifics an industry discount ticket or
	// agency discount codes.
	//
	// It can be one of the following values:
	//   0 - IDN1 positive space
	//   1 - IDN2 space available
	//   2 - IDB1 positive space
	//   3 - IDB2 space available
	//   4 - AD
	//   5 - DG
	//   6 - DM
	//   7 - GE
	//   8 - IG
	//   9 - RG
	//   A - UD
	//   B - ID
	//
	// The following are industry discounts not followed by any classification
	//   C - IDFS1
	//   D - IDFS2
	//   E - IDR1
	//
	// Values G-Z are reserved for future industry use.
	IDADIndicator string `json:"idad_indicator,omitempty"`

	// FreeBaggageAllowance specifies the weight, either in K (kilos) or
	// L (pounds), or PC (number of pieces).
	//
	// For example, it can be 20K, 40L, or 2PC.
	FreeBaggageAllowance string `json:"free_baggage_allowance,omitempty"`

	// FastTrack is a flag that specifies if the passenger is entitled to use
	// a priority, security, or immigration lane.
	//
	// It can be one of the following values:
	//   Y - Yes
	//   N - No
	//
	// A whitespace means unqualified.
	FastTrack string `json:"fast_track,omitempty"`

	// ForIndividualAirlineUse is a special field that airlines may use to
	// populate with different entries such as but not limited to:
	//   Frequent flyer tier
	//   Passenger preferences
	ForIndividualAirlineUse string `json:"for_individual_airline_use,omitempty"`
}

// FromStr creates a new BCBP from s.
func FromStr(s string) (BCBP, error) {
	if len(s) < 60 {
		return BCBP{}, InsufficientData(s, len(s))
	}

	if pos, ok := ascii(s); !ok {
		val, _ := utf8.DecodeRuneInString(s[pos:])
		return BCBP{}, NonASCII(s, pos+1, val)
	}

	if s[0:1] != "M" {
		return BCBP{}, UnsupportedBoardingPass(s, s[0:1])
	}

	return fromStr(s)
}

// ascii checks s to determine if it contains only ASCII characters.
// If a unicode character is found then the index of the rune is
// returned.
func ascii(s string) (int, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return i, false
		}
	}
	return 0, true
}

func fromStr(s string) (BCBP, error) {
	if !spec[numberOfLegsEncoded].validate(s[1:2]) {
		return BCBP{},
			InvalidDataFormat(s, 2, spec[numberOfLegsEncoded], s[1:2])
	}

	// No need to check error as data validation happens above
	legs, _ := strconv.Atoi(s[1:2])

	// Dates use RFC-3339 full-date format. These are 10 bytes long.
	// Allocate a 16 byte array, create a slice, and assign to dateBuf.
	buf := [16]byte{}
	b := BCBP{
		data:                s,
		NumberOfLegsEncoded: uint(legs),
		dateBuf:             buf[:0],
		pos:                 1,
	}

	// Iterate over the number of legs specified and recursively process the
	// items defined in spec.
	for leg := 0; leg < legs; leg++ {
		b.Legs[leg] = Leg{}

		for _, item := range spec {
			switch item.id {
			// Security items are unique and are at the end of the Bar Coded
			// Boarding Pass. Set those fields last.
			case beginningOfSecurityData,
				typeOfSecurityData,
				lengthOfSecurityData,
				securityData:
				continue
			}

			processed, err := b.setFieldByItem(s, item, leg)
			if err != nil {
				return b, err
			}
			s = s[processed:]
		}
	}

	// If len of s is 0 then there is nothing more to process.
	if len(s) == 0 {
		return b, nil
	}

	// Otherwise, there is more to process. However, if the prefix isn't the "^"
	// character, which marks the beginning of security section, then return
	// ErrProcessItemFailed.
	if s[0:1] != "^" {
		return b, InvalidDataFormat(b.data, b.pos, spec[fieldSizeOfVariableSizeField+1], s[0:1])
	}

	// Security items start after fieldSizeOfVariableSizeField in spec.
	for _, item := range spec[fieldSizeOfVariableSizeField+1:] {
		processed, err := b.setFieldByItem(s, item, 0)
		if err != nil {
			return b, err
		}
		s = s[processed:]
	}

	// At this point decoding should be successfully completed and s should
	// be empty. If not, then that means the barcode data has extra unprocessed
	// characters.
	if s != "" {
		return b, UnknownData(b.data, b.pos, s)
	}
	return b, nil
}

// setFieldByItem sets the value of the appropriate field based on the item ID.
// The length of the field is returned if successfully set. Otherwise, 0 and
// and error is returned.
func (b *BCBP) setFieldByItem(s string, item item, leg int) (int, error) {
	// Unique items appear only once in a Bar Coded Boarding Pass.
	// Return immediately if the item is unique and current leg is
	// greater than 0.
	switch item.id {
	case formatCode,
		numberOfLegsEncoded,
		passengerName,
		electronicTicketIndicator,
		beginningOfVersionNumber,
		versionNumber,
		fieldSizeOfFollowingStructuredMessageUnique,
		passengerDescription,
		sourceOfCheckin,
		sourceOfBoardingPassIssuance,
		dateOfIssueOfBoardingPass,
		documentType,
		airlineDesignatorOfBoardingPassIssuer,
		baggageTagLicensePlateNumber,
		firstNonConsecutiveBaggageTagLicensePlateNumber,
		secondNonConsecutiveBaggageTagLicensePlateNumber:
		if leg > 0 {
			return 0, nil
		}
	}

	itemLen := item.length
	// forIndividualAirlineUse and lengthOfSecurityData do not have a static
	// length. It's length is the remainder of the conditional section of the
	// Bar Coded Boarding Pass.
	switch item.id {
	case forIndividualAirlineUse,
		securityData:
		itemLen = len(s)
	}

	// Validate that the data matches the item's format.
	if !item.validate(s[:itemLen]) {
		return 0, InvalidDataFormat(b.data, b.pos, item, s[:itemLen])
	}

	// Substring the value and assign to the appropriate BCBP field based on
	// item.id.
	val := strings.TrimSpace(s[:itemLen])
	switch item.id {
	case formatCode:
		b.FormatCode = val
	case passengerName:
		b.PassengerName = val
	case electronicTicketIndicator:
		b.ElectronicTicketIndicator = val
	case operatingCarrierPNRCode:
		b.Legs[leg].OperatingCarrierPNRCode = val
	case fromCityAirportCode:
		b.Legs[leg].FromCityAirportCode = val
	case toCityAirportCode:
		b.Legs[leg].ToCityAirportCode = val
	case operatingCarrierDesignator:
		b.Legs[leg].OperatingCarrierDesignator = val
	case flightNumber:
		b.Legs[leg].FlightNumber = val
	case dateOfFlight:
		// Re-slice dateBuf so that we append the date format at the start
		// of the buffer instead of at the end.
		b.dateBuf = b.dateBuf[:0]

		// item.validate() ensures val is a number, no need to check error
		d, _ := strconv.Atoi(val)
		t := time.Date(time.Now().Year(), time.January, 0, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(0, 0, d)
		b.dateBuf = t.AppendFormat(b.dateBuf, "2006-01-02")

		// See https://github.com/golang/go/issues/25484#issuecomment-391415660.
		// This copies strings.Builder.String() way of copying byte array to string.
		b.Legs[leg].DateOfFlight = *(*string)(unsafe.Pointer(&b.dateBuf))
	case compartmentCode:
		b.Legs[leg].CompartmentCode = val
	case seatNumber:
		b.Legs[leg].SeatNumber = val
	case checkinSequenceNumber:
		b.Legs[leg].CheckInSequenceNumber = val
	case passengerStatus:
		b.Legs[leg].PassengerStatus = val
	case versionNumber:
		// item.validate() ensures val is a number, no need to check error
		n, _ := strconv.Atoi(val)
		b.VersionNumber = uint(n)
	case passengerDescription:
		b.PassengerDescription = val
	case sourceOfCheckin:
		b.SourceOfCheckIn = val
	case sourceOfBoardingPassIssuance:
		b.SourceOfBoardingPassIssuance = val
	case dateOfIssueOfBoardingPass:
		// Re-slice dateBuf so that we append the date format at the start
		// of the buffer instead of at the end.
		b.dateBuf = b.dateBuf[:0]

		// item.validate() ensures val is a number, no need to check error
		y, _ := strconv.Atoi(val[:1])
		n := time.Now().Year() % 10
		y -= n

		// item.validate() ensures val is a number
		d, _ := strconv.Atoi(val[1:])
		t := time.Date(time.Now().Year(), time.January, 0, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(y, 0, d)
		b.dateBuf = t.AppendFormat(b.dateBuf, "2006-01-02")

		// See https://github.com/golang/go/issues/25484#issuecomment-391415660.
		// This copies strings.Builder.String() way of copying byte array to string.
		b.DateOfIssueOfBoardingPass = *(*string)(unsafe.Pointer(&b.dateBuf))
	case documentType:
		b.DocumentType = val
	case airlineDesignatorOfBoardingPassIssuer:
		b.AirlineDesignatorOfBoardingPassIssuer = val
	case baggageTagLicensePlateNumber:
		b.BaggageTagLicensePlateNumber = val
	case firstNonConsecutiveBaggageTagLicensePlateNumber:
		b.FirstNonConsecutiveBaggageTagLicensePlateNumber = val
	case secondNonConsecutiveBaggageTagLicensePlateNumber:
		b.SecondNonConsecutiveBaggageTagLicensePlateNumber = val
	case airlineNumericCode:
		b.Legs[leg].AirlineNumericCode = val
	case documentFormSerialNumber:
		b.Legs[leg].DocumentFormSerialNumber = val
	case selecteeIndicator:
		b.Legs[leg].SelecteeIndicator = val
	case internationalDocumentationVerification:
		b.Legs[leg].InternationalDocumentationVerification = val
	case marketingCarrierDesignator:
		b.Legs[leg].MarketingCarrierDesignator = val
	case frequentFlyerAirlineDesignator:
		b.Legs[leg].FrequentFlyerAirlineDesignator = val
	case frequentFlyerNumber:
		b.Legs[leg].FrequentFlyerNumber = val
	case idadIndicator:
		b.Legs[leg].IDADIndicator = val
	case freeBaggageAllowance:
		b.Legs[leg].FreeBaggageAllowance = val
	case fastTrack:
		b.Legs[leg].FastTrack = val
	case forIndividualAirlineUse:
		b.Legs[leg].ForIndividualAirlineUse = val
	case typeOfSecurityData:
		b.TypeOfSecurityData = val
	case securityData:
		b.SecurityData = val
	}

	// Reassign s to the remaining unprocessed characters.
	s = s[itemLen:]

	if item.items == nil {
		// Set the position of the next character to be processed.
		b.pos += itemLen
		return itemLen, nil
	}

	// If item.items does not equal nil then there is a sub-section to process.
	//
	// The beginning of sub-sections are indicated by the following items:
	//   fieldSizeOfVariableSizeField
	//   fieldSizeOfFollowingStructuredMessageUnique
	//   fieldSizeOfFollowingStructuredMessageRepeated
	//   lengthOfSecurityData
	//
	// If the current item is neither of these then ErrMalformedSpec is
	// returned. Otherwise, convert val from a hex string to int and slice
	// s up to the length of the section.
	if item.id != fieldSizeOfVariableSizeField &&
		item.id != fieldSizeOfFollowingStructuredMessageUnique &&
		item.id != fieldSizeOfFollowingStructuredMessageRepeated &&
		item.id != lengthOfSecurityData {
		return itemLen, MalformedSpec(b.data, b.pos, item)
	}
	sectionLen, err := strconv.ParseInt(val, 16, 32)
	if err != nil {
		return itemLen, InvalidDataFormat(b.data, b.pos, item, val)
	}

	// If the sub-section length is greater than the length of s then
	// the Bar Coded Boarding Pass is malformed and is missing data.
	if int(sectionLen) > len(s) {
		return itemLen, UnexpectedEndOfInput(b.data, b.pos+item.length, item, s, int(sectionLen))
	}

	// Substring s based on the length of the sub-section.
	sectionStr := s[:sectionLen]
	// Set the position of the next character to be processed.
	b.pos += itemLen
	for _, subItem := range item.items {
		// If sectionStr is empty then processing of the sub-section is
		// complete. No need to continue processing.
		if sectionStr == "" {
			break
		}

		subItemLen, err := b.setFieldByItem(sectionStr, subItem, leg)
		if err != nil {
			return subItemLen, err
		}

		// Add the sub-item length to itemLen since we are recursively
		// processing the sub-section and thus treating it as one field
		// being processed.
		itemLen += subItemLen

		// Reassign sectionStr to the remaining unprocessed characters.
		sectionStr = sectionStr[subItemLen:]
	}

	return itemLen, nil
}
