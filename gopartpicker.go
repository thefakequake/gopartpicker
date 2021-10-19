package gopartpicker

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/gocolly/colly"
)

var (
	partListClassMappings = map[string]string{
		"Base":     ".td__base",
		"Promo":    ".td__promo",
		"Shipping": ".td__shipping",
		"Tax":      ".td__tax",
		"Price":    ".td__price",
	}

	partClassMappings = map[string]string{
		"Base":     ".td__base",
		"Promo":    ".td__promo",
		"Shipping": ".td__shipping",
		"Tax":      ".td__tax",
		"Total":    ".td__finalPrice",
	}
)

type Scraper struct {
	// The Colly collector used for scraping.
	Collector *colly.Collector
	Headers   map[string]map[string]string
}

type RedirectError struct {
	URL string
}

func (r RedirectError) Error() string {
	return r.URL
}

// Extracts the name of a vendor from a PCPartPicker affiliate link.
func ExtractVendorName(URL string) string {
	return strings.Split(URL, "/")[2]
}

// Creates a new Scraper instance.
func NewScraper() Scraper {
	col := colly.NewCollector()
	col.Async = true
	col.AllowURLRevisit = true

	s := Scraper{
		Collector: col,
	}
	s.Headers = map[string]map[string]string{
		"global": {},
	}

	return s
}

// Sets headers for subsequent requests on a specific site - set site to "pcpartpicker.com" for PCPP or "global" for all sites.
func (s *Scraper) SetHeaders(site string, newHeaders map[string]string) {
	s.Headers[site] = newHeaders

	for k, v := range newHeaders {
		s.Headers[site][k] = v
	}

	s.Collector.OnRequest(func(r *colly.Request) {
		headers := s.Headers["global"]
		for k, v := range s.Headers[r.URL.Hostname()] {
			headers[k] = v
		}

		for k, v := range headers {
			if len(k) > 0 && len(v) > 0 {
				r.Headers.Set(k, v)
			}
		}
	})
}

// Fetches data on a part list via URL.
func (s Scraper) GetPartList(URL string) (*PartList, error) {
	if !MatchPCPPURL(URL) {
		return nil, errors.New("invalid PCPartPicker URL")
	}
	URL = ConvertListURL(URL)

	var partList PartList

	s.Collector.OnHTML(".partlist__wrapper", func(el *colly.HTMLElement) {
		parts := []PartListPart{}

		el.ForEach(".tr__product", func(i int, prod *colly.HTMLElement) {
			partVendor := Vendor{
				InStock: false,
				Price:   Price{},
			}

			for k, v := range partListClassMappings {
				toParse := prod.ChildText(v)

				if strings.HasSuffix(toParse, "No Prices Available") || toParse == "FREE" {
					continue
				}
				stringPrice := strings.Replace(toParse, k, "", 1)
				price, curr, _ := StringPriceToFloat(stringPrice)

				switch k {
				case "Base":
					partVendor.Price.Base = price
				case "Promo":
					partVendor.Price.Discounts = price
				case "Shipping":
					partVendor.Price.Shipping = price
				case "Tax":
					partVendor.Price.Shipping = price
				case "Price":
					partVendor.Price.TotalString = strings.TrimSpace(stringPrice)
					partVendor.Price.Total = price
					partVendor.Price.Currency = curr
					partVendor.InStock = true
				}
			}

			if partVendor.InStock {
				partVendor.URL = "https://" + el.Request.URL.Host + prod.ChildAttr(".td__where a", "href")
				partVendor.Image = "https:" + prod.ChildAttr(".td__where a img", "src")
				partVendor.Name = ExtractVendorName(partVendor.URL)
			}

			part := PartListPart{
				Type:   prod.ChildText(".td__component"),
				Name:   prod.ChildText(".td__name"),
				Image:  prod.ChildAttr(".td__image a", "href"),
				URL:    "https://" + el.Request.URL.Host + prod.ChildAttr(".td__name a", "href"),
				Vendor: partVendor,
			}

			parts = append(parts, part)
		})

		listPrice := Price{}

		el.ForEach(".tr__total", func(i int, pr *colly.HTMLElement) {
			stringPrice := pr.ChildText(".td__price")
			val, curr, _ := StringPriceToFloat(stringPrice)

			switch pr.ChildText(".td__label") {
			case "Base Total:":
				listPrice.Base = val
			case "Tax:":
				listPrice.Tax = val
			case "Promo Discounts:":
				listPrice.Discounts = val
			case "Shipping:":
				listPrice.Shipping = val
			case "Total:":
				listPrice.Total = val
				listPrice.TotalString = stringPrice
				listPrice.Currency = curr
			}
		})

		compNotes := []CompNote{}

		el.ForEach("#compatibility_notes .info-message", func(i int, note *colly.HTMLElement) {
			mode := note.ChildText("span")
			compNotes = append(compNotes, CompNote{
				Message: strings.TrimLeft(strings.TrimSpace(note.Text), mode),
				Level:   strings.TrimRight(mode, ":"),
			})
		})

		partList = PartList{
			URL:           el.Request.URL.String(),
			Parts:         parts,
			Price:         listPrice,
			Wattage:       strings.TrimPrefix(el.ChildText(".partlist__keyMetric"), "Estimated Wattage:\n"),
			Compatibility: compNotes,
		}
	})
	err := s.Collector.Visit(URL)
	s.Collector.Wait()

	if err != nil {
		return nil, err
	}

	return &partList, nil
}

// Searches for parts using PCPartPicker's search function.
func (s Scraper) SearchParts(searchTerm string, region string) ([]SearchPart, error) {
	if region != "" {
		region += "."
	}

	fullURL := "https://" + region + "pcpartpicker.com/search?q=" + url.QueryEscape(searchTerm)

	if !MatchPCPPURL(fullURL) {
		return nil, errors.New("invalid region")
	}

	searchResults := []SearchPart{}

	var reqURL string

	s.Collector.OnHTML(".pageTitle", func(h *colly.HTMLElement) {
		reqURL = h.Request.URL.String()
	})

	s.Collector.OnHTML(".search-results__pageContent .block", func(el *colly.HTMLElement) {
		el.ForEach(".list-unstyled li", func(i int, searchResult *colly.HTMLElement) {
			vendorURL := searchResult.ChildAttr(".search_results--price a", "href")
			stringPrice := searchResult.ChildText(".search_results--price a")

			price, curr, _ := StringPriceToFloat(stringPrice)

			var vendorName string

			if len(stringPrice) > 0 {
				vendorName = ExtractVendorName(vendorURL)
			}

			partVendor := Vendor{
				URL:  "https://" + el.Request.URL.Host + vendorURL,
				Name: vendorName,
				Price: Price{
					Total:       price,
					TotalString: stringPrice,
					Currency:    curr,
				},
				InStock: len(stringPrice) > 0,
			}

			searchResults = append(searchResults, SearchPart{
				Name:   searchResult.ChildText(".search_results--link a"),
				Image:  "https:" + searchResult.ChildAttr(".search_results--img a img", "src"),
				URL:    "https://" + el.Request.URL.Host + searchResult.ChildAttr(".search_results--link a", "href"),
				Vendor: partVendor,
			})
		})
	})

	err := s.Collector.Visit(fullURL)
	s.Collector.Wait()

	if err != nil {
		return nil, err
	}

	if MatchProductURL(reqURL) {
		return nil, &RedirectError{
			URL: reqURL,
		}
	}

	return searchResults, nil
}

// Fetches a PCPartPicker product via a URL.
func (s Scraper) GetPart(URL string) (*Part, error) {
	if !MatchProductURL(URL) {
		return nil, errors.New("invalid part URL")
	}

	var images []string

	s.Collector.OnHTML(".single_image_gallery_box", func(image *colly.HTMLElement) {
		images = append(images, "https:"+image.ChildAttr("a img", "src"))
	})

	if len(images) < 1 {
		s.Collector.OnHTML("script", func(script *colly.HTMLElement) {
			r := regexp2.MustCompile(`(?<=src:\s").*(?=")`, 0)

			for _, match := range regexp2FindAllString(r, script.Text) {
				if strings.HasPrefix(match, "//") {
					match = "https:" + match
				}
				images = append(images, match)
			}
		})
	}

	rating := Rating{}
	var name string

	s.Collector.OnHTML(".wrapper__pageTitle section.xs-col-11", func(ratingContainer *colly.HTMLElement) {
		var stars uint
		ratingContainer.ForEach(".product--rating li", func(i int, _ *colly.HTMLElement) {
			stars += 1
		})

		rating.Stars = stars
		name = ratingContainer.ChildText(".pageTitle")

		splitParts := strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(ratingContainer.Text, name, ""), ratingContainer.ChildText(".breadcrumb"), "")), ",")

		if len(splitParts) < 2 {
			return
		}

		countParse, _ := strconv.Atoi(strings.Trim(strings.ReplaceAll(splitParts[0], "Ratings", ""), "( "))
		rating.Count = uint(countParse)

		averageParse, _ := strconv.ParseFloat(strings.Trim(strings.ReplaceAll(splitParts[1], "Average", ""), ") "), 64)
		rating.Average = float64(averageParse)
	})

	var vendors []Vendor

	s.Collector.OnHTML("#prices table tbody tr", func(vendor *colly.HTMLElement) {
		if vendor.Attr("class") != "" {
			return
		}

		price := Price{}

		for k, v := range partClassMappings {
			stringPrice := vendor.ChildText(v)
			val, curr, _ := StringPriceToFloat(stringPrice)

			switch k {
			case "Base":
				price.Base = val
			case "Shipping":
				price.Shipping = val
			case "Tax":
				price.Tax = val
			case "Discounts":
				price.Discounts = val
			case "Total":
				price.Total = val
				price.Currency = curr
				price.TotalString = stringPrice
			}
		}

		vendors = append(vendors, Vendor{
			Name:    vendor.ChildAttr(".td__logo a img", "alt"),
			Image:   "https:" + vendor.ChildAttr(".td__logo a img", "src"),
			InStock: vendor.ChildText(".td__availability") == "In stock",
			URL:     "https://" + vendor.Request.URL.Host + vendor.ChildAttr(".td__finalPrice a", "href"),
			Price:   price,
		})
	})

	var specs []Spec

	s.Collector.OnHTML(".specs", func(specsContainer *colly.HTMLElement) {
		if len(specs) > 0 {
			return
		}
		specsContainer.ForEach(".group", func(i int, spec *colly.HTMLElement) {
			var values []string

			spec.ForEach(".group__content li", func(i int, specValue *colly.HTMLElement) {
				values = append(values, specValue.Text)
			})

			if len(values) == 0 {
				values = []string{spec.ChildText(".group__content")}
			}

			specs = append(specs, Spec{
				Name:   spec.ChildText(".group__title"),
				Values: values,
			})
		})

	})

	err := s.Collector.Visit(URL)
	s.Collector.Wait()

	if err != nil {
		return nil, err
	}

	return &Part{
		Name:    name,
		Rating:  rating,
		Specs:   specs,
		Vendors: vendors,
		Images:  images,
	}, nil
}
