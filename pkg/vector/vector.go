package vector

import (
	"fmt"
	"github.com/digital-dream-labs/hugh/grpc/client"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Vector is the struct containing info about Vector
type Vector struct {
	Conn vectorpb.ExternalInterfaceClient
}

// New returns either a vector struct, or an error on failure
func New(opts ...Option) (*Vector, error) {
	cfg := options{}

	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.target == "" || cfg.token == "" {
		return nil, fmt.Errorf("configuration options missing")
	}

	//var trusted_certs, _ = ioutil.ReadFile()

	/*
		trusted_certs = cert.read()
			channel_credentials = aiogrpc.ssl_channel_credentials(root_certificates=trusted_certs)
			            # Add authorization header for all the calls
			            call_credentials = aiogrpc.access_token_call_credentials(self._guid)

			            credentials = aiogrpc.composite_channel_credentials(channel_credentials, call_credentials)

			            self._logger.info(f"Connecting to {self.host} for {self.name} using {self.cert_file}")
			            self._channel = aiogrpc.secure_channel(self.host, credentials,
			                                                   options=(("grpc.ssl_target_name_override", self.name,),))
	*/
	creds, err0 := credentials.NewClientTLSFromFile("005070ac", "true")
	if err0 != nil {
		return nil, err0
	}

	c, err := client.New(
		client.WithTarget(cfg.target),
		client.WithInsecureSkipVerify(),
		client.WithDialopts(
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(
				tokenAuth{
					token: "q0MulUycT6aThkfFavXoog==",
				},
			),
		),
	)
	print("Vector created")
	if err != nil {
		return nil, err
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}

	print("Vector connected")
	r := Vector{
		Conn: vectorpb.NewExternalInterfaceClient(c.Conn()),
	}

	return &r, nil
}
