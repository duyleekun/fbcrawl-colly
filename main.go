package main

import "C"
import (
	context "context"
	"flag"
	"github.com/google/logger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"qnetwork.net/fbcrawl/fbcrawl"
	"qnetwork.net/fbcrawl/fbcrawl/pb"
)

const logPath = "parse.log"

var limiter = rate.NewLimiter(1, 1)

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
	limiter.Wait(ctx)
	p := getColly(nil)
	return p.Login(request.Email, request.Password, request.TotpSecret)
}

func (s server) FetchMyGroups(ctx context.Context, request *pb.FetchMyGroupsRequest) (*pb.FacebookGroupList, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchMyGroups()
}

func (s server) FetchGroupInfo(ctx context.Context, request *pb.FetchGroupInfoRequest) (*pb.FacebookGroup, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchGroupInfo(request.GroupUsername)
}

func (s server) FetchUserInfo(ctx context.Context, request *pb.FetchUserInfoRequest) (*pb.FacebookUser, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchUserInfo(request.Username)
}

func (s server) FetchGroupFeed(ctx context.Context, request *pb.FetchGroupFeedRequest) (*pb.FacebookPostList, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchGroupFeed(request.GroupId, request.NextCursor)
}

func (s server) FetchPost(ctx context.Context, request *pb.FetchPostRequest) (*pb.FacebookPost, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchPost(request.GroupId, request.PostId, request.CommentNextCursor)
}

func (s server) FetchContentImages(ctx context.Context, request *pb.FetchContentImagesRequest) (*pb.FacebookImageList, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchContentImages(request.PostId, request.NextCursor)
}

func (s server) FetchImageUrl(ctx context.Context, request *pb.FetchImageUrlRequest) (*pb.FacebookImage, error) {
	limiter.Wait(ctx)
	p := getColly(request.Context)
	return p.FetchImageUrl(request.ImageId)
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
