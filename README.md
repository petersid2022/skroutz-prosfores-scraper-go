# skroutz-prosfores-scraper-go
CLI based web scraper written in Go, that scrapes www.skroutz.gr/price-drops for any new deals.

![image](./examples/demo.gif)
> Gif made using [vhs](https://github.com/charmbracelet/vhs) 
## Installation using Make
Make sure you are running ```go version 1.20.x```

```bash
git clone https://github.com/petersid2022/skroutz-prosfores-scraper-go.git
cd skroutz-prosfores-scraper-go
make
```
**Pro Tip:** Place the binary in the ```~/.local/bin``` directory. This enables you to run the command from the terminal regardless of your location!

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

## Things I am working on  
* ~~Command-Line flags (change category, add filters, prices ascending/descending etc.)~~ [f155812](https://github.com/petersid2022/skroutz-web-scraper-go/commit/f155812c9a3daa9e19d77b2e0e172f07b4546021)

* ~~Use a different TUI Library~~

* ~~Use goroutines (Concurrency)~~

## License
This project is licensed under the MIT License. Please see the [LICENSE](./LICENSE) file for more details.
