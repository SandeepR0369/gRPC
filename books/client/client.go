package main

import (
	v1 "books/pb/v1"
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewBooksServiceClient(conn)

	book := &v1.AddRequest{
		ID:     13,
		Title:  "Power of Subconcious Mind",
		Author: "Joseph Murphy",
	}

	id := &v1.IDRequest{
		ID: 11,
	}

	addBook(client, book)
	getBookByID(client, id)

}

func addBook(client v1.BooksServiceClient, add *v1.AddRequest) {
	resp, err := client.AddBook(context.Background(), add)

	if err != nil {
		log.Fatalf("Cannot add book: %v", err)
	} else {
		log.Printf("A new book has been added with id: %d", resp.ID)
	}

}

func getBookByID(client v1.BooksServiceClient, id *v1.IDRequest) {
	resp1, err := client.GetBookByID(context.Background(), id)
	if err != nil {
		log.Fatalf("Cannot get book by id: %v", err)
	} else {
		log.Printf("A book has been found with id: %d", resp1.ID)
	}

}
