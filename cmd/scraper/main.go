package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
)

type productInformation struct {
	name, link, oldprice, newprice string
}

type result struct {
	productInfo productInformation
	err         error
}

func shortenURL(longURL string) (shortURL string, err error) {
	resp, err := http.Get("http://tinyurl.com/api-create.php?url=" + longURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func link(x string, y string) string {
	shortLink, err := shortenURL(x)
	if err != nil {
		fmt.Println(err)
	}
    //create a hyperlink in the terminal
	out := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", shortLink, y)
	return out
}

func scrapeProductInfo(c *colly.Collector, category string, wg *sync.WaitGroup, results chan<- result) {
	defer wg.Done()
    scrapedPages := 5

    //scrape the product information
	c.OnHTML(".sku-card.js-sku", func(e *colly.HTMLElement) {
		productInfo := productInformation{}
		productInfo.name = e.ChildAttr("a", "title")
		productInfo.link = "https://skroutz.gr" + e.ChildAttr("a", "href")
		u, _ := url.Parse(productInfo.link)
		productInfo.link = "https://www.skroutz.gr" + u.Path
		productInfo.oldprice = e.ChildText(".sku-card-info")
		productInfo.oldprice = strings.Split(productInfo.oldprice, " ")[0]
		productInfo.newprice = e.ChildText(".product-link.js-sku-link")

		results <- result{
			productInfo: productInfo,
			err:         nil,
		}
	})

    //scrape the number of pages
	for i := 1; i <= scrapedPages; i++ {
		url := "https://www.skroutz.gr/prosfores?order_by=" + category + "&recent=1&page=" + strconv.Itoa(i)
		c.Visit(url)
	}
}

func main() {
    c := colly.NewCollector()
    startTime := time.Now()

    //create a flag to choose the category
    categoryPtr := flag.String("f", "Recommended", "[Recommended], [price_asc], [price_desc], [newest]")

    //create a slice to store the product information
	var productsInfo []productInformation

    //create a waitgroup to wait for all the goroutines to finish and a channel to store the results
	var wg sync.WaitGroup
	results := make(chan result)

	go func() {
		wg.Wait()
		close(results)
	}()

    //create 10 goroutines to scrape the product information
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go scrapeProductInfo(c.Clone(), *categoryPtr, &wg, results)
	}

	flag.Parse()

    //store the results in the slice
	for res := range results {
		if res.err != nil {
			fmt.Println(res.err)
		} else {
			productsInfo = append(productsInfo, res.productInfo)
		}
	}

	rand.Seed(time.Now().UnixNano())
	numProductsToPrint := 5
	if len(productsInfo) < numProductsToPrint {
		numProductsToPrint = len(productsInfo)
	}
	indices := rand.Perm(len(productsInfo))[:numProductsToPrint]

    //create a table to print the results, currently there is a bug in the TableWriter library that 
    //does not render the border correctly if there hyperlinks in the table. see issue #212
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Old Price", "New Price", "Link"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})
	table.SetBorder(false)
	table.SetRowLine(true)
	table.SetCenterSeparator("|")
	table.SetCaption(true, "petrside 2023 / Category: "+*categoryPtr)

    //if no flags are given (i.e. recommended category) print 5 random products,
    //else, if there are flags being given, print all the products
	if *categoryPtr == "Recommended" {
		for _, i := range indices {
			productInfo := productsInfo[i]

			out := []string{
				strings.TrimSpace(productInfo.name),
				strings.TrimSpace(productInfo.oldprice + " €"),
				strings.TrimSpace(productInfo.newprice),
				link(productInfo.link, "link"),
			}
			table.Append(out)
		}
	} else {
		for i := 0; i < len(productsInfo); i++ {
			productInfo := productsInfo[i]
			out := []string{
				strings.TrimSpace(productInfo.name),
				strings.TrimSpace(productInfo.oldprice + " €"),
				strings.TrimSpace(productInfo.newprice),
				link(productInfo.link, "link"),
			}
			table.Append(out)
		}
	}

    //print the table and the elapsed time
	table.Render()
	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time: " + elapsedTime.String())
}
