package main

import "C"
import (
	context "context"
	"flag"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/logger"
	lru "github.com/hashicorp/golang-lru"
	"google.golang.org/grpc"
	"log"
	"net"
	"qnetwork.net/fbcrawl/fbcrawl"
	"qnetwork.net/fbcrawl/fbcrawl/pb"
	"unsafe"
)

const logPath = "parse.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")
var email = flag.String("email", "change_me@gmail.com", "facebook email")
var password = flag.String("password", "change_me", "facebook password")
var otp = flag.String("otp", "123456", "facebook otp")
var groupId = flag.String("groupId", "334294967318328", "facebook group id, default is 334294967318328")

var allInstances, _ = lru.New(20)

func getColly(pointer *pb.Pointer) *fbcolly.Fbcolly {
	c, ok := allInstances.Get(pointer.Address)
	if ok {
		return c.(*fbcolly.Fbcolly)
	}
	return nil
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.GrpcServer
}

func (s server) Init(ctx context.Context, empty *empty.Empty) (*pb.Pointer, error) {
	instance := fbcolly.New()
	ptr := (uintptr)(unsafe.Pointer(instance))
	allInstances.Add(int64(ptr), instance)
	return &pb.Pointer{Address: int64(ptr)}, nil
}

func (s server) FreeColly(ctx context.Context, pointer *pb.Pointer) (*empty.Empty, error) {
	logger.Info("FreeColly")
	allInstances.Remove(pointer.Address)
	return &empty.Empty{}, nil
}

func (s server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	p := getColly(request.Pointer)

	cookies, err := p.Login(request.Email, request.Password, "")
	if err == nil {
		return &pb.LoginResponse{Cookies: cookies}, err
	}
	return nil, err
}

func (s server) LoginWithCookies(ctx context.Context, request *pb.LoginWithCookiesRequest) (*empty.Empty, error) {
	p := getColly(request.Pointer)
	err := p.LoginWithCookies(request.Cookies)
	return &empty.Empty{}, err
}

func (s server) FetchGroupInfo(ctx context.Context, request *pb.FetchGroupInfoRequest) (*pb.FacebookGroup, error) {
	p := getColly(request.Pointer)
	err, groupInfo := p.FetchGroupInfo(request.GroupUsername)
	return groupInfo, err
}

func (s server) FetchGroupFeed(ctx context.Context, request *pb.FetchGroupFeedRequest) (*pb.FacebookPostList, error) {
	p := getColly(request.Pointer)
	err, postsList := p.FetchGroupFeed(request.GroupId, request.NextCursor)
	return postsList, err
}

func (s server) FetchPost(ctx context.Context, request *pb.FetchPostRequest) (*pb.FacebookPost, error) {
	p := getColly(request.Pointer)
	err, post := p.FetchPost(request.GroupId, request.PostId, request.CommentNextCursor)
	return post, err
}

func (s server) FetchContentImages(ctx context.Context, request *pb.FetchContentImagesRequest) (*pb.FacebookImageList, error) {
	p := getColly(request.Pointer)
	err, imageList := p.FetchContentImages(request.PostId, request.NextCursor)
	return imageList, err
}

func (s server) FetchImageUrl(ctx context.Context, request *pb.FetchImageUrlRequest) (*pb.FacebookImage, error) {
	p := getColly(request.Pointer)
	err, image := p.FetchImageUrl(request.ImageId)
	return image, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGrpcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
