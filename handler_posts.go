package main

import (
	"net/http"
	"strconv"

	"github.com/quintui/rssagg/internal/database"
)

func (apiCfg apiConfig) handlerUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	strLimit := r.URL.Query().Get("limit")
	limit := 10

	if customLimit, err := strconv.Atoi(strLimit); err == nil {
		limit = customLimit
	}

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
