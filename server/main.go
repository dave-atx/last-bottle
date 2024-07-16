package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

var supabaseUrl = os.Getenv("SUPABASE_URL")
var supabaseKey = os.Getenv("SUPABASE_KEY")
var apnsKeyID = os.Getenv("APNS_KEY_ID")
var apnsTeamID = os.Getenv("APNS_TEAM_ID")
var apnsAuthKey = os.Getenv("APNS_AUTH_KEY")
var apnsTopic = os.Getenv("APNS_TOPIC")

type Offer struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Varietal    string      `json:"varietal"`
	Vintage     string      `json:"vintage"`
	Price       json.Number `json:"price"`
	Retail      json.Number `json:"retail"`
	BestWeb     json.Number `json:"best_web"`
	Image       string      `json:"image"`
	PurchaseURL string      `json:"purchase_url"`
}

type PushNotificationRegistration struct {
	DeviceToken string `json:"device_token"`
}

type Device struct {
	DeviceToken string `json:"device_token"`
}

func fetchAndParse() {
	fmt.Println("Fetching...")
	client := resty.New()
	resp, err := client.R().Get("https://www.lastbottlewines.com/")
	if err != nil {
		log.Printf("Failed to fetch page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		return
	}

	name := strings.TrimSpace(doc.Find(".offer-name").Text())
	varietal := strings.ReplaceAll(strings.TrimSpace(doc.Find("li:contains('Varietal')").Text()), "Varietal: ", "")
	vintageText := findYear(name)
	vintage := ""
	if vintageText != "" {
		vintage = vintageText + "-01-01"
	}

	price := strings.TrimSpace(doc.Find(".price-holder .amount.lb").First().Text())
	retail := strings.TrimSpace(doc.Find(".price-holder:has(.retail) .amount").First().Text())
	bestWeb := strings.TrimSpace(doc.Find(".price-holder:has(.bestweb) .amount").First().Text())
	image, _ := doc.Find(".offer-image img").Attr("src")
	purchaseURL, _ := doc.Find(".purchase-it a").Attr("href")

	id := extractID(purchaseURL)

	// Check the most recent entry in the database
	recentOffers := []Offer{}
	offersUrl := fmt.Sprintf("%s/offers", supabaseUrl)
	resp, err = client.R().
		SetAuthToken(supabaseKey).
		SetHeader("apikey", supabaseKey).
		SetQueryParam("select", "id").
		SetQueryParam("order", "created_at.desc").
		SetQueryParam("limit", "1").
		SetResult(&recentOffers).
		Get(offersUrl)

	if err != nil {
		log.Printf("unable to select latest offer: %v", err)
		return
	}

	if len(recentOffers) == 0 || recentOffers[0].ID != id {
		newOffer := Offer{
			ID:          id,
			Name:        name,
			Varietal:    varietal,
			Vintage:     vintage,
			Price:       json.Number(price),
			Retail:      json.Number(retail),
			BestWeb:     json.Number(bestWeb),
			Image:       image,
			PurchaseURL: purchaseURL,
		}

		resp, err = client.R().
			SetAuthToken(supabaseKey).
			SetHeader("apikey", supabaseKey).
			SetBody(newOffer).
			Post(offersUrl)

		if err != nil || !resp.IsSuccess() {
			log.Printf("unable to insert order: %v", err)
			return
		}

		log.Printf("Inserted Offer: %+v\n", newOffer)
		// sendPushNotification(newOffer)
	} else {
		log.Println("No new offer to insert")
	}
}

func findYear(text string) string {
	// Extract the year (e.g., 2019) from the text
	year := ""
	for _, word := range strings.Fields(text) {
		if len(word) == 4 && word[0] >= '0' && word[0] <= '9' {
			year = word
			break
		}
	}
	return year
}

func extractID(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2]
}

func sendPushNotification(offer Offer) {
	authKey, err := token.AuthKeyFromBytes([]byte(apnsAuthKey))
	if err != nil {
		log.Fatalf("Failed to parse APNS auth key: %v", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   apnsKeyID,
		TeamID:  apnsTeamID,
	}

	client := apns2.NewTokenClient(token)

	// Fetch device tokens from the database
	var devices []Device
	// _, err = supabase.From("devices").Select("device_token", "", false).ExecuteTo(&devices)
	// if err != nil {
	// 	log.Fatalf("Failed to fetch device tokens: %v", err)
	// }

	for _, device := range devices {
		notification := &apns2.Notification{}
		notification.DeviceToken = device.DeviceToken
		notification.Topic = apnsTopic
		notification.Payload = payload.NewPayload().AlertTitle("New Offer Available!").AlertBody(fmt.Sprintf("%s is now available for $%.2f", offer.Name, offer.Price))

		res, err := client.Push(notification)
		if err != nil {
			log.Printf("Failed to send push notification to %s: %v", device.DeviceToken, err)
		} else {
			fmt.Printf("APNs Response: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		}
	}
}

func getOffers(c *gin.Context) {
	var offers []Offer

	offersUrl := fmt.Sprintf("%s/offers", supabaseUrl)
	resp, err := resty.New().R().
		SetAuthToken(supabaseKey).
		SetHeader("apikey", supabaseKey).
		SetQueryParam("select", "*").
		SetQueryParam("order", "created_at.desc").
		SetResult(&offers).
		Get(offersUrl)

	if err != nil || !resp.IsSuccess() {
		log.Printf("error getting offers: %v", err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusOK, offers)
}

func registerForPushNotifications(c *gin.Context) {
	var registration PushNotificationRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Insert the device token into the devices table
	// _, _, err := supabase.From("devices").Insert(map[string]interface{}{
	// 	"device_token": registration.DeviceToken,
	// }, false, "", "", "").Execute()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Device registered for push notifications"})
}

func healthcheck(c *gin.Context) {
	offersUrl := fmt.Sprintf("%s/offers", supabaseUrl)
	resp, err := resty.New().R().
		SetAuthToken(supabaseKey).
		SetHeader("apikey", supabaseKey).
		SetQueryParam("select", "id").
		SetQueryParam("limit", "1").
		Head(offersUrl)

	if err != nil || !resp.IsSuccess() {
		c.AbortWithStatus(500)
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}

func main() {
	if supabaseUrl == "" || supabaseKey == "" || apnsKeyID == "" || apnsTeamID == "" || apnsAuthKey == "" || apnsTopic == "" {
		log.Fatalf("Environment variables SUPABASE_URL, SUPABASE_KEY, APNS_KEY_ID, APNS_TEAM_ID, APNS_AUTH_KEY, and APNS_TOPIC must be set")
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		fetchAndParse()
		for {
			select {
			case <-ticker.C:
				fetchAndParse()
			}
		}
	}()

	router := gin.Default()
	router.GET("/api/v1/healthcheck", healthcheck)
	router.GET("/api/v1/offers", getOffers)
	router.POST("/api/v1/register", registerForPushNotifications)
	router.Run(":8080")
}
