package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Confirm(w http.ResponseWriter, r *http.Request) {
	log.Debug("Confirm attempt")

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{":26379", ":26380", ":26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetRedisPassword(),
		Password: utils.GetRedisPassword(),
	})

	tokenStr, exists := mux.Vars(r)["token"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	EMAIL_SECRET, exists := os.LookupEnv("EMAIL_KEY")
	if !exists {
		log.Error("EMAIL_SECRET is not set")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.ConfirmClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(EMAIL_SECRET), nil
	})

	if err != nil {
		log.Error("invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*models.ConfirmClaims)
	if !ok || !token.Valid {
		log.Error("invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := rdb.HSet(ctx, claims.Key, "confirm").Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Account Confirmed!")); err != nil {
		log.Error(err)
	}
}
