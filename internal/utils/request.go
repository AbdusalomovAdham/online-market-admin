package utils

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetQuery(c *gin.Context, query string) (*string, error) {
	q := c.Query(query)

	if q != "" {
		err := checkSQLInjection(q)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": err.Error(),
				"status":  false,
			})
			return nil, err
		}

		return &q, nil
	}

	return nil, nil
}

func checkSQLInjection(text string) error {
	containsOr := strings.Contains(text, " or ")
	if containsOr {
		return errors.New("sql injection detected!")
	}
	containsDropTable := strings.Contains(text, "drop table")
	if containsDropTable {
		return errors.New("sql injection detected!")
	}
	containsDropDatabase := strings.Contains(text, "drop database")
	if containsDropDatabase {
		return errors.New("sql injection detected!")
	}
	containsSelectTable := strings.Contains(text, "select ")
	if containsSelectTable {
		return errors.New("sql injection detected!")
	}
	containsComment := strings.Contains(text, " -- ")
	if containsComment {
		return errors.New("sql injection detected!")
	}
	match, _ := regexp.MatchString("([^a-zA-Z0-9]+)=([^a-zA-Z0-9]+)", text)
	if match {
		return errors.New("sql injection detected!")
	}
	return nil
}
