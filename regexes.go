package bcbp

import "regexp"

const (
	formatCodeRegexString                             = "^[mM]$"
	numberOfLegsEncodedRegexString                    = "^[1-4]$"
	passengerNameRegexString                          = "^[a-zA-Z ]*/[a-zA-Z ]+$"
	electronicTicketRegexString                       = "^[eElL]$"
	operatingCarrierPNRCodeRegexString                = "^[a-zA-Z0-9]+ *$"
	airportCodeRegexString                            = "^[a-zA-Z]{3}$"
	operatingCarrierDesignatorRegexString             = "^[a-zA-Z0-9]{2,3} *$"
	flightNumberRegexString                           = "^[0-9]{4}[a-zA-Z ]{1}$"
	dateOfFlightRegexString                           = "^[0-2][0-9]{2}|3[0-5][0-9]|36[0-6]$"
	compartmentCodeRegexString                        = "^[a-aA-Z]$"
	seatNumberRegexString                             = "^[0-9]{3}[a-zA-Z]{1}$|^(?i)[INF]$|^(?i)[GATE]$|^(?i)[STBY]$"
	checkInSequenceNumberRegexString                  = "^[0-9]{4}[a-zA-Z ]{1}$"
	passengerStatusRegexString                        = "^[a-zA-Z0-9]$"
	hexRegexString                                    = "^[a-fA-f0-9]{2}$"
	beginningOfVersionNumberRegexString               = "^>$"
	versionNumberRegexString                          = "^[1-8]$"
	passengerDescriptionRegexString                   = "^[a-zA-Z0-9 ]$"
	sourceOfCheckInRegexString                        = "(?i)^[WKXRMOTVA ]$"
	sourceOfBoardingPassIssuanceRegexString           = "(?i)^[WKXRMOTV ]$"
	dateOfIssueOfBoardingPassRegexString              = "^[0-9][0-2][0-9]{2}$|^[0-9]3[0-5][0-9]$|^[0-9]36[0-6]$|^ {4}$"
	documentTypeRegexString                           = "^[bBiI]$"
	airlineDesignatorOfBoardingPassIssuerRegexString  = "^[a-zA-Z0-9]{2,3} *$|^ {3}$"
	baggageTagLicensePlateNumberRegexString           = "^[0-2]{1}[0-9]{12}$|^ {13}$"
	airlineNumericCodeRegexString                     = "^[0-9]{3}$|^ {3}$"
	documentFormSerialNumberRegexString               = "^0*[a-zA-Z0-9]*$|^ {10}$"
	selecteeIndicatorRegexString                      = "^[0-2]$|^ {1}$"
	internationalDocumentationVerificationRegexString = "^[0-2]$|^ {1}$"
	marketingCarrierDesignatorRegexString             = "^[a-zA-Z0-9]{2,3} *$|^ {3}$"
	frequentFlyerAirlineDesignatorRegexString         = "^[a-zA-Z0-9]{2,3} *$|^ {3}$"
	frequentFlyerNumberRegexString                    = "^[a-zA-Z0-9]+ *$|^ {16}$"
	idadIndicatorRegexString                          = "^[a-zA-Z0-9 ]$"
	freeBaggageAllowanceRegexString                   = "^[0-9]{2}[kKlL]|[0-9](?i)(PC)$|^ {3}$"
	fastTrackRegexString                              = "^[yYnN ]$"
	dotRegexString                                    = "^.*$"
	beginningOfSecurityDataRegexString                = "^[\\^]$"
	typeOfSecurityDataRegexString                     = "^[a-zA-Z0-9]$"
)

var (
	formatCodeRegex                             = regexp.MustCompile(formatCodeRegexString)
	numberOfLegsEncodedRegex                    = regexp.MustCompile(numberOfLegsEncodedRegexString)
	passengerNameRegex                          = regexp.MustCompile(passengerNameRegexString)
	electronicTicketRegex                       = regexp.MustCompile(electronicTicketRegexString)
	operatingCarrierPNRCodeRegex                = regexp.MustCompile(operatingCarrierPNRCodeRegexString)
	airportCodeRegex                            = regexp.MustCompile(airportCodeRegexString)
	operatingCarrierDesignatorRegex             = regexp.MustCompile(operatingCarrierDesignatorRegexString)
	flightNumberRegex                           = regexp.MustCompile(flightNumberRegexString)
	dateOfFlightRegex                           = regexp.MustCompile(dateOfFlightRegexString)
	compartmentCodeRegex                        = regexp.MustCompile(compartmentCodeRegexString)
	seatNumberRegex                             = regexp.MustCompile(seatNumberRegexString)
	checkInSequenceNumberRegex                  = regexp.MustCompile(checkInSequenceNumberRegexString)
	passengerStatusRegex                        = regexp.MustCompile(passengerStatusRegexString)
	hexRegex                                    = regexp.MustCompile(hexRegexString)
	beginningOfVersionNumberRegex               = regexp.MustCompile(beginningOfVersionNumberRegexString)
	versionNumberRegex                          = regexp.MustCompile(versionNumberRegexString)
	passengerDescriptionRegex                   = regexp.MustCompile(passengerDescriptionRegexString)
	sourceOfCheckInRegex                        = regexp.MustCompile(sourceOfCheckInRegexString)
	sourceOfBoardingPassIssuanceRegex           = regexp.MustCompile(sourceOfBoardingPassIssuanceRegexString)
	dateOfIssueOfBoardingPassRegex              = regexp.MustCompile(dateOfIssueOfBoardingPassRegexString)
	documentTypeRegex                           = regexp.MustCompile(documentTypeRegexString)
	airlineDesignatorOfBoardingPassIssuerRegex  = regexp.MustCompile(airlineDesignatorOfBoardingPassIssuerRegexString)
	baggageTagLicensePlateNumberRegex           = regexp.MustCompile(baggageTagLicensePlateNumberRegexString)
	airlineNumericCodeRegex                     = regexp.MustCompile(airlineNumericCodeRegexString)
	documentFormSerialNumberRegex               = regexp.MustCompile(documentFormSerialNumberRegexString)
	selecteeIndicatorRegex                      = regexp.MustCompile(selecteeIndicatorRegexString)
	internationalDocumentationVerificationRegex = regexp.MustCompile(internationalDocumentationVerificationRegexString)
	marketingCarrierDesignatorRegex             = regexp.MustCompile(marketingCarrierDesignatorRegexString)
	frequentFlyerAirlineDesignatorRegex         = regexp.MustCompile(frequentFlyerAirlineDesignatorRegexString)
	frequentFlyerNumberRegex                    = regexp.MustCompile(frequentFlyerNumberRegexString)
	idadIndicatorRegex                          = regexp.MustCompile(idadIndicatorRegexString)
	freeBaggageAllowanceRegex                   = regexp.MustCompile(freeBaggageAllowanceRegexString)
	fastTrackRegex                              = regexp.MustCompile(fastTrackRegexString)
	dotRegex                                    = regexp.MustCompile(dotRegexString)
	beginningOfSecurityDataRegex                = regexp.MustCompile(beginningOfSecurityDataRegexString)
	typeOfSecurityDataRegex                     = regexp.MustCompile(typeOfSecurityDataRegexString)
)
