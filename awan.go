package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "path/to/your/docs" // Import the docs generated by swaggo
	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var jwtKey = []byte("awanpay")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TopUp struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type Earnings struct {
	TotalIncome float64 `json:"total_income"`
	TotalExpenditure float64 `json:"total_expenditure"`
}

type FinancialStatus struct {
	TotalBalance float64 `json:"total_balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpenditure float64 `json:"total_expenditure"`
}

// Store the total earnings and expenditure (for demonstration purposes)
var earnings = Earnings{
	TotalIncome:     1000000, // Example initial amount
	TotalExpenditure: 500000,  // Example initial amount
}

// Store Top Up records (for demonstration purposes)
var topUps []TopUp

// Function to generate QRIS (for payment)
func generateQrisHandler(c *gin.Context) {
	qrisData := map[string]interface{}{
		"merchant_name": "AwanPay",
		"amount":        100000, // Example amount
		"currency":      "IDR",
	}
	requestBody, err := json.Marshal(qrisData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QRIS data"})
		return
	}

	apiUrl := "https://api.qris.example.com/generate" // Example API URL for QRIS generation
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YOUR_API_KEY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request QRIS generation"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse QRIS response"})
		return
	}

	qrisImageURL := responseData["qris_image_url"].(string)
	c.JSON(http.StatusOK, gin.H{
		"qris_image_url": qrisImageURL,
	})
}

// Endpoint to Top Up
func topUpHandler(c *gin.Context) {
	var topUpRequest TopUp
	if err := c.ShouldBindJSON(&topUpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Simulate top-up action
	topUps = append(topUps, topUpRequest)
	earnings.TotalIncome += topUpRequest.Amount // Increase total income with top-up

	c.JSON(http.StatusOK, gin.H{
		"message": "Top up successful",
		"total_balance": earnings.TotalIncome,
	})
}

// Endpoint to get earnings and expenditures
func getFinancialStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_balance":     earnings.TotalIncome - earnings.TotalExpenditure,
		"total_income":      earnings.TotalIncome,
		"total_expenditure": earnings.TotalExpenditure,
	})
}

// Endpoint to get total top-ups
func getTotalTopUpsHandler(c *gin.Context) {
	var totalTopUp float64
	for _, topUp := range topUps {
		totalTopUp += topUp.Amount
	}

	c.JSON(http.StatusOK, gin.H{
		"total_top_up": totalTopUp,
	})
}
func loginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if creds.Username != "admin" || creds.Password != "brayy" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "awanbrayy",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
func main() {
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.GET("/generate-qris", generateQrisHandler)
	r.POST("/top-up", topUpHandler)
	r.GET("/financial-status", getFinancialStatusHandler)
	r.GET("/total-top-ups", getTotalTopUpsHandler)
	r.StaticFile("/", "./index.html")
	r.StaticFile("/signup", "./signup.html")
	r.StaticFile("/dashboard", "./dashboard.html")
	r.GET("/swagger/*any", swagger.WrapHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}