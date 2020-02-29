package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
)

type Book struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	json.NewEncoder(resp).Encode(books)
}

func getBookInfo(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	params := mux.Vars(r)
	var result Book
	for _, value := range books {
		if value.Id == params["id"] {
			result = value
		}
	}
	if (result == Book{}) {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(resp).Encode(result)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	var book Book

	json.NewDecoder(r.Body).Decode(&book)
	if book.Id != "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	book.Id = fmt.Sprint(rand.Intn(10000000))

	books = append(books, book)
	json.NewEncoder(resp).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	params := mux.Vars(r)
	for i, existingBook := range books {
		if existingBook.Id == params["id"] {
			var updatedBook Book
			json.NewDecoder(r.Body).Decode(updatedBook)
			updatedBook.Id = params["id"]
			books[i] = updatedBook
			json.NewEncoder(resp).Encode(updatedBook)
			return
		}
	}
	resp.WriteHeader(http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	params := mux.Vars(r)
	var requestedId Book
	for i, existingBook := range books {
		if existingBook.Id == params["id"] {
			json.NewEncoder(resp).Encode(existingBook)
			books = append(books[:i], books[i+1:]...)
			break
			resp.WriteHeader(http.StatusOK)
		}
	}
	if (requestedId == Book{}) {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
}

func checkBook(w http.ResponseWriter, r *http.Request) {
	resp := prepareResponseHeaders(w)
	params := mux.Vars(r)
	var result Book
	for _, value := range books {
		if value.Id == params["id"] {
			result = value
			json.NewEncoder(resp).Encode(value) //НЕ ПЕРЕДАЁТ
		}
	}
	if (result == Book{}) {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(resp).Encode(result)
}

func prepareResponseHeaders (w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}

func seed(books *[]Book) {
	*books = append(*books, Book{Id: "1", Title: "Евгений Онегин", Description: "Книга о любви", Author: &Author{FirstName: "Александр", LastName: "Пушкин"}})
	*books = append(*books, Book{Id: "2", Title: "Гарри Поттер", Description: "Книга о волшебнике", Author: &Author{FirstName: "Джоан", LastName: "Rowling"}})
	*books = append(*books, Book{Id: "3", Title: "Темная Башня", Description: "Книга о башне", Author: &Author{FirstName: "Стивен", LastName: "Кинг"}})

}


func main() {
	r := mux.NewRouter()
	seed(&books)
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBookInfo).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	r.HandleFunc("/books/{id}", checkBook).Methods("HEAD")

	log.Fatal(http.ListenAndServe(":8080", r))
}
