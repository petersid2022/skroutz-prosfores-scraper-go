package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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

	out := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", shortLink, y)
	return out
}

func scrapeProductInfo(c *colly.Collector, category string, wg *sync.WaitGroup, results chan<- result) {
	defer wg.Done()

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

	c.Visit("https://www.skroutz.gr/prosfores?order_by=" + category + "&recent=1")
}

func main() {
	startTime := time.Now()

	categoryPtr := flag.String("f", "Recommended", "[Recommended], [price_asc], [price_desc], [newest], [pricedrop]")

	var productsInfo []productInformation
	c := colly.NewCollector()

	var wg sync.WaitGroup
	results := make(chan result)

	go func() {
		wg.Wait()
		close(results)
	}()

	numWorkers := 10 
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go scrapeProductInfo(c.Clone(), *categoryPtr, &wg, results)
	}

	flag.Parse()

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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Old Price", "New Price", "Link"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetRowLine(true)
	table.SetCaption(true, "petrside 2023 / Category: "+*categoryPtr)

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
		table.Render()
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
		table.Render()

	}
	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time: " + elapsedTime.String())
}

