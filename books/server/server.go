package main

import (
	v1 "books/pb/v1"
	"context"
	"errors"
	"fmt"
	"net"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	util "books/util"
)

type Item interface {
	key() string
	val() string
}

type BooksService struct {
	v1.UnimplementedBooksServiceServer
	books []*v1.AddRequest
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}

	S := grpc.NewServer()
	ss := &BooksService{}
	v1.RegisterBooksServiceServer(S, ss)

	reflection.Register(S)

	if e := S.Serve(listener); e != nil {
		panic(e)
	}

}

func (S *BooksService) AddBook(ctx context.Context, Addrequest *v1.AddRequest) (*v1.AddResponse, error) {
	for _, b := range S.books {
		if Addrequest.ID == b.ID {
			return &v1.AddResponse{}, errors.New(("ID already exists"))
		}
	}
	S.books = append(S.books, Addrequest)
	for _, v := range S.books {
		if v.ID != Addrequest.ID {
			return &v1.AddResponse{
				ID:     Addrequest.ID,
				Title:  Addrequest.Title,
				Author: Addrequest.Author,
			}, nil
		}
	}
	return &v1.AddResponse{
		ID:     Addrequest.ID,
		Title:  Addrequest.Title,
		Author: Addrequest.Author}, nil
}

func (S *BooksService) UpdateBook(ctx context.Context, Addrequest *v1.AddRequest) (*v1.UpdateResponse, error) {

	for _, v := range S.books {
		if v.ID == Addrequest.ID {
			v.ID = Addrequest.ID
			v.Title = Addrequest.Title
			v.Author = Addrequest.Author
			return &v1.UpdateResponse{
				ID:     Addrequest.ID,
				Title:  Addrequest.Title,
				Author: Addrequest.Author,
			}, nil
		}
	}
	return &v1.UpdateResponse{}, nil
}

func (S *BooksService) GetAllBooks(ctx context.Context, _ *emptypb.Empty) (*v1.GetAllBooksResponse, error) {

	out1 := make([]*v1.AddResponse, 0)
	for _, v := range S.books {
		a := v1.AddResponse{
			ID:     v.ID,
			Title:  v.Title,
			Author: v.Author,
		}
		out1 = append(out1, &a)
	}

	out2 := &v1.GetAllBooksResponse{
		Books: out1,
	}
	return out2, nil
}

func (S *BooksService) GetBookByID(ctx context.Context, IDrequest *v1.IDRequest) (*v1.IDResponse, error) {

	for _, book := range S.books {
		if IDrequest.ID == book.ID {
			return &v1.IDResponse{
				ID:     book.ID,
				Title:  book.Title,
				Author: book.Author,
			}, nil
		}
	}
	return &v1.IDResponse{}, errors.New(("ID doesn't exists"))
}

func (S *BooksService) GetBookByAuthor(ctx context.Context, Authorrequest *v1.AuthorRequest) (*v1.AuthorResponse, error) {

	for _, book := range S.books {
		if Authorrequest.Author == book.Author {
			return &v1.AuthorResponse{
				ID:     book.ID,
				Title:  book.Title,
				Author: book.Author,
			}, nil
		}
	}
	return &v1.AuthorResponse{}, errors.New(("author doesn't exists"))
}

func (S *BooksService) GetBookByTitle(ctx context.Context, Titlerequest *v1.TitleRequest) (*v1.TitleResponse, error) {

	for _, book := range S.books {
		if Titlerequest.Title == book.Title {
			return &v1.TitleResponse{
				ID:     book.ID,
				Title:  book.Title,
				Author: book.Author,
			}, nil
		}
	}

	return &v1.TitleResponse{}, errors.New(("title Doesn't exists"))
}

func (S *BooksService) Arbitrary(ctx context.Context, Datarequest *v1.DataRequest) (*v1.DataResponse, error) {

	mapper := Datarequest.ArbitraryData.AsMap()

	for k, v := range mapper {
		if k == "brand" && v == "Subaru" {
			mapper["waiver"] = true
			mapper["country"] = "america"
		} else {
			continue
		}
	}

	out, err := util.MaptoStruct(mapper)
	if err != nil {
		fmt.Println("Cannot Map to Struct")
	}

	return &v1.DataResponse{ArbitraryData: out}, nil
}

func (S *BooksService) DataArbitrary(ctx context.Context, request *structpb.Struct) (*structpb.Struct, error) {

	var response *structpb.Struct

	mapper := request.AsMap()

	for k := range mapper {
		if k == "ArbitraryData" {
			a := mapper["ArbitraryData"]
			b, ok := a.(map[string]interface{})
			if ok {
				b["familySharing"] = false
				b["country"] = "Africa"
			} else {
				continue
			}
		}
	}

	out, err := util.MaptoStruct(mapper)
	if err != nil {
		fmt.Println("Cannot Map to Struct")
	}
	response = out
	return response, nil
}
