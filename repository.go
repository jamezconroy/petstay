package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// Get a single post
func retrieve(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

func retrievePet(id int) (pet Pet, err error) {
	pet = Pet{}
	err = Db.QueryRow("select id, name, owner from pet where id = $1", id).Scan(&pet.Id, &pet.Name, &pet.Owner)
	return
}

func retrievePets(offset int, limit int) (pets []Pet, err error) {

	rows, err := Db.Query("select id, name, owner from pet offset $1 limit $2", offset, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var pet Pet
		err := rows.Scan(&pet.Id, &pet.Name, &pet.Owner)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(pet.Id, pet.Name, pet.Owner)
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// Create a new post
func (post *Post) create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (pet *Pet) create() (err error) {
	statement := "insert into pet (name, owner) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pet.Name, pet.Owner).Scan(&pet.Id)
	return
}

// Update a post
func (post *Post) update() (err error) {
	_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (pet *Pet) updatePet() (err error) {
	_, err = Db.Exec("update pet set name = $2, owner = $3 where id = $1", pet.Id, pet.Name, pet.Owner)
	return
}


// Delete a post
func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}

func (pet *Pet) deletePet() (err error) {
	_, err = Db.Exec("delete from pet where id = $1", pet.Id)
	return
}
