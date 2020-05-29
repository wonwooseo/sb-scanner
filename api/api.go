package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sb-scanner/services"
	"time"

	"github.com/labstack/echo/v4"
)

var respBytes []byte
var lastUpdate time.Time

// GetLatestHandler _
func GetLatestHandler(c echo.Context) error {
	log.Printf("GET %s: %s", c.Path(), c.RealIP())
	if respBytes == nil || time.Now().Sub(lastUpdate) >= time.Hour {
		cacheCommitList, err := services.SearchLatestCommit()
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "github search failed..")
		}
		respBytes, err = json.Marshal(cacheCommitList)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "github search failed..")
		}
		lastUpdate = time.Now()
	}
	return c.JSONBlob(http.StatusOK, respBytes)
}
