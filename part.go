package gopartpicker

// The seller of a part.
type Vendor struct {
	// The name of the vendor.
	Name string
	// The full URL for the vendor's image.
	Image string
	// Whether the vendor has the part in stock or not.
	InStock bool
	// The vendor's price for the part.
	Price Price
	// The full, affiliate URL to buy the part from.
	URL string
}

// A part from PCPartPicker search results. It has a lot of data missing, such as full pricing information.
type SearchPart struct {
	// The part's name.
	Name string
	// The full URL for the part's image.
	Image string
	// The full URL for the part's PCPartPicker page.
	URL string
	// The vendor for the part that is offering it for the lowest price.
	Vendor Vendor
}

type Rating struct {
	// The amount of stars the part has on average, out of 5.
	Stars uint
	// The amount of ratings the part has.
	Count uint
	// The average rating of the part.
	Average float64
}

type Spec struct {
	// The name of the spec.
	Name string
	// A slice containing the values of the spec.
	Values []string
}

// A full part object with multiple vendors and reviews.
type Part struct {
	// The part's type, e.g. CPU, Motherboard.
	Type string
	// The part's name.
	Name string
	// A slice of full image URLs.
	Images []string
	// The full URL for the part's PCPartPicker page.
	URL string
	// A slice of vendors supplying the part.
	Vendors []Vendor
	// The specifications of the part.
	Specs []Spec
	// The part's rating.
	Rating Rating
}

// A part from a part list. It has some data missing, such as all available vendors and specs.
type PartListPart struct {
	// The part's type, e.g. CPU, Motherboard.
	Type string
	// The part's name.
	Name string
	// The full URL for the part's image.
	Image string
	// The full URL for the part's PCPartPicker page.
	URL string
	// The vendor for the part that is offering it for the lowest price.
	Vendor Vendor
}

// A compatibility note for a part list.
type CompNote struct {
	// The compatibility note's message.
	Message string
	// The level or severity of the note, e.g. Warning, Note.
	Level string
}

// Represents a PCPartPicker part list.
type PartList struct {
	// The full URL for the part list.
	URL string
	// The parts in the part list.
	Parts []PartListPart
	// The price of all the part list's parts combined.
	Price Price
	// The total estimated wattage for the part list.
	Wattage string
	// Compatibility notes for the part list.
	Compatibility []CompNote
}
