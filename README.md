# gopartpicker

A scraper for [pcpartpicker.com](https://pcpartpicker.com) for Go.

It is implemented using [Colly](https://github.com/gocolly/colly).

[![GoDoc](https://godoc.org/github.com/quakecodes/gopartpicker?status.svg)](https://pkg.go.dev/github.com/quakecodes/gopartpicker)
[![Report Card](https://img.shields.io/badge/report%20card-A-00b800?style=flat-square)](http://goreportcard.com/report/quakecodes/gopartpicker)

## Features

- Extract data from part list URLs
- Search for parts
- Extract data from part URLs
- Match PCPartPicker URLs using RegExp
- Support for multiple regions and currencies

## Installation

It is assumed that you have Go 1.17 installed.

```
go get github.com/quakecodes/gopartpicker
```

## Usage

Import the library.

```go
import "github.com/quakecodes/gopartpicker"
```

Create a new scraper instance.

```go
scraper := gopartpicker.NewScraper()
```

Fetch a part via URL.

```go
part, err := scraper.GetPart("https://uk.pcpartpicker.com/product/g94BD3/amd-ryzen-5-5600x-37-ghz-6-core-processor-100-100000065box")
if err != nil {
    log.Fatal(err)
}

fmt.Println(part.Name)
```

Fetch a part list via URL.

```go
partList, err := scraper.GetPartList("https://uk.pcpartpicker.com/list/LNqWbh")
if err != nil {
    log.Fatal(err)
}

fmt.Println(partList.Parts[0].Name)
```

Search for parts via a search term. The second argument is the region to search with.

```go
parts, err := scraper.SearchParts("ryzen 5 3600", "uk")

// Some searches redirect to a product page, if you know that what you are searching will not redirect
// then you do not need to do the type assertion and if statement
_, ok := err.(*gopartpicker.RedirectError)

if ok {
    // RedirectError.Error returns the URL of the redirect
    part, err := scraper.GetPart(err.Error())

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(part.Name)
} else if err != nil {
    log.Fatal(err)
} else {
    fmt.Println(parts[0].Name)
}
```

Set headers for subsequent requests.

```go
// "global" sets headers for all requests
scraper.SetHeaders("global", map[string]string{
  "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36",
})

// Alternatively, make headers site specific by passing in a domain
scraper.SetHeaders("pcpartpicker.com", map[string]string{
  "accept": "accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
})

// Make sure you set headers before calling request methods
scraper.GetPart("...")
```

# Documentation

As of right now, there is no full documentation.

Feel free to ask for help by asking QuaKe in [this Discord server](https://discord.com/invite/WM9pHp8).
