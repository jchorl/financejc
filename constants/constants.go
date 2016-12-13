package constants

import (
	"errors"
)

const (
	CTX_DB      = "database"
	CTX_ES      = "elasticsearch"
	CTX_USER_ID = "user"

	ES_INDEX = "financejc"

	FIXED_INTERVAL  = "fixedInterval"
	FIXED_DAY_WEEK  = "fixedDayWeek"
	FIXED_DAY_MONTH = "fixedDayMonth"
	FIXED_DAY_YEAR  = "fixedDayYear"

	ADMIN_UID = uint(1)
)

var CTX_KEYS = [...]string{
	CTX_DB,
	CTX_ES,
	CTX_USER_ID,
}

var (
	Forbidden       = errors.New("User does not have permission to access this resource.")
	NotLoggedIn     = errors.New("User is not logged in.")
	BadRequest      = errors.New("Request contains malformed data.")
	InvalidCurrency = errors.New("The specified currency is not recognized.")
)

type currency struct {
	Name               string `json:"name"`
	Code               string `json:"code"`
	DigitsAfterDecimal int    `json:"digitsAfterDecimal"`
}

var CurrencyInfo = map[string]currency{
	"SDG": currency{
		Name:               "Sudanese Pound",
		Code:               "SDG",
		DigitsAfterDecimal: 2,
	},
	"PAB": currency{
		Name:               "Balboa",
		Code:               "PAB",
		DigitsAfterDecimal: 2,
	},
	"SAR": currency{
		Name:               "Saudi Riyal",
		Code:               "SAR",
		DigitsAfterDecimal: 2,
	},
	"PEN": currency{
		Name:               "Sol",
		Code:               "PEN",
		DigitsAfterDecimal: 2,
	},
	"SRD": currency{
		Name:               "Surinam Dollar",
		Code:               "SRD",
		DigitsAfterDecimal: 2,
	},
	"CAD": currency{
		Name:               "Canadian Dollar",
		Code:               "CAD",
		DigitsAfterDecimal: 2,
	},
	"COP": currency{
		Name:               "Colombian Peso",
		Code:               "COP",
		DigitsAfterDecimal: 2,
	},
	"MYR": currency{
		Name:               "Malaysian Ringgit",
		Code:               "MYR",
		DigitsAfterDecimal: 2,
	},
	"SVC": currency{
		Name:               "El Salvador Colon",
		Code:               "SVC",
		DigitsAfterDecimal: 2,
	},
	"UYI": currency{
		Name:               "Uruguay Peso en Unidades Indexadas (URUIURUI)",
		Code:               "UYI",
		DigitsAfterDecimal: 0,
	},
	"UYU": currency{
		Name:               "Peso Uruguayo",
		Code:               "UYU",
		DigitsAfterDecimal: 2,
	},
	"TZS": currency{
		Name:               "Tanzanian Shilling",
		Code:               "TZS",
		DigitsAfterDecimal: 2,
	},
	"LAK": currency{
		Name:               "Kip",
		Code:               "LAK",
		DigitsAfterDecimal: 2,
	},
	"PHP": currency{
		Name:               "Philippine Peso",
		Code:               "PHP",
		DigitsAfterDecimal: 2,
	},
	"BYN": currency{
		Name:               "Belarusian Ruble",
		Code:               "BYN",
		DigitsAfterDecimal: 2,
	},
	"BBD": currency{
		Name:               "Barbados Dollar",
		Code:               "BBD",
		DigitsAfterDecimal: 2,
	},
	"KRW": currency{
		Name:               "Won",
		Code:               "KRW",
		DigitsAfterDecimal: 0,
	},
	"BRL": currency{
		Name:               "Brazilian Real",
		Code:               "BRL",
		DigitsAfterDecimal: 2,
	},
	"COU": currency{
		Name:               "Unidad de Valor Real",
		Code:               "COU",
		DigitsAfterDecimal: 2,
	},
	"EUR": currency{
		Name:               "Euro",
		Code:               "EUR",
		DigitsAfterDecimal: 2,
	},
	"ALL": currency{
		Name:               "Lek",
		Code:               "ALL",
		DigitsAfterDecimal: 2,
	},
	"DJF": currency{
		Name:               "Djibouti Franc",
		Code:               "DJF",
		DigitsAfterDecimal: 0,
	},
	"SEK": currency{
		Name:               "Swedish Krona",
		Code:               "SEK",
		DigitsAfterDecimal: 2,
	},
	"MAD": currency{
		Name:               "Moroccan Dirham",
		Code:               "MAD",
		DigitsAfterDecimal: 2,
	},
	"PKR": currency{
		Name:               "Pakistan Rupee",
		Code:               "PKR",
		DigitsAfterDecimal: 2,
	},
	"XOF": currency{
		Name:               "CFA Franc BCEAO",
		Code:               "XOF",
		DigitsAfterDecimal: 0,
	},
	"RON": currency{
		Name:               "Romanian Leu",
		Code:               "RON",
		DigitsAfterDecimal: 2,
	},
	"PGK": currency{
		Name:               "Kina",
		Code:               "PGK",
		DigitsAfterDecimal: 2,
	},
	"CHW": currency{
		Name:               "WIR Franc",
		Code:               "CHW",
		DigitsAfterDecimal: 2,
	},
	"GYD": currency{
		Name:               "Guyana Dollar",
		Code:               "GYD",
		DigitsAfterDecimal: 2,
	},
	"KPW": currency{
		Name:               "North Korean Won",
		Code:               "KPW",
		DigitsAfterDecimal: 2,
	},
	"HRK": currency{
		Name:               "Kuna",
		Code:               "HRK",
		DigitsAfterDecimal: 2,
	},
	"USN": currency{
		Name:               "US Dollar (Next day)",
		Code:               "USN",
		DigitsAfterDecimal: 2,
	},
	"AFN": currency{
		Name:               "Afghani",
		Code:               "AFN",
		DigitsAfterDecimal: 2,
	},
	"CUC": currency{
		Name:               "Peso Convertible",
		Code:               "CUC",
		DigitsAfterDecimal: 2,
	},
	"KWD": currency{
		Name:               "Kuwaiti Dinar",
		Code:               "KWD",
		DigitsAfterDecimal: 3,
	},
	"USD": currency{
		Name:               "US Dollar",
		Code:               "USD",
		DigitsAfterDecimal: 2,
	},
	"LRD": currency{
		Name:               "Liberian Dollar",
		Code:               "LRD",
		DigitsAfterDecimal: 2,
	},
	"NAD": currency{
		Name:               "Namibia Dollar",
		Code:               "NAD",
		DigitsAfterDecimal: 2,
	},
	"QAR": currency{
		Name:               "Qatari Rial",
		Code:               "QAR",
		DigitsAfterDecimal: 2,
	},
	"CZK": currency{
		Name:               "Czech Koruna",
		Code:               "CZK",
		DigitsAfterDecimal: 2,
	},
	"LKR": currency{
		Name:               "Sri Lanka Rupee",
		Code:               "LKR",
		DigitsAfterDecimal: 2,
	},
	"MOP": currency{
		Name:               "Pataca",
		Code:               "MOP",
		DigitsAfterDecimal: 2,
	},
	"WST": currency{
		Name:               "Tala",
		Code:               "WST",
		DigitsAfterDecimal: 2,
	},
	"LSL": currency{
		Name:               "Loti",
		Code:               "LSL",
		DigitsAfterDecimal: 2,
	},
	"ZWL": currency{
		Name:               "Zimbabwe Dollar",
		Code:               "ZWL",
		DigitsAfterDecimal: 2,
	},
	"AED": currency{
		Name:               "UAE Dirham",
		Code:               "AED",
		DigitsAfterDecimal: 2,
	},
	"GHS": currency{
		Name:               "Ghana Cedi",
		Code:               "GHS",
		DigitsAfterDecimal: 2,
	},
	"BMD": currency{
		Name:               "Bermudian Dollar",
		Code:               "BMD",
		DigitsAfterDecimal: 2,
	},
	"SLL": currency{
		Name:               "Leone",
		Code:               "SLL",
		DigitsAfterDecimal: 2,
	},
	"RSD": currency{
		Name:               "Serbian Dinar",
		Code:               "RSD",
		DigitsAfterDecimal: 2,
	},
	"SCR": currency{
		Name:               "Seychelles Rupee",
		Code:               "SCR",
		DigitsAfterDecimal: 2,
	},
	"INR": currency{
		Name:               "Indian Rupee",
		Code:               "INR",
		DigitsAfterDecimal: 2,
	},
	"NIO": currency{
		Name:               "Cordoba Oro",
		Code:               "NIO",
		DigitsAfterDecimal: 2,
	},
	"MDL": currency{
		Name:               "Moldovan Leu",
		Code:               "MDL",
		DigitsAfterDecimal: 2,
	},
	"XPF": currency{
		Name:               "CFP Franc",
		Code:               "XPF",
		DigitsAfterDecimal: 0,
	},
	"CDF": currency{
		Name:               "Congolese Franc",
		Code:               "CDF",
		DigitsAfterDecimal: 2,
	},
	"MNT": currency{
		Name:               "Tugrik",
		Code:               "MNT",
		DigitsAfterDecimal: 2,
	},
	"MXV": currency{
		Name:               "Mexican Unidad de Inversion (UDI)",
		Code:               "MXV",
		DigitsAfterDecimal: 2,
	},
	"IDR": currency{
		Name:               "Rupiah",
		Code:               "IDR",
		DigitsAfterDecimal: 2,
	},
	"SHP": currency{
		Name:               "Saint Helena Pound",
		Code:               "SHP",
		DigitsAfterDecimal: 2,
	},
	"BTN": currency{
		Name:               "Ngultrum",
		Code:               "BTN",
		DigitsAfterDecimal: 2,
	},
	"NGN": currency{
		Name:               "Naira",
		Code:               "NGN",
		DigitsAfterDecimal: 2,
	},
	"BOV": currency{
		Name:               "Mvdol",
		Code:               "BOV",
		DigitsAfterDecimal: 2,
	},
	"BHD": currency{
		Name:               "Bahraini Dinar",
		Code:               "BHD",
		DigitsAfterDecimal: 3,
	},
	"BYR": currency{
		Name:               "Belarusian Ruble",
		Code:               "BYR",
		DigitsAfterDecimal: 0,
	},
	"KZT": currency{
		Name:               "Tenge",
		Code:               "KZT",
		DigitsAfterDecimal: 2,
	},
	"CLP": currency{
		Name:               "Chilean Peso",
		Code:               "CLP",
		DigitsAfterDecimal: 0,
	},
	"ERN": currency{
		Name:               "Nakfa",
		Code:               "ERN",
		DigitsAfterDecimal: 2,
	},
	"LBP": currency{
		Name:               "Lebanese Pound",
		Code:               "LBP",
		DigitsAfterDecimal: 2,
	},
	"SSP": currency{
		Name:               "South Sudanese Pound",
		Code:               "SSP",
		DigitsAfterDecimal: 2,
	},
	"SOS": currency{
		Name:               "Somali Shilling",
		Code:               "SOS",
		DigitsAfterDecimal: 2,
	},
	"AUD": currency{
		Name:               "Australian Dollar",
		Code:               "AUD",
		DigitsAfterDecimal: 2,
	},
	"TJS": currency{
		Name:               "Somoni",
		Code:               "TJS",
		DigitsAfterDecimal: 2,
	},
	"CHE": currency{
		Name:               "WIR Euro",
		Code:               "CHE",
		DigitsAfterDecimal: 2,
	},
	"MXN": currency{
		Name:               "Mexican Peso",
		Code:               "MXN",
		DigitsAfterDecimal: 2,
	},
	"ISK": currency{
		Name:               "Iceland Krona",
		Code:               "ISK",
		DigitsAfterDecimal: 0,
	},
	"GTQ": currency{
		Name:               "Quetzal",
		Code:               "GTQ",
		DigitsAfterDecimal: 2,
	},
	"BOB": currency{
		Name:               "Boliviano",
		Code:               "BOB",
		DigitsAfterDecimal: 2,
	},
	"CUP": currency{
		Name:               "Cuban Peso",
		Code:               "CUP",
		DigitsAfterDecimal: 2,
	},
	"NPR": currency{
		Name:               "Nepalese Rupee",
		Code:               "NPR",
		DigitsAfterDecimal: 2,
	},
	"SYP": currency{
		Name:               "Syrian Pound",
		Code:               "SYP",
		DigitsAfterDecimal: 2,
	},
	"UZS": currency{
		Name:               "Uzbekistan Sum",
		Code:               "UZS",
		DigitsAfterDecimal: 2,
	},
	"ARS": currency{
		Name:               "Argentine Peso",
		Code:               "ARS",
		DigitsAfterDecimal: 2,
	},
	"UAH": currency{
		Name:               "Hryvnia",
		Code:               "UAH",
		DigitsAfterDecimal: 2,
	},
	"CRC": currency{
		Name:               "Costa Rican Colon",
		Code:               "CRC",
		DigitsAfterDecimal: 2,
	},
	"IQD": currency{
		Name:               "Iraqi Dinar",
		Code:               "IQD",
		DigitsAfterDecimal: 3,
	},
	"KHR": currency{
		Name:               "Riel",
		Code:               "KHR",
		DigitsAfterDecimal: 2,
	},
	"BZD": currency{
		Name:               "Belize Dollar",
		Code:               "BZD",
		DigitsAfterDecimal: 2,
	},
	"SZL": currency{
		Name:               "Lilangeni",
		Code:               "SZL",
		DigitsAfterDecimal: 2,
	},
	"FJD": currency{
		Name:               "Fiji Dollar",
		Code:               "FJD",
		DigitsAfterDecimal: 2,
	},
	"MVR": currency{
		Name:               "Rufiyaa",
		Code:               "MVR",
		DigitsAfterDecimal: 2,
	},
	"HKD": currency{
		Name:               "Hong Kong Dollar",
		Code:               "HKD",
		DigitsAfterDecimal: 2,
	},
	"BND": currency{
		Name:               "Brunei Dollar",
		Code:               "BND",
		DigitsAfterDecimal: 2,
	},
	"TRY": currency{
		Name:               "Turkish Lira",
		Code:               "TRY",
		DigitsAfterDecimal: 2,
	},
	"CVE": currency{
		Name:               "Cabo Verde Escudo",
		Code:               "CVE",
		DigitsAfterDecimal: 2,
	},
	"MUR": currency{
		Name:               "Mauritius Rupee",
		Code:               "MUR",
		DigitsAfterDecimal: 2,
	},
	"YER": currency{
		Name:               "Yemeni Rial",
		Code:               "YER",
		DigitsAfterDecimal: 2,
	},
	"BWP": currency{
		Name:               "Pula",
		Code:               "BWP",
		DigitsAfterDecimal: 2,
	},
	"SBD": currency{
		Name:               "Solomon Islands Dollar",
		Code:               "SBD",
		DigitsAfterDecimal: 2,
	},
	"GBP": currency{
		Name:               "Pound Sterling",
		Code:               "GBP",
		DigitsAfterDecimal: 2,
	},
	"DZD": currency{
		Name:               "Algerian Dinar",
		Code:               "DZD",
		DigitsAfterDecimal: 2,
	},
	"KYD": currency{
		Name:               "Cayman Islands Dollar",
		Code:               "KYD",
		DigitsAfterDecimal: 2,
	},
	"XCD": currency{
		Name:               "East Caribbean Dollar",
		Code:               "XCD",
		DigitsAfterDecimal: 2,
	},
	"VEF": currency{
		Name:               "Bolívar",
		Code:               "VEF",
		DigitsAfterDecimal: 2,
	},
	"TMT": currency{
		Name:               "Turkmenistan New Manat",
		Code:               "TMT",
		DigitsAfterDecimal: 2,
	},
	"MRO": currency{
		Name:               "Ouguiya",
		Code:               "MRO",
		DigitsAfterDecimal: 2,
	},
	"ANG": currency{
		Name:               "Netherlands Antillean Guilder",
		Code:               "ANG",
		DigitsAfterDecimal: 2,
	},
	"GNF": currency{
		Name:               "Guinea Franc",
		Code:               "GNF",
		DigitsAfterDecimal: 0,
	},
	"JMD": currency{
		Name:               "Jamaican Dollar",
		Code:               "JMD",
		DigitsAfterDecimal: 2,
	},
	"KGS": currency{
		Name:               "Som",
		Code:               "KGS",
		DigitsAfterDecimal: 2,
	},
	"SGD": currency{
		Name:               "Singapore Dollar",
		Code:               "SGD",
		DigitsAfterDecimal: 2,
	},
	"DKK": currency{
		Name:               "Danish Krone",
		Code:               "DKK",
		DigitsAfterDecimal: 2,
	},
	"KMF": currency{
		Name:               "Comoro Franc",
		Code:               "KMF",
		DigitsAfterDecimal: 0,
	},
	"ETB": currency{
		Name:               "Ethiopian Birr",
		Code:               "ETB",
		DigitsAfterDecimal: 2,
	},
	"TOP": currency{
		Name:               "Pa’anga",
		Code:               "TOP",
		DigitsAfterDecimal: 2,
	},
	"AMD": currency{
		Name:               "Armenian Dram",
		Code:               "AMD",
		DigitsAfterDecimal: 2,
	},
	"EGP": currency{
		Name:               "Egyptian Pound",
		Code:               "EGP",
		DigitsAfterDecimal: 2,
	},
	"BIF": currency{
		Name:               "Burundi Franc",
		Code:               "BIF",
		DigitsAfterDecimal: 0,
	},
	"BGN": currency{
		Name:               "Bulgarian Lev",
		Code:               "BGN",
		DigitsAfterDecimal: 2,
	},
	"FKP": currency{
		Name:               "Falkland Islands Pound",
		Code:               "FKP",
		DigitsAfterDecimal: 2,
	},
	"LYD": currency{
		Name:               "Libyan Dinar",
		Code:               "LYD",
		DigitsAfterDecimal: 3,
	},
	"RUB": currency{
		Name:               "Russian Ruble",
		Code:               "RUB",
		DigitsAfterDecimal: 2,
	},
	"TWD": currency{
		Name:               "New Taiwan Dollar",
		Code:               "TWD",
		DigitsAfterDecimal: 2,
	},
	"HTG": currency{
		Name:               "Gourde",
		Code:               "HTG",
		DigitsAfterDecimal: 2,
	},
	"HNL": currency{
		Name:               "Lempira",
		Code:               "HNL",
		DigitsAfterDecimal: 2,
	},
	"JPY": currency{
		Name:               "Yen",
		Code:               "JPY",
		DigitsAfterDecimal: 0,
	},
	"XAF": currency{
		Name:               "CFA Franc BEAC",
		Code:               "XAF",
		DigitsAfterDecimal: 0,
	},
	"UGX": currency{
		Name:               "Uganda Shilling",
		Code:               "UGX",
		DigitsAfterDecimal: 0,
	},
	"CLF": currency{
		Name:               "Unidad de Fomento",
		Code:               "CLF",
		DigitsAfterDecimal: 4,
	},
	"GMD": currency{
		Name:               "Dalasi",
		Code:               "GMD",
		DigitsAfterDecimal: 2,
	},
	"TND": currency{
		Name:               "Tunisian Dinar",
		Code:               "TND",
		DigitsAfterDecimal: 3,
	},
	"AZN": currency{
		Name:               "Azerbaijanian Manat",
		Code:               "AZN",
		DigitsAfterDecimal: 2,
	},
	"JOD": currency{
		Name:               "Jordanian Dinar",
		Code:               "JOD",
		DigitsAfterDecimal: 3,
	},
	"KES": currency{
		Name:               "Kenyan Shilling",
		Code:               "KES",
		DigitsAfterDecimal: 2,
	},
	"STD": currency{
		Name:               "Dobra",
		Code:               "STD",
		DigitsAfterDecimal: 2,
	},
	"ZAR": currency{
		Name:               "Rand",
		Code:               "ZAR",
		DigitsAfterDecimal: 2,
	},
	"GIP": currency{
		Name:               "Gibraltar Pound",
		Code:               "GIP",
		DigitsAfterDecimal: 2,
	},
	"RWF": currency{
		Name:               "Rwanda Franc",
		Code:               "RWF",
		DigitsAfterDecimal: 0,
	},
	"AWG": currency{
		Name:               "Aruban Florin",
		Code:               "AWG",
		DigitsAfterDecimal: 2,
	},
	"MMK": currency{
		Name:               "Kyat",
		Code:               "MMK",
		DigitsAfterDecimal: 2,
	},
	"THB": currency{
		Name:               "Baht",
		Code:               "THB",
		DigitsAfterDecimal: 2,
	},
	"NOK": currency{
		Name:               "Norwegian Krone",
		Code:               "NOK",
		DigitsAfterDecimal: 2,
	},
	"CHF": currency{
		Name:               "Swiss Franc",
		Code:               "CHF",
		DigitsAfterDecimal: 2,
	},
	"HUF": currency{
		Name:               "Forint",
		Code:               "HUF",
		DigitsAfterDecimal: 2,
	},
	"IRR": currency{
		Name:               "Iranian Rial",
		Code:               "IRR",
		DigitsAfterDecimal: 2,
	},
	"NZD": currency{
		Name:               "New Zealand Dollar",
		Code:               "NZD",
		DigitsAfterDecimal: 2,
	},
	"BSD": currency{
		Name:               "Bahamian Dollar",
		Code:               "BSD",
		DigitsAfterDecimal: 2,
	},
	"AOA": currency{
		Name:               "Kwanza",
		Code:               "AOA",
		DigitsAfterDecimal: 2,
	},
	"CNY": currency{
		Name:               "Yuan Renminbi",
		Code:               "CNY",
		DigitsAfterDecimal: 2,
	},
	"BDT": currency{
		Name:               "Taka",
		Code:               "BDT",
		DigitsAfterDecimal: 2,
	},
	"MWK": currency{
		Name:               "Malawi Kwacha",
		Code:               "MWK",
		DigitsAfterDecimal: 2,
	},
	"MGA": currency{
		Name:               "Malagasy Ariary",
		Code:               "MGA",
		DigitsAfterDecimal: 2,
	},
	"BAM": currency{
		Name:               "Convertible Mark",
		Code:               "BAM",
		DigitsAfterDecimal: 2,
	},
	"MKD": currency{
		Name:               "Denar",
		Code:               "MKD",
		DigitsAfterDecimal: 2,
	},
	"GEL": currency{
		Name:               "Lari",
		Code:               "GEL",
		DigitsAfterDecimal: 2,
	},
	"MZN": currency{
		Name:               "Mozambique Metical",
		Code:               "MZN",
		DigitsAfterDecimal: 2,
	},
	"PLN": currency{
		Name:               "Zloty",
		Code:               "PLN",
		DigitsAfterDecimal: 2,
	},
	"TTD": currency{
		Name:               "Trinidad and Tobago Dollar",
		Code:               "TTD",
		DigitsAfterDecimal: 2,
	},
	"PYG": currency{
		Name:               "Guarani",
		Code:               "PYG",
		DigitsAfterDecimal: 0,
	},
	"OMR": currency{
		Name:               "Rial Omani",
		Code:               "OMR",
		DigitsAfterDecimal: 3,
	},
	"VUV": currency{
		Name:               "Vatu",
		Code:               "VUV",
		DigitsAfterDecimal: 0,
	},
	"ZMW": currency{
		Name:               "Zambian Kwacha",
		Code:               "ZMW",
		DigitsAfterDecimal: 2,
	},
	"ILS": currency{
		Name:               "New Israeli Sheqel",
		Code:               "ILS",
		DigitsAfterDecimal: 2,
	},
	"VND": currency{
		Name:               "Dong",
		Code:               "VND",
		DigitsAfterDecimal: 0,
	},
	"DOP": currency{
		Name:               "Dominican Peso",
		Code:               "DOP",
		DigitsAfterDecimal: 2,
	},
}
