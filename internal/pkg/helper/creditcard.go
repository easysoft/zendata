package helper

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCardType = "visa"
)

func GenerateCreditCard(cardType string) (ret string) {
	cardProperties, exists := AvailableCardTypes[cardType]

	if !exists {
		cardType = defaultCardType
		cardProperties, _ = AvailableCardTypes[cardType]
	}

	rand.Seed(time.Now().UnixNano())

	var card = CreditCard{
		Issuer:     cardProperties.LongName,
		Pan:        GeneratePAN(cardProperties),
		ExpiryDate: GenerateExpiryDate(),
		CVV:        GenerateCVV(cardProperties.CvvSize),
	}

	ret = card.Pan.Raw

	return
}

type CreditCard struct {
	Issuer     string     `json:"issuer"`
	Pan        PAN        `json:"pan"`
	ExpiryDate ExpiryDate `json:"expiryDate"`
	CVV        string     `json:"cvv"`
}

type ExpiryDate struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type PAN struct {
	Raw       string `json:"raw"`
	Formatted string `json:"formatted"`
}
type CardProperties struct {
	LongName string
	Prefix   []string
	PanSize  int
	CvvSize  int
}

var AvailableCardTypes = map[string]CardProperties{
	"amex": {
		LongName: "American Express",
		Prefix:   []string{"37", "34"},
		PanSize:  15,
		CvvSize:  4,
	},
	"visa": {
		LongName: "Visa",
		Prefix:   []string{"4"},
		PanSize:  16,
		CvvSize:  3,
	},
	"mc": {
		LongName: "Mastercard",
		Prefix:   []string{"51", "52", "53", "54", "55"},
		PanSize:  16,
		CvvSize:  3,
	},
	"dci": {
		LongName: "Diners Club International",
		Prefix:   []string{"36", "38"},
		PanSize:  16,
		CvvSize:  3,
	},
	"jcb": {
		LongName: "Japan Credit Bureau",
		Prefix:   []string{"35"},
		PanSize:  16,
		CvvSize:  3,
	},
	"discover": {
		LongName: "Discover",
		Prefix:   []string{"6011", "65"},
		PanSize:  16,
		CvvSize:  3,
	},
}

func calculateLuhn(digits string) string {
	number, _ := strconv.Atoi(digits)
	checkNumber := checksum(number)
	if checkNumber == 0 {
		return "0"
	}

	return strconv.Itoa(10 - checkNumber)
}

func checksum(number int) int {
	var sum int
	for i := 0; number > 0; i++ {
		digit := number % 10
		if i%2 == 0 { // even
			digit = digit * 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		sum += digit
		number = number / 10
	}

	return sum % 10
}

func GeneratePAN(properties CardProperties) PAN {
	var prefix = properties.Prefix[rand.Intn(len(properties.Prefix))]
	for len(prefix) < properties.PanSize-1 {
		prefix = prefix + generateRandomNumberOfSize(1)
	}

	number := prefix + calculateLuhn(prefix)
	return PAN{
		Raw:       number,
		Formatted: FormatPan(properties.LongName, number),
	}
}

func FormatPan(cardType string, pan string) string {
	var formattedPan string
	for i := 0; i < len(pan); i++ {
		formattedPan = formattedPan + string(pan[i])
		switch cardType {
		case "American Express":
			if i == 3 || i == 9 {
				formattedPan = formattedPan + " "
			}
		default:
			if (i+1)%4 == 0 {
				formattedPan = formattedPan + " "
			}
		}
	}

	return strings.TrimSpace(formattedPan)
}

func generateRandomNumberOfSize(size int) string {
	var numbers = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var randomNumber []string
	for i := 0; i < size; i++ {
		randomNumber = append(randomNumber, numbers[rand.Intn(len(numbers))])
	}
	return strings.Join(randomNumber, "")
}

func GenerateCVV(size int) string {
	return generateRandomNumberOfSize(size)
}

func GenerateExpiryDate() ExpiryDate {
	var year = time.Now().Year() + rand.Intn(6)
	var month = rand.Intn(12) + 1

	return ExpiryDate{
		Month: month,
		Year:  year,
	}
}
