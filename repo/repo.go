package repo

import (
	"encoding/csv"
	"os"
	"strconv"
)

type currency struct {
	currency_name string
	price         float64
}

type currencyRepo struct {
	currencies map[string]*currency
}

func NewCurrencyRepo(currencyPath string) (cr *currencyRepo, err error) {
	currencyRecords, err := readCSV(currencyPath)
	if err != nil {
		return nil, err
	}

	currencies := make(map[string]*currency, 0)
	for _, record := range currencyRecords {
		price, _ := strconv.ParseFloat(record[1], 64)

		p := &currency{
			currency_name: record[0],
			price:         price,
		}

		currencies[p.currency_name] = p
	}

	cr = &currencyRepo{
		currencies: currencies,
	}

	return cr, nil
}

func readCSV(path string) (lines [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (cr *currencyRepo) FetchPrice(currency_name string) (price float64, found bool) {
	p, ok := cr.currencies[currency_name]
	if !ok {
		return 0.0, false
	}

	return p.price, true
}
