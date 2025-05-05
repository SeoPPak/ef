package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"ef/gemini"
)

type CategorizeRequest struct {
	PayerName string `json:"payerName"`
	Amount    string `json:"amount"`
}

type CategorizeResponse struct {
	Category string  `json:"category"`
	Figure   float64 `json:"figure"`
}

func CategorizeHandler(c *gin.Context) {
	var req CategorizeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	prompt := fmt.Sprintf(`{ "입금자명": %s, "거래 금액": %s } 인 거래내역의 소비 카테고리와 환경발자국 수치를 알려줘. 아래 json 형식으로만 답해줘. { "category": string, "figure" : float }`, 
		req.PayerName, req.Amount)

	geminiResponse, err := gemini.GetGeminiResponse(c.Request.Context(), prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get response from Gemini: %v", err)})
		return
	}

	jsonStr, err := extractJSON(geminiResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to extract JSON from response: %v", err)})
		return
	}

	var response CategorizeResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse JSON response: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
}

func extractJSON(text string) (string, error) {
	re := regexp.MustCompile(`\{[^{}]*"category"[^{}]*"figure"[^{}]*\}`)
	match := re.FindString(text)
	
	if match == "" {
		return "", fmt.Errorf("no valid JSON found in response")
	}
	
	var obj interface{}
	if err := json.Unmarshal([]byte(match), &obj); err != nil {
		return "", fmt.Errorf("found text is not valid JSON: %v", err)
	}
	
	return match, nil
}