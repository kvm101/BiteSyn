package nlp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"restaurant_reviews/internal"
	"strings"
	"time"
)

func CheckMessage(text string) internal.RatingResponse {
	nlp_url := "http://127.0.0.1:8000/rate"

	jsonStr := fmt.Sprintf(`{"text": "%s"}`, text)
	payload := strings.NewReader(jsonStr)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", nlp_url, payload)
	if err != nil {
		return internal.RatingResponse{
			TextReview: "Failed to create request",
			Status:     false,
			Rating:     0,
		}
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return internal.RatingResponse{
				TextReview: "NLP service timeout - request took too long",
				Status:     false,
				Rating:     0,
			}
		}
		return internal.RatingResponse{
			TextReview: "Failed to connect to NLP service",
			Status:     false,
			Rating:     0,
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return internal.RatingResponse{
			TextReview: "Failed to read response",
			Status:     false,
			Rating:     0,
		}
	}

	var rating internal.RatingResponse
	err = json.Unmarshal(body, &rating)
	if err != nil {
		return internal.RatingResponse{
			TextReview: fmt.Sprintf("Failed to parse response: %v", err),
			Status:     false,
			Rating:     0,
		}
	}

	return rating
}
