package common

type ItemHashStruct struct {
	Response Response `json:"Response,omitempty"`
}
type Response struct {
	DisplayProperties          DisplayProperties `json:"displayProperties,omitempty"`
	Screenshot                 string            `json:"screenshot,omitempty"`
	FlavorText                 string            `json:"flavorText,omitempty"`
	ItemTypeAndTierDisplayName string            `json:"itemTypeAndTierDisplayName,omitempty"`
}

type DisplayProperties struct {
	Name string `json:"name,omitempty"`
}
type XurLocation struct {
	LocationName string `json:"locationName"`
	PlaceName    string `json:"placeName"`
}

type ManifestHashStruct struct {
	ResponseManifest ResponseManifest `json:"Response,omitempty"`
}
type ResponseManifest struct {
	Sales Sales `json:"sales,omitempty"`
}
type Sales struct {
	Data Data `json:"data,omitempty"`
}
type Data struct {
	The2190858386 The2190858386 `json:"2190858386,omitempty"`
}
type The2190858386 struct {
	SaleItems map[string]SaleItem `json:"saleItems,omitempty"`
}
type SaleItem struct {
	ItemHash int `json:"itemHash,omitempty"`
}
