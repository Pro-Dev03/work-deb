package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type GeocodingService struct {
	client *resty.Client
	apiKey string
}

type GeocodeResult struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	LabelHe     string  `json:"label_he"`      // الاسم بالعبرية
	Language    string  `json:"language"`
	CountryCode string  `json:"country_code"`
	State       string  `json:"state"`
	City        string  `json:"city"`
	CityHe      string  `json:"city_he"`       // المدينة بالعبرية
	Street      string  `json:"street"`
	StreetHe    string  `json:"street_he"`     // الشارع بالعبرية
	HouseNumber string  `json:"house_number"`
	PostalCode  string  `json:"postal_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Type        string  `json:"type"`
}

func NewGeocodingService(apiKey string) *GeocodingService {
	return &GeocodingService{
		client: resty.New().SetTimeout(10 * time.Second),
		apiKey: apiKey,
	}
}

func (s *GeocodingService) Autocomplete(query, language string) ([]GeocodeResult, error) {
	if len(query) < 2 {
		return nil, nil
	}

	log.Printf("🔍 البحث: %s (اللغة: %s)", query, language)

	// Geoapify API - مع دعم العبرية
	url := "https://api.geoapify.com/v1/geocode/autocomplete"
	
	// استخدام اللغة العبرية للبحث
	searchLang := "he"
	if language == "ar" {
		searchLang = "ar"
	} else if language == "en" {
		searchLang = "en"
	} else {
		searchLang = "he" // افتراضي عبري
	}
	
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"text":   query,
			"apiKey": s.apiKey,
			"limit":  "15",
			"lang":   searchLang,
			"filter": "countrycode:il",
		}).
		Get(url)

	if err != nil {
		log.Printf("❌ فشل الاتصال بـ Geoapify: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("❌ خطأ من Geoapify: %s", resp.String())
		return nil, fmt.Errorf("Geoapify error: %s", resp.Status())
	}

	var response struct {
		Features []struct {
			Properties struct {
				PlaceID     string `json:"place_id"`
				Formatted   string `json:"formatted"`
				CountryCode string `json:"country_code"`
				State       string `json:"state"`
				City        string `json:"city"`
				Street      string `json:"street"`
				HouseNumber string `json:"housenumber"`
				Postcode    string `json:"postcode"`
				AddressLine1 string `json:"address_line1"`
				AddressLine2 string `json:"address_line2"`
				ResultType  string `json:"result_type"`
				Lat         float64 `json:"lat"`
				Lon         float64 `json:"lon"`
				// أسماء بالعبرية من المصدر
				Name        string `json:"name"`
			} `json:"properties"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Printf("❌ فشل تحليل الرد: %v", err)
		return nil, err
	}

	var results []GeocodeResult
	for _, feature := range response.Features {
		props := feature.Properties
		
		var lat, lon float64
		if len(feature.Geometry.Coordinates) >= 2 {
			lon = feature.Geometry.Coordinates[0]
			lat = feature.Geometry.Coordinates[1]
		} else {
			lat = props.Lat
			lon = props.Lon
		}

		// استخراج الأسماء بالعبرية
		cityHe := props.City
		streetHe := props.Street
		
		// إذا كان الاسم فارغاً، استخدم المنسق
		labelHe := props.AddressLine1
		if labelHe == "" {
			labelHe = props.Street
			if props.HouseNumber != "" {
				labelHe = labelHe + " " + props.HouseNumber
			}
		}
		if labelHe == "" {
			labelHe = props.City
		}
		if labelHe == "" {
			labelHe = props.Formatted
		}

		results = append(results, GeocodeResult{
			ID:          props.PlaceID,
			Label:       labelHe,
			LabelHe:     labelHe,
			Language:    searchLang,
			CountryCode: props.CountryCode,
			State:       props.State,
			City:        props.City,
			CityHe:      cityHe,
			Street:      props.Street,
			StreetHe:    streetHe,
			HouseNumber: props.HouseNumber,
			PostalCode:  props.Postcode,
			Latitude:    lat,
			Longitude:   lon,
			Type:        props.ResultType,
		})
	}

	log.Printf("✅ تم العثور على %d نتيجة", len(results))
	return results, nil
}
