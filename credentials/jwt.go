package credentials

import (
	"context"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type JWT struct {
	token string
}

func NewTokenFromFile(token string) (credentials.PerRPCCredentials, error) {
	data, err := ioutil.ReadFile(token)
	if err != nil {
		return JWT{}, err
	}
	return JWT{string(data)}, nil
}

func (j JWT) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j JWT) RequireTransportSecurity() bool {
	return true
}
