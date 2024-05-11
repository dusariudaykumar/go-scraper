package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

const URL = "https://www.etsy.com/c/clothing/mens-clothing"

type Product struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Cost        float32 `json:"cost"`
	Currency    string  `json:"currency"`
	Image       string  `json:"image"`
	ProductLink string  `json:"product_link"`
}

func main() {
	fmt.Println("Go Scraper....")

	c := colly.NewCollector(colly.AllowedDomains("etsy.com", "www.etsy.com"))

	products, err := EtsyScraper(c, URL)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total product: ", len(*products))

	// converting products slice into json
	content, err := json.Marshal(products)

	if err != nil {
		fmt.Println("Something went wrong while converting into json: ", err)
	}

	os.WriteFile("products.json", content, 0644)
}

func EtsyScraper(c *colly.Collector, url string) (*[]Product, error) {

	var productsData []Product

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping: ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visted: ", r.Request.URL)
		fmt.Println("Status: ", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Requested URL: ", r.Request.URL, "\nfailed with status: ", r.StatusCode, "\nError: ", err)
	})

	c.OnHTML(".v2-listing-card", func(h *colly.HTMLElement) {
		title := h.ChildText(".v2-listing-card__title")
		id := h.Attr("data-palette-listing-id")
		cost, _ := strconv.ParseFloat(h.ChildText(".n-listing-card__price .lc-price .currency-value"), 32)
		currency := h.ChildText(".n-listing-card__price .lc-price .currency-symbol")
		image := h.ChildAttr(".v2-listing-card__img img", "src")
		productLink := fmt.Sprintf("https://www.etsy.com/listing/%v", id)
		product := &Product{Id: id, Title: title, Cost: float32(cost), Currency: currency, Image: image, ProductLink: productLink}
		productsData = append(productsData, *product)

	})
	c.Wait()

	err := c.Visit(URL)

	if err != nil {
		return nil, err
	}

	return &productsData, nil
}
