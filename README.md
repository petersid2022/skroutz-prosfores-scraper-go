# skroutz-prosfores-scraper-go
CLI based web scraper written in Go, that scrapes www.skroutz.gr/prosfores for any new deals.

![image](./examples/demo.gif)
> Gif made using [vhs](https://github.com/charmbracelet/vhs) 
## Installation 
Make sure you are running ```go version 1.20.x```

```bash
git clone https://github.com/petersid2022/skroutz-prosfores-scraper-go.git
cd skroutz-prosfores-scraper-go
go build ./cmd/scraper/main.go 
```
**TIP:** Move the binary to ```~/.local/bin```. That way, you can execute the command from the terminal, wherever you are!

## CLI Usage
```skroutz [Options]```

Options:

```
-h, --help      Output usage information
-f, [String]    Set filtering options: [Recommended], [price_asc], [price_desc], [newest]
-p, [Number]    Set the number of pages to scrape (default: 5)
-n, [Number]    Set the number of products to print when filtering_option=Recommended (default: 5)
-w, [Number]    Set the number of workers (default: 10)
```

## Possible Optimizations

* ~~Use goroutines (Concurrency)~~

* Reuse HTTP client, instead of creating a new one for every request. 

* Use a caching mechanism

* Use a better URL shortener service (Bitly or Google URL Shortener)

* Use a better HTML parsing library (GoQuery)

## Things I am working on  
* ~~Command-Line flags (change category, add filters, prices ascending/descending etc.)~~ [f155812](https://github.com/petersid2022/skroutz-web-scraper-go/commit/f155812c9a3daa9e19d77b2e0e172f07b4546021)

* ~~Use a different TUI Library~~

* ~~Use goroutines (Concurrency)~~

* Reuse HTTP client, instead of creating a new one for every request. 

## License
This project is licensed under the MIT License. Please see the [LICENSE](./LICENSE) file for more details.
