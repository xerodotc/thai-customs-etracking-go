package etracking

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const initialSelector = `form[name="frmETSQ2040"] > table > tbody > tr:nth-child(2) > td > table > tbody`

func Lookup(trackingNo string) (TaxResult, error) {
	return LookupWithClient(http.DefaultClient, trackingNo)
}

func LookupWithClient(client *http.Client, trackingNo string) (TaxResult, error) {
	q := url.Values{}
	q.Add("act", "SRH")
	q.Add("pge", "")
	q.Add("rowPerPge", "")
	q.Add("pclNum", trackingNo)
	postBody := q.Encode()

	req, err := http.NewRequest(http.MethodPost, ETrackingPostURL, strings.NewReader(postBody))
	if err != nil {
		return TaxResult{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return TaxResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return TaxResult{}, fmt.Errorf("status code %d returned", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return TaxResult{}, err
	}

	tables := doc.Find(initialSelector).Find(`tr`)
	infoTable := tables.First().Find(`table`)

	if html, err := infoTable.Html(); err != nil || html == "" {
		if err != nil {
			return TaxResult{}, err
		}
		return TaxResult{}, errors.New("something went wrong")
	}

	if warning := infoTable.Find(`font[color="red"]`).Text(); warning != "" {
		return TaxResult{Available: false}, nil
	}

	result := TaxResult{Available: true}

	row := infoTable.Find(`tr`).First()
	col := row.Find(`td`).First().Next()
	result.Barcode = strings.TrimSpace(col.Text())

	col = col.Next().Next()
	stg, err := convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.ImportTax = stg

	row = row.Next()
	col = row.Find(`td`).First().Next()
	result.CustomID = strings.TrimSpace(col.Text())

	col = col.Next().Next()
	stg, err = convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.ExciseTax = stg

	row = row.Next()
	col = row.Find(`td`).First().Next()
	str, err := convertTIS620ToUTF8(col.Text())
	if err != nil {
		return TaxResult{}, err
	}
	result.Recipient = strings.TrimSpace(str)

	col = col.Next().Next()
	stg, err = convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.InteriorTax = stg

	row = row.Next()
	col = row.Find(`td`).First().Next()
	str, err = convertTIS620ToUTF8(col.Text())
	if err != nil {
		return TaxResult{}, err
	}
	result.ReceivingLocation = strings.TrimSpace(str)

	col = col.Next().Next()
	stg, err = convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.ValueAddedTax = stg

	row = row.Next()
	col = row.Find(`td`).First().Next().Next()
	stg, err = convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.OtherFee = stg

	row = row.Next()
	col = row.Find(`td`).First().Next().Next().Next()
	stg, err = convertDecimalStringToSatang(strings.TrimSpace(col.Text()))
	if err != nil {
		return TaxResult{}, err
	}
	result.TotalTax = stg

	stepTable := tables.First().Next().Next().Next().Find(`table`)

	result.Steps = make([]CustomStepEntry, 0)

	for iterator := stepTable.Find(`tr`).First().Next(); iterator.Text() != ""; iterator = iterator.Next() {
		col := iterator.Find(`td`).First().Next()
		desc, err := convertTIS620ToUTF8(col.Text())
		if err != nil {
			return TaxResult{}, err
		}
		desc = strings.TrimSpace(desc)
		col = col.Next()
		t, err := parseETrackingTime(strings.TrimSpace(col.Text()))
		if err != nil {
			return TaxResult{}, err
		}
		s := CustomStepEntry{
			Step: desc,
			Time: t,
		}
		result.Steps = append(result.Steps, s)
	}

	return result, nil
}
