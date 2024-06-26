package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
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
	Duration      string `json:"duration"`
	SubHeadline   string `json:"subHeadline"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Expertise     string `json:"expertise"`
	ContentType   string `json:"contentType"`
	Headline      string `json:"headline"`
	EventDate     string `json:"eventDate"`
	DurationStart string `json:"durationStart"`
	DurationEnd   string `json:"durationEnd"`
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

type Response struct {
	Items []ItemWithTags `json:"items"`
}
type ExpertiseLevel int

const (
	Beginner ExpertiseLevel = iota
	Intermediate
	Advanced
	Expert
)

func (el ExpertiseLevel) String() string {
	switch el {
	case Beginner:
		return "100 – Beginner"
	case Intermediate:
		return "200 – Intermediate"
	case Advanced:
		return "300 – Advanced"
	case Expert:
		return "400 – Expert"
	default:
		return "Unknown"
	}
}

type Filter struct {
	ExpertiseLevel ExpertiseLevel
	EventDate      string
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

	var responseObj Response
	var items []ItemWithTags
	err = json.Unmarshal(body, &responseObj)
	items = responseObj.Items
	if err != nil {
		fmt.Println("Failed to unmarshal response:", err)
		return
	}

	for i, item := range items {
		duration := strings.Split(item.Item.AdditionalFields.Duration, "-")
		startTime := strings.TrimSpace(duration[0])
		endTime := strings.TrimSpace(duration[1])
		startTimeFormatted, err := time.Parse("03:04 PM", startTime)
		if err != nil {
			fmt.Println("Failed to parse start time:", err)
			return
		}
		endTimeFormatted, err := time.Parse("03:04 PM", endTime)
		if err != nil {
			fmt.Println("Failed to parse end time:", err)
			return
		}
		items[i].Item.AdditionalFields.DurationStart = startTimeFormatted.Format("15:04:05")
		items[i].Item.AdditionalFields.DurationEnd = endTimeFormatted.Format("15:04:05")
	}

	fmt.Println("Count:", len(items))

	activeFilter := Filter{ExpertiseLevel: Expert, EventDate: "June 26th"}

	var filteredItems []ItemWithTags
	for _, item := range items {
		if item.Item.AdditionalFields.EventDate == activeFilter.EventDate && item.Item.AdditionalFields.Expertise == activeFilter.ExpertiseLevel.String() {
			filteredItems = append(filteredItems, item)
		}
	}

	fmt.Println("Filtered Count:", len(filteredItems))

	// Sort filteredItems in ascending order based on duration start
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].Item.AdditionalFields.DurationStart < filteredItems[j].Item.AdditionalFields.DurationStart
	})

	for _, item := range filteredItems {
		fmt.Println("----------------------------------------")
		fmt.Println("Name:", item.Item.Name)
		fmt.Println("Teacher:", item.Item.AdditionalFields.SubHeadline)
		fmt.Println("Description:", item.Item.AdditionalFields.Description)
		fmt.Println("Category:", item.Item.AdditionalFields.Category)
		fmt.Println("Expertise:", item.Item.AdditionalFields.Expertise)
		fmt.Println("Headline:", item.Item.AdditionalFields.Headline)
		fmt.Println("Event Date:", item.Item.AdditionalFields.EventDate)
		fmt.Println("Duration Start:", item.Item.AdditionalFields.DurationStart)
		fmt.Println("Duration End:", item.Item.AdditionalFields.DurationEnd)
	}
}
