package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)


/*create table


*/
type Advert struct {
	ID 			int `json:"ID,omitempty"`
	Title 		string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Price		float32 `json:"price,omitempty"`
	PhotoLink 	[]string `json:"photo_link,omitempty"`

}


func (advert *Advert) Validate() error {
	if len(advert.Title) > 200 {
		return errors.New("Title no more than 200 characters")
	}
	if len(advert.Description) > 1000 {
		return errors.New("Description no more than 200 characters")
	}
	if len(advert.PhotoLink) >3 {
		return errors.New("no more than 3 photo links")
	}
	return nil
}

func (a *Advert) getAdvert(db *sql.DB) error{
	 return db.QueryRow(
		 "SELECT Title, Price, Photo_link FROM advert WHERE ID=$1",
		 a.ID).Scan(&a.Title, &a.Price, pq.Array(&a.PhotoLink))

}

type Fields struct {
	onDescription bool
	numberPhotoLink int
}

func (a *Advert) getAdvertAdditionalFields(db *sql.DB, field Fields )  error{
	sqlQuery := "SELECT Title, Price, %s Photo_link[1:%d] FROM advert WHERE ID=$1"
	if field.onDescription {
		sqlQuery := fmt.Sprintf(sqlQuery, "Description,", field.numberPhotoLink)
		return db.QueryRow(
			sqlQuery, a.ID).Scan(&a.Title, &a.Price, &a.Description,pq.Array(&a.PhotoLink))
	} else {
		sqlQuery := fmt.Sprintf(sqlQuery, "", field.numberPhotoLink)
		return db.QueryRow(
			sqlQuery,
			a.ID).Scan(&a.Title, &a.Price, pq.Array(&a.PhotoLink))
	}
}

type SortType int
const(
	None SortType = iota
	Desc
	Ask
)

type ISort interface {
	getSqlString() string
}

func (s SortType)String() string  {
	return [...] string{"None", "Desc", "Ask"}[s]
}

type Sort struct {
	Column string
	Type SortType
}

func (s *Sort)getSqlString() string{
	sort := "ORDER BY "
	if s.Type == None && s.Column != ""{
		return sort + s.Column
	}
	switch s.Type{
	case None:
		return ""
	case Ask:
		sort = sort + s.Column + " ASC"
	case Desc:
		sort = sort + s.Column + " DESC"
	}

	return sort
}


func getAllAdverts(db *sql.DB, start, count int, sort ISort ) ([]Advert, error) {

	sqlQuery := fmt.Sprintf("SELECT Title, Price, Photo_link[1:1] FROM Advert %s LIMIT $1 OFFSET $2", sort.getSqlString())

	rows, err := db.Query(sqlQuery, count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	adverts := []Advert{}

	for rows.Next() {
		var a Advert

		if err := rows.Scan(&a.Title, &a.Price, pq.Array(&a.PhotoLink)); err != nil {
			return nil, err
		}
		adverts = append(adverts, a)
	}

	return adverts, nil
}

func (a *Advert) createAdvert(db *sql.DB) error{
	err := db.QueryRow(
		"INSERT INTO advert (title, description, price, photo_link, time_create) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		a.Title,a.Description,a.Price, pq.Array(a.PhotoLink), time.Now()).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}


