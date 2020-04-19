package etracking

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var tis620Decoder = charmap.Windows874.NewDecoder()

func convertDecimalStringToSatang(text string) (int64, error) {
	textPart := strings.Split(text, ".")
	if len(textPart) != 2 || len(textPart[1]) != 2 {
		return 0, errors.New("invalid format")
	}
	baht, err := strconv.ParseInt(textPart[0], 10, 64)
	if err != nil {
		return 0, err
	}
	satang, err := strconv.ParseInt(textPart[1], 10, 64)
	if err != nil {
		return 0, err
	}

	satang += baht * 100

	return satang, nil
}

func convertTIS620ToUTF8(text string) (string, error) {
	return tis620Decoder.String(text)
}
