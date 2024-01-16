package main

import (
	//"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type productInformation struct {
	name, link, oldprice, newprice string
}

type result struct {
	productInfo productInformation
	err         error
}

// func shortenURL(longURL string) (shortURL string, err error) {
//     resp, err := http.Get("http://tinyurl.com/api-create.php?url=" + longURL)
//     if err != nil {
//         return "", err
//     }
//     defer resp.Body.Close()
//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return "", err
//     }
//     return string(body), nil
// }

func createLink(x string, y string) string {
	// shortLink, err := shortenURL(x)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// create a hyperlink in the terminal
	out := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", x, y)
	return out
}

func scrapeProductInfo(c *colly.Collector, category string, wg *sync.WaitGroup, results chan<- result, scrapedPages int) {
	defer wg.Done()

	// scrape the product information
	c.OnHTML(".sku-card.js-sku", func(e *colly.HTMLElement) {
		productInfo := productInformation{}
		productInfo.name = e.ChildAttr("a", "title")
		productInfo.link = "https://skroutz.gr" + e.ChildAttr("a", "href")
		u, err := url.Parse(productInfo.link)
		if err != nil {
			log.Fatalf("Failed to parse url: %s\n", err)
		}
		productInfo.link = "https://www.skroutz.gr" + u.Path
		productInfo.oldprice = e.ChildText(".sku-card-info")
		productInfo.oldprice = strings.Split(productInfo.oldprice, " ")[0]
		productInfo.newprice = e.ChildText(".product-link.js-sku-link")

		results <- result{
			productInfo: productInfo,
			err:         nil,
		}
	})

	// scrape the number of pages
	url := "https://www.skroutz.gr/price-drops?order_by=" + category + "&recent=1&page=" + strconv.Itoa(scrapedPages)
	c.Visit(url)
}

func replaceNaN(x string) string {
	out := ""
	x = strings.ReplaceAll(x, ".", "")
	x = strings.ReplaceAll(x, ",", ".")
	out = x
	return out
}

func main() {
	c := colly.NewCollector()
	startTime := time.Now()

	// create a flag to choose the category
	categoryPtr := flag.String("f", "Recommended", "[Recommended], [price_asc], [price_desc], [newest]")

	// create a flag to choose the number of pages to scrapedPages
	pagesPtr := flag.Int("p", 5, "number of pages to scrape")

	// create a flag to choose the number of products to Println
	productsPtr := flag.Int("n", 5, "number of products to print")

	// create a flag to choose the number of numWorkers
	numWorkersPtr := flag.Int("w", 10, "number of workers")

	// create a slice to store the product information
	var productsInfo []productInformation

	// create a waitgroup to wait for all the goroutines to finish and a channel to store the results
	var wg sync.WaitGroup
	results := make(chan result)

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *numWorkersPtr; i++ {
		wg.Add(1)
		go scrapeProductInfo(c.Clone(), *categoryPtr, &wg, results, *pagesPtr)
	}

	flag.Parse()

	// store the results in the slice
	for res := range results {
		if res.err != nil {
			log.Fatalf("Failed to store the results in the slice: %s\n", res.err)
		} else {
			productsInfo = append(productsInfo, res.productInfo)
		}
	}

	numProductsToPrint := *productsPtr
	if len(productsInfo) < numProductsToPrint {
		numProductsToPrint = len(productsInfo)
	}
	indices := rand.Perm(len(productsInfo))[:numProductsToPrint]

	// Create a new table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	// Add table headers
	t.AppendHeader(table.Row{"Name", "Old Price", "New Price", "Link"})

	// Add table rows
	for _, i := range indices {
		productInfo := productsInfo[i]

		newOldPrice := replaceNaN(productInfo.oldprice)
		newNewPrice := replaceNaN(productInfo.newprice)

		oldprice_val, err := strconv.ParseFloat(strings.Split(newOldPrice, " ")[0], 64)
		if err != nil {
			log.Fatalf("Failed to parse oldprice value: %s\n", err)
		}

		newprice_val, err := strconv.ParseFloat(strings.Split(newNewPrice, " ")[0], 64)
		if err != nil {
			log.Fatalf("Failed to parse oldprice value: %s\n", err)
		}

		productInfo.newprice += fmt.Sprintf("\n(-%.2f%%)", (100 * (oldprice_val - newprice_val) / oldprice_val))

		row := make([]interface{}, 4)
		row[0] = strings.TrimSpace(productInfo.name)
		t.AppendSeparator()
		row[1] = strings.TrimSpace(productInfo.oldprice + " €")
		t.AppendSeparator()
		row[2] = strings.TrimSpace(productInfo.newprice)
		t.AppendSeparator()
		row[3] = createLink(productInfo.link, "link")

		t.AppendRow(row)
	}

	// Render the table
	t.SetStyle(table.StyleRounded)
	t.SetAutoIndex(true)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignHeader: text.AlignCenter, WidthMax: 48},
		{Number: 2, AlignHeader: text.AlignCenter, WidthMax: 48},
		{Number: 3, AlignHeader: text.AlignCenter, WidthMax: 48},
		{Number: 4, AlignHeader: text.AlignCenter, WidthMax: 48},
	})
	t.SetCaption("petrside 2024 / Category: " + *categoryPtr)
	t.Render()

	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time: " + elapsedTime.String())
}
