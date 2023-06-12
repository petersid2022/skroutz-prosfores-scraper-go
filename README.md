# skroutz-prosfores-scraper-go
CLI based web scraper written in Go, that scrapes www.skroutz.gr/prosfores for any new deals.

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

* Use a different TUI Library 

## Things I am working on  
* ~~Command-Line flags (change category, add filters, prices ascending/descending etc.)~~ [f155812](https://github.com/petersid2022/skroutz-web-scraper-go/commit/f155812c9a3daa9e19d77b2e0e172f07b4546021)

* ~~Use goroutines (Concurrency)~~

* Reuse HTTP client, instead of creating a new one for every request. 

## License

(The MIT License)

Copyright (c) 2023 Peter Sideris petersid2022@gmail.com 

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
