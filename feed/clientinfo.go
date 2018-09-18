package feed

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type key int

var userKey key = 0
var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type User struct {
	UserName string
	ClientId string
	Scope    []string
}

func getKey(t *jwt.Token) (interface{}, error) {
	return []byte("dashur"), nil
}

func valid(authorization []string) (*User, error) {
	if len(authorization) < 1 {
		return nil, errInvalidToken
	}
	tokenStr := strings.TrimPrefix(authorization[0], "Bearer ")
	token, err := jwt.Parse(tokenStr, getKey)
	fmt.Println(tokenStr)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, errInvalidToken
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range claims {
		fmt.Printf("%s\t%v\n", key, value)
	}

	s := make([]string, len(claims["scope"].([]interface{})))
	for i, v := range claims["scope"].([]interface{}) {
		s[i] = fmt.Sprint(v)
	}

	return &User{
		UserName: claims["user_name"].(string),
		ClientId: claims["client_id"].(string),
		Scope:    s,
	}, nil
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	user, err := valid(md["authorization"])

	if err != nil {
		return nil, err
	} else {
		newCtx := context.WithValue(ctx, userKey, user)
		return handler(newCtx, req)
	}
}
