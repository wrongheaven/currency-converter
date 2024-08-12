package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/joho/godotenv"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

var (
	fromCurr  string
	toCurr    string
	amount    float64
	converted float64
)

func init() {
	rootPath := os.Getenv("CCONV_PATH")
	envPath := filepath.Join(rootPath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}
}

func convert() {
	app_id := os.Getenv("APP_ID")
	url := "https://openexchangerates.org/api/latest.json"

	resp, err := http.Get(fmt.Sprintf("%s?app_id=%s", url, app_id))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("foo")
		log.Fatal(err)
	}

	var data ApiResponse
	json.Unmarshal(body, &data)

	if data.IsError {
		fmt.Println(data.ErrorMessage)
		os.Exit(1)
	}

	fromRate, toRate := data.Rates[fromCurr], data.Rates[toCurr]
	converted = amount * (toRate / fromRate)
}

func main() {
	var amountStr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("From").
				Value(&fromCurr),
			huh.NewInput().
				Title("To").
				Value(&toCurr),
			huh.NewInput().
				Title("Amount").
				Value(&amountStr).
				Validate(func(str string) error {
					var err error
					if amount, err = strconv.ParseFloat(amountStr, 64); err != nil {
						return err
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fromCurr = strings.ToUpper(fromCurr)
	toCurr = strings.ToUpper(toCurr)

	spinner.New().
		Title("Converting ...").
		Action(convert).
		Run()

	// fmt.Printf("%s %.2f -> %s %.2f\n", fromCurr, amount, toCurr, converted)

	printer := message.NewPrinter(language.English)
	printer.Printf(
		"%s %v -> %s %v\n",
		fromCurr,
		number.Decimal(amount, number.MinFractionDigits(2), number.MaxFractionDigits(2)),
		toCurr,
		number.Decimal(converted, number.MinFractionDigits(2), number.MaxFractionDigits(2)),
	)
}
