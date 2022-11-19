package vector

import (
	"context"
	"flag"
	"fmt"
	"github.com/digital-dream-labs/hugh/grpc/client"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

// Vector is the struct containing info about Vector
type Vector struct {
	Conn vectorpb.ExternalInterfaceClient
}

// New returns either a vector struct, or an error on failure
func New(opts ...Option) (*Vector, error) {
	cfg := options{}

	homedir, _ := os.UserHomeDir()
	dirname := filepath.Join(homedir, ".anki_vector", "sdk_config.ini")

	if initData, _ := ini.Load(dirname); initData != nil {
		sec, _ := initData.GetSection(ini.DefaultSection)
		sec.MapTo(&cfg)
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Target == "" || cfg.Token == "" {
		return nil, fmt.Errorf("configuration options missing")
	}

	creds, err0 := credentials.NewClientTLSFromFile("005070ac", "true")
	if err0 != nil {
		return nil, err0
	}

	c, err := client.New(
		client.WithTarget(cfg.Target),
		client.WithInsecureSkipVerify(),
		client.WithDialopts(
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(
				tokenAuth{
					token: cfg.Token,
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

// NewEP returns either a vector struct for escape pod vector, or an error on failure
func NewEP(opts ...Option) (*Vector, error) {
	cfg := options{}

	homedir, _ := os.UserHomeDir()
	dirname := filepath.Join(homedir, ".anki_vector", "sdk_config.ini")

	if initData, _ := ini.Load(dirname); initData != nil {
		sec, _ := initData.GetSection(ini.DefaultSection)
		sec.MapTo(&cfg)
	} else {
		return nil, fmt.Errorf("INI file missing")
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Target == "" {
		var host = flag.String("host", "", "Vector's IP address")
		flag.Parse()
		if *host == "" {
			log.Fatal("please use the -host argument and set it to your robots IP address")
		}
		cfg.Target = fmt.Sprintf("%s:443", *host)
	}
	if cfg.Target == "" {
		return nil, fmt.Errorf("configuration options missing")
	}

	c, err := client.New(
		client.WithTarget(cfg.Target),
		client.WithInsecureSkipVerify(),
	)
	if err != nil {
		return nil, err
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}

	vc := vectorpb.NewExternalInterfaceClient(c.Conn())

	login, err := vc.UserAuthentication(context.Background(),
		&vectorpb.UserAuthenticationRequest{
			UserSessionId: []byte("bullshit1"),
			ClientName:    []byte("bullshit2"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return New(
		WithTarget(cfg.Target),
		WithToken(string(login.ClientTokenGuid)),
	)
}
