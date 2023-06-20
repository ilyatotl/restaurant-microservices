package authentication

import (
	"context"
	"google.golang.org/grpc"
	"order_pocessor/internal/app/pb"
)

type Client struct {
	authClient pb.AuthenticationServiceClient
	conn       *grpc.ClientConn
}

func NewClient(_ context.Context, target string) (*Client, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		authClient: pb.NewAuthenticationServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Authorize(ctx context.Context, token string) (*pb.AuthenticationResponse, error) {
	return c.authClient.Authorize(ctx, &pb.AuthenticationRequest{Token: token})
}
