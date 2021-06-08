package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

type XurLocationMessage struct {
	LocationName string
	PlaceName    string
	SaleItems    map[string]SaleItem
}

type Message struct {
	Screenshot                 string
	Name                       string
	FlavorText                 string
	ItemTypeAndTierDisplayName string
	Hash                       int
}

func GetItemDefinition(body []byte) *Message {
	pj := ItemHashStruct{}

	if err := json.Unmarshal(body, &pj); err != nil {
		log.Println(err)
	}

	itemDefinition := &Message{
		Screenshot:                 pj.Response.Screenshot,
		Name:                       pj.Response.DisplayProperties.Name,
		FlavorText:                 pj.Response.FlavorText,
		ItemTypeAndTierDisplayName: pj.Response.ItemTypeAndTierDisplayName,
	}
	return itemDefinition
}

func (u *UrlStruct) getData(itemHashes []int) []*Message {
	var dataSlice []*Message

	for _, item := range itemHashes {
		hash := strconv.Itoa(item)
		u.Url = fmt.Sprint("https://www.bungie.net/Platform/Destiny2/Manifest/DestinyInventoryItemDefinition/" + hash + "/")
		u.HeaderKey = "X-API-Key"
		u.HeaderValue = os.Getenv("DESTINY_API_KEY")
		data := u.responseGetter()
		m := GetItemDefinition(data)
		m.Hash = item
		dataSlice = append(dataSlice, m)
	}

	return dataSlice
}

func GetXurLocation(body []byte) (*XurLocationMessage, error) {
	pj := XurLocation{}

	if err := json.Unmarshal(body, &pj); err != nil {
		log.Println(err)
	}
	if reflect.ValueOf(pj).IsZero() {
		return nil, errors.New("xur is not here")
	}

	xurPlace := &XurLocationMessage{
		LocationName: pj.LocationName,
		PlaceName:    pj.PlaceName,
	}
	return xurPlace, nil
}

func (u *UrlStruct) GetXurInventory() []int {
	u.Url = "https://www.bungie.net/Platform/Destiny2/Vendors/?components=402"
	u.HeaderKey = "X-API-Key"
	u.HeaderValue = os.Getenv("DESTINY_API_KEY")

	data := u.responseGetter()
	pj := ManifestHashStruct{}
	if err := json.Unmarshal(data, &pj); err != nil {
		log.Println(err)
	}
	var itemHashes []int
	for _, item := range pj.ResponseManifest.Sales.Data.The2190858386.SaleItems {
		validcheck := IsValidHash(item.ItemHash)
		if !validcheck {
			itemHashes = append(itemHashes, item.ItemHash)
		}
	}
	return itemHashes
}

func IsValidHash(category int) bool {
	switch category {
	case
		bannedHashEngram,
		bannedHashQuest:
		return true
	}
	return false
}

func GetXurData() *XurLocationMessage {
	u := &UrlStruct{
		Url:         "https://paracausal.science/xur/current.json",
		HeaderKey:   "Content-Type",
		HeaderValue: "'application/json",
	}
	data := u.responseGetter()
	xurPlace, err := GetXurLocation(data)
	if err != nil {
		log.Println(err)
	}

	return xurPlace
}

func ParseHashesData() []string {
	u := &UrlStruct{}
	itemHashes := u.GetXurInventory()
	data := u.getData(itemHashes)
	var itemsSlice []string
	for _, i := range data {
		text := fmt.Sprintf("<a href=\"https://www.light.gg/db/items/%v\">%v</a> - %v /item_%v",
			i.Hash,
			i.Name,
			i.ItemTypeAndTierDisplayName,
			i.Hash)
		itemsSlice = append(itemsSlice, text)
	}
	return itemsSlice
}

func ParseHashDataOneItem(hash string) (*Message, error) {
	u := &UrlStruct{}
	itemHash, err := strconv.Atoi(hash)
	if err != nil {
		return nil, errors.New("incorrect hash")
	}
	hashslice := []int{itemHash}
	data := u.getData(hashslice)
	result := data[0]
	return result, nil
}

func HashRegexp(s string) string {
	re := regexp.MustCompile(`^item_(\d+)$`)
	regexpResult := re.FindString(s)
	if regexpResult == "" {
		return ""
	} else {
		return regexpResult[5:]
	}
}
