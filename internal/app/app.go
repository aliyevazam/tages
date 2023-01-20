package app

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aliyevazam/tages/internal/pkg/db"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/juju/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	tages "github.com/aliyevazam/tages/genproto"
	"github.com/aliyevazam/tages/internal/controller/service"
	"github.com/aliyevazam/tages/internal/controller/storage/filestore"
	"github.com/aliyevazam/tages/internal/pkg/config"
	"github.com/aliyevazam/tages/internal/pkg/logger"
)

const (
	gatherTime = 5 * time.Second
	streamcap  = 10
	unarycap   = 100
)

type rateLimiterInterceptor struct {
	TokenBucket *ratelimit.Bucket
}

func (r *rateLimiterInterceptor) Limit() bool {
	// debug
	fmt.Printf("Token Avail %d \n", r.TokenBucket.Available())

	// if zero we reached rate limit, so return true ( report error to Grpc)
	tokenRes := r.TokenBucket.TakeAvailable(1)
	if tokenRes == 0 {
		fmt.Printf("Reached Rate-Limiting %d \n", r.TokenBucket.Available())
		return true
	}

	// if tokenRes is not zero, means gRpc request can continue to flow without rate limiting :)
	return false
}

func Run(cfg config.Config) {
	l := logger.New(cfg.LogLevel)

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("error while connDB, err := db.ConnectToDB(cfg): %v"), err)
	}
	fileStore := filestore.NewDiskFileStore("files")

	TagesService := service.NewTagesService(connDB, l, fileStore)

	lis, err := net.Listen("tcp", ":"+cfg.ServicePort)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpcclient.New: %w", err))
	}

	limiterUnary := &rateLimiterInterceptor{}
	limiterStream := &rateLimiterInterceptor{}

	limiterUnary.TokenBucket = ratelimit.NewBucket(gatherTime, int64(unarycap))
	limiterStream.TokenBucket = ratelimit.NewBucket(gatherTime, int64(streamcap))

	c := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ratelimit.UnaryServerInterceptor(limiterUnary),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ratelimit.StreamServerInterceptor(limiterStream)),
	)
	reflection.Register(c)
	tages.RegisterTagesServiceServer(c, TagesService)
	fmt.Println(cfg.ServicePort)
	l.Info("Server is running on" + "port" + ": " + cfg.ServicePort)

	if err := c.Serve(lis); err != nil {
		log.Fatal("Error while listening: ", err)
	}

}
