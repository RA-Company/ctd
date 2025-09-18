package ctd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type StatisticRating struct {
	ScoreValue         json.Number `json:"score_value"`
	RatingScaleScore   int64       `json:"rating_scale_score"`
	ValuationRequestID int64       `json:"valuation_request_id"`
}

type StatisticRatingsResponse struct {
	Data []StatisticRating `json:"data"`
	Meta MetaResponse      `json:"meta"`
	BasicResponse
}

// GetScoreValue converts the ScoreValue from json.Number to int64.
// If the conversion fails, it returns -1 to indicate an invalid score.
//
// Returns:
//   - An int64 representing the score value, or -1 if the conversion fails.
func (dst *StatisticRating) GetScoreValue() int64 {
	if result, err := dst.ScoreValue.Int64(); err != nil {
		return -1
	} else {
		return result
	}
}

// GetRangeValue categorizes the score value into three ranges based on the provided limits.
// It uses the GetScoreValue method to retrieve the score value.
// If the score value is less than or equal to limit1, it returns 1.
// If the score value is less than or equal to limit2, it returns 2.
// If the score value is greater than limit2, it returns 3.
// If the score value is invalid (i.e., -1), it returns 0.
//
// Parameters:
//   - limit1: The first limit for categorization.
//   - limit2: The second limit for categorization.
//
// Returns:
//   - A uint8 representing the category of the score value (0, 1, 2, or 3).
func (dst *StatisticRating) GetRangeValue(limit1, limit2 int64) uint8 {
	res := dst.GetScoreValue()

	if res == -1 {
		return 0
	}

	if res <= limit1 {
		return 1
	} else if res <= limit2 {
		return 2
	}

	return 3
}

// APIStatisticsRating retrieves a list of statistic ratings from the Chat2Desk API.
// It constructs the API endpoint URL with the provided date, offset, and limit,
// sends a GET request to the API, and returns the response data as a StatisticRatingsResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the StatisticRatingsResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - date: The date for which to retrieve statistics. If zero, the current date is used.
//   - offset: The offset for pagination, indicating where to start fetching ratings.
//   - limit: The maximum number of ratings to return.
//
// Returns:
//   - A pointer to a StatisticRatingsResponse struct containing the list of statistic ratings and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIStatisticsRating(ctx context.Context, date time.Time, offset, limit int) (*StatisticRatingsResponse, error) {
	if date.IsZero() {
		date = time.Now()
	}

	url := fmt.Sprintf("%sv1/statistics?report=rating&date=%s&offset=%d&limit=%d", dst.Url, date.Format("2006-01-02"), offset, limit)
	timeout := dst.Timeout
	defer func() { dst.Timeout = timeout }()

	response := StatisticRatingsResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed to get statistics: %v", err)
		return nil, err
	}

	return &response, nil
}

// StatisticsRating retrieves a list of statistic ratings from the Chat2Desk API.
// It uses the APIStatisticsRating method to fetch the ratings and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a slice of StatisticRating, which contains the ratings.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - date: The date for which to retrieve statistics. If zero, the current date is used.
//   - offset: The offset for pagination, indicating where to start fetching ratings.
//   - limit: The maximum number of ratings to return.
//
// Returns:
//   - A pointer to a slice of StatisticRating containing the list of statistic ratings.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) StatisticsRating(ctx context.Context, date time.Time, offset int, limit int) (*[]StatisticRating, error) {
	data, err := dst.APIStatisticsRating(ctx, date, offset, limit)
	if err != nil {
		return nil, err
	}

	fmt.Printf("StatisticsRating: %+v\n", data)

	if data.Status == "error" {
		dst.Error(ctx, "Failed to get statistics: %s", data.Errors)
		return nil, fmt.Errorf("failed to get statistics: %s", data.Errors)
	}

	return &data.Data, nil
}

// AllStatisticsRating retrieves all statistic ratings from the Chat2Desk API by handling pagination.
// It repeatedly calls the StatisticsRating method with increasing offsets until all ratings are fetched.
// It returns a pointer to a slice of StatisticRating, which contains all the ratings.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - date: The date for which to retrieve statistics. If zero, the current date is used.
//
// Returns:
//   - A pointer to a slice of StatisticRating containing all the statistic ratings.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) AllStatisticsRating(ctx context.Context, date time.Time) (*[]StatisticRating, error) {
	ratings := []StatisticRating{}
	offset := 0
	limit := 200

	for {
		data, err := dst.StatisticsRating(ctx, date, offset, limit)
		if err != nil {
			return nil, err
		}

		if data == nil || len(*data) == 0 {
			break
		}

		ratings = append(ratings, *data...)
		if len(*data) < limit {
			break
		}
		offset += limit
	}

	return &ratings, nil
}
