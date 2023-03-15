package main

import (
    "github.com/gocolly/colly"
    "github.com/olekukonko/tablewriter"
    "os"
    "strings"
    "fmt"
    "net/url"
    "io/ioutil"
    "net/http"
)

type productInformation struct {
    name, link, oldprice, newprice string
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


func main() {
    var productsInfo []productInformation

    c := colly.NewCollector(
        colly.MaxDepth(2),
    )

    c.OnHTML(".js-sku.cf.card", func(e *colly.HTMLElement) {
        productInfo := productInformation{}

        productInfo.name = e.ChildAttr("a", "title")
        productInfo.link = "https://skroutz.gr" + e.ChildAttr("a", "href")
        u, _ := url.Parse(productInfo.link)
        productInfo.link  = "https://www.skroutz.gr" + u.Path
        productInfo.oldprice = e.ChildText(".price")
        productInfo.oldprice = strings.Split(productInfo.oldprice, " ")[0]
        productInfo.newprice = e.ChildText(".product-link.js-sku-link")

        productsInfo = append(productsInfo, productInfo)
    })

    //c.Visit("https://www.skroutz.gr/prosfores?order_by=newest")
    c.Visit("https://www.skroutz.gr/prosfores?order_by=recommended&recent=1")

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Old Price", "New Price", "Link"})
    table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
    table.SetAlignment(tablewriter.ALIGN_LEFT)
    table.SetBorder(true)
    table.SetRowLine(true)

    for _, productInfo := range productsInfo {
        out := []string{
            strings.TrimSpace(productInfo.name),
            strings.TrimSpace(productInfo.newprice),
            strings.TrimSpace(productInfo.oldprice + " €"),
            link(productInfo.link, "link"),
        }
        table.Append(out)
    }
    table.Render()
}
