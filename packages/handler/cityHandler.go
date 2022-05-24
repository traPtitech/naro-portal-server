package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int    `json:"id,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  string `json:"population,omitempty" db:"Population"`
}
type Country struct {
	Code        string  `json:"code,omitempty" db:"Code"`
	Name        string  `json:"name,omitempty" db:"Name"`
	Continent   string  `json:"continent,omitempty" db:"Continent"`
	Region      string  `json:"region,omitempty" db:"Region"`
	SurfaceArea float32 `json:"surfaceArea,omitempty" db:"SurfaceArea"`
	// IndepYear      sql.NullInt16   `json:"indepYear,omitempty" db:"IndepYear"`
	Population int `json:"population,omitempty" db:"Population"`
	// LifeExpectancy uint8           `json:"lifeExpectancy,omitempty" db:"LifeExpectancy"`
	// GNP            sql.NullFloat64 `json:"gnp,omitempty" db:"GNP"`
	// GNPOld         sql.NullFloat64 `json:"gnpOld,omitempty" db:"GNPOld"`
	// LocalName      sql.NullString  `json:"localName,omitempty" db:"LocalName"`
	// GovernmentForm sql.NullString  `json:"governmentForm,omitempty" db:"GovernmentForm"`
	// HeadOfState    sql.NullString  `json:"headOfState,omitempty" db:"HeadOfState"`
	Capital sql.NullInt64 `json:"capital,omitempty" db:"Capital"`
	// Code2          sql.NullString  `json:"code2,omitempty" db:"Code2"`
}

func AddCityHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := &City{}
		err := c.Bind(data)
		if err != nil { // エラーが発生した際
			// fmt.Sprintf("%+v", data): dataをstringに変換
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}
		cityState := `INSERT INTO city(ID, Name, CountryCode, District, Population) VALUES (?,?,?,?,?)`
		_, err2 := db.Exec(cityState, data.ID, data.Name, data.CountryCode, data.District, data.Population)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err2.Error()))
		}

		cityID := data.ID
		fmt.Println(cityID)
		var city City
		if err := db.Get(&city, "SELECT * FROM city WHERE ID=?", cityID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City ID=%d", cityID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.JSON(http.StatusOK, city)
	}
}

func GetCityInfoHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cityName := c.Param("cityName")
		cities := []City{}
		if err := db.Select(&cities, "SELECT * FROM city WHERE Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City Name=%s", cityName)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.JSON(http.StatusOK, cities)
	}
}

func GetCountryListHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		countries := []Country{}
		if err := db.Select(&countries, "SELECT Code, Name, Continent, Region, SurfaceArea, Population, Capital FROM country"); errors.Is(err, sql.ErrNoRows) {
			log.Printf("Something Error")
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.JSON(http.StatusOK, countries)
	}
}
func GetCityListInCountryHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		countryConde := c.Param("countryCode")
		cities := []City{}
		if err := db.Select(&cities, "SELECT * FROM city WHERE CountryCode=?", countryConde); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City Name=%s", countryConde)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}
		return c.JSON(http.StatusOK, cities)
	}
}
func DeleteCityInfoHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cityID := c.Param("cityID")
		fmt.Println(cityID)
		i, atoiErr := strconv.Atoi(cityID)
		if atoiErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("IDをURLの末尾に入れてください"))
		}
		if _, err := db.Exec("DELETE FROM city WHERE ID=?", i); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City Name=%s", cityID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.String(http.StatusOK, "Delete City ID:"+cityID)
	}
}
