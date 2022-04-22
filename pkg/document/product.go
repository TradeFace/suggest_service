package document

type Product struct {
	IbCode      string         `jsonapi:"primary,product"`
	Description string         `jsonapi:"attr,description"`
	Images      []ProductImage `jsonapi:"attr,images"`
}

type ProductImage struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	DisplayName  string `json:"displayName"`
}
