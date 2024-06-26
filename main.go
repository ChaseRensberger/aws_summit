package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const customTimeLayout = "2006-01-02T15:04:05-0700"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}
	t, err := time.Parse(customTimeLayout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type AdditionalFields struct {
	Duration    string `json:"duration"`
	SubHeadline string `json:"subHeadline"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Expertise   string `json:"expertise"`
	ContentType string `json:"contentType"`
	Headline    string `json:"headline"`
	EventDate   string `json:"eventDate"`
}

type Item struct {
	ID               string           `json:"id"`
	Locale           string           `json:"locale"`
	DirectoryID      string           `json:"directoryId"`
	Name             string           `json:"name"`
	Author           string           `json:"author"`
	CreatedBy        string           `json:"createdBy"`
	LastUpdatedBy    string           `json:"lastUpdatedBy"`
	DateCreated      CustomTime       `json:"dateCreated"`
	DateUpdated      CustomTime       `json:"dateUpdated"`
	AdditionalFields AdditionalFields `json:"additionalFields"`
}

type Tag struct {
	ID             string     `json:"id"`
	Locale         string     `json:"locale"`
	TagNamespaceID string     `json:"tagNamespaceId"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CreatedBy      string     `json:"createdBy"`
	LastUpdatedBy  string     `json:"lastUpdatedBy"`
	DateCreated    CustomTime `json:"dateCreated"`
	DateUpdated    CustomTime `json:"dateUpdated"`
}

type ItemWithTags struct {
	Item Item  `json:"item"`
	Tags []Tag `json:"tags"`
}

type responseType struct {
	Items []ItemWithTags `json:"items"`
}

func main() {
	urlString := "https://aws.amazon.com/api/dirs/items/search?item.directoryId=amer-summit&size=302&item.locale=en_US&tags.id=amer-summit%23location%23washington-dc&tags.id=amer-summit%23day%232024-06-26&tags.id=amer-summit%23day%232024-06-27&page=0"

	response, err := http.Get(urlString)
	if err != nil {
		fmt.Println("Failed to query URL:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var responseObj responseType
	var items []ItemWithTags
	err = json.Unmarshal(body, &responseObj)
	items = responseObj.Items
	if err != nil {
		fmt.Println("Failed to unmarshal response:", err)
		return
	}

	fmt.Println("Count:", len(items))

	var filteredItems []ItemWithTags
	for _, item := range items {
		if item.Item.AdditionalFields.EventDate == "June 26th" && item.Item.AdditionalFields.Expertise == "300 – Advanced" {
			filteredItems = append(filteredItems, item)
		}
	}

	fmt.Println("Filtered Count:", len(filteredItems))

	for _, item := range filteredItems {
		fmt.Println("----------------------------------------")
		fmt.Println("Name:", item.Item.Name)
		fmt.Println("Duration:", item.Item.AdditionalFields.Duration)
		fmt.Println("Teacher:", item.Item.AdditionalFields.SubHeadline)
		fmt.Println("Description:", item.Item.AdditionalFields.Description)
		fmt.Println("Category:", item.Item.AdditionalFields.Category)
		fmt.Println("Expertise:", item.Item.AdditionalFields.Expertise)
		fmt.Println("Headline:", item.Item.AdditionalFields.Headline)
		fmt.Println("Event Date:", item.Item.AdditionalFields.EventDate)
	}

}
