package main

import "C"
import (
	context "context"
	"flag"
	"github.com/google/logger"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"qnetwork.net/fbcrawl/fbcrawl"
	"qnetwork.net/fbcrawl/fbcrawl/pb"
)

const logPath = "parse.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")
var email = flag.String("email", "change_me@gmail.com", "facebook email")
var password = flag.String("password", "change_me", "facebook password")
var otp = flag.String("otp", "123456", "facebook otp")
var groupId = flag.String("groupId", "334294967318328", "facebook group id, default is 334294967318328")

func getColly(context *pb.Context) *fbcolly.Fbcolly {
	instance := fbcolly.New()
	if context != nil && len(context.Cookies) > 0 {
		_ = instance.LoginWithCookies(context.Cookies)
	}

	return instance
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.GrpcServer
}

func (s server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	p := getColly(nil)

	cookies, err := p.Login(request.Email, request.Password, request.TotpSecret)
	if err == nil {
		return &pb.LoginResponse{Cookies: cookies}, err
	}
	return nil, err
}

func (s server) FetchMyGroups(ctx context.Context, request *pb.FetchMyGroupsRequest) (*pb.FacebookGroupList, error) {
	p := getColly(request.Context)
	err, groupInfo := p.FetchMyGroups()
	return groupInfo, err
}

func (s server) FetchGroupInfo(ctx context.Context, request *pb.FetchGroupInfoRequest) (*pb.FacebookGroup, error) {
	p := getColly(request.Context)
	err, groupInfo := p.FetchGroupInfo(request.GroupUsername)
	return groupInfo, err
}

func (s server) FetchUserInfo(ctx context.Context, request *pb.FetchUserInfoRequest) (*pb.FacebookUser, error) {
	p := getColly(request.Context)
	err, userInfo := p.FetchUserInfo(request.Username)
	return userInfo, err
}

func (s server) FetchGroupFeed(ctx context.Context, request *pb.FetchGroupFeedRequest) (*pb.FacebookPostList, error) {
	p := getColly(request.Context)
	err, postsList := p.FetchGroupFeed(request.GroupId, request.NextCursor)
	return postsList, err
}

func (s server) FetchPost(ctx context.Context, request *pb.FetchPostRequest) (*pb.FacebookPost, error) {
	p := getColly(request.Context)
	err, post := p.FetchPost(request.GroupId, request.PostId, request.CommentNextCursor)
	return post, err
}

func (s server) FetchContentImages(ctx context.Context, request *pb.FetchContentImagesRequest) (*pb.FacebookImageList, error) {
	p := getColly(request.Context)
	err, imageList := p.FetchContentImages(request.PostId, request.NextCursor)
	return imageList, err
}

func (s server) FetchImageUrl(ctx context.Context, request *pb.FetchImageUrlRequest) (*pb.FacebookImage, error) {
	p := getColly(request.Context)
	err, image := p.FetchImageUrl(request.ImageId)
	return image, err
}

func main() {
	logger.Init("fb-colly", true, false, ioutil.Discard)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	logger.Info("Port listened at", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGrpcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
