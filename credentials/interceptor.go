package credentials

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authorization unary interceptor function to handle authorize per RPC call
func JWTServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("=====JWTServerInterceptor=====")
	start := time.Now()
	// Skip authorize when GetJWT is requested
	if info.FullMethod != "/proto.EventStoreService/GetJWT" {
		if err := authorize(ctx); err != nil {
			log.Println("JWTServerInterceptor, authentication error: ", err)
			return nil, err
		}
	}
	log.Println("=====DONE JWTServerInterceptor===")
	// Calls the handler
	h, err := handler(ctx, req)

	// Logging with grpclog (grpclog.LoggerV2)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}

// authorize function authorizes the token received from Metadata
func authorize(ctx context.Context) error {
	log.Println("authorize")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]
	secr, err := ioutil.ReadFile(JWT_SECRET_KEY)
	if err != nil {
		log.Printf("Err open secret_key: %v", err)
	}
	// validateToken function validates the token
	_, err = ValidateToken(token, string(secr))

	/*
		Can add more logic here to validate the returned token.
		Like checking more Claims
	*/

	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}
	return nil
}
