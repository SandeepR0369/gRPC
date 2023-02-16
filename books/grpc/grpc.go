package grpc

import v1 "books/pb/v1"

type BooksService struct {
	v1.UnimplementedBooksServiceServer
}
