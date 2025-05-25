package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
	"untitled/internal/database"
)

func (apiCfg *apiConfig) handlerFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	feed, err := apiCfg.DB.GetFeedByURL(r.Context(), params.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Feed not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	feed_follows, err := apiCfg.DB.CreateFollowsFeed(r.Context(), database.CreateFollowsFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feeds, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsArrToFeedFollowsArr(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	fmt.Printf("user_id: %v feed_id: %v", user.ID, params.FeedID)
	feedId, err := uuid.Parse(params.FeedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "bad feed id asshole")
		return
	}

	apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		FeedID: feedId,
		UserID: user.ID,
	})
	respondWithJSON(w, http.StatusOK, "ok")
}
