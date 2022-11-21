package vector

import (
	"fmt"
	"github.com/digital-dream-labs/hugh/grpc/client"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"google.golang.org/grpc"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

// Vector is the struct containing info about Vector
type Vector struct {
	Conn vectorpb.ExternalInterfaceClient
	Cfg  options
}

// New returns either a vector struct, or an error on failure
func New(opts ...Option) (*Vector, error) {
	cfg := options{}

	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Target == "" || cfg.Token == "" {
		return nil, fmt.Errorf("configuration options missing")
	}

	c, err := client.New(
		client.WithTarget(cfg.Target),
		client.WithInsecureSkipVerify(),
		client.WithDialopts(
			grpc.WithPerRPCCredentials(
				tokenAuth{
					token: cfg.Token,
				},
			),
		),
	)
	if err != nil {
		return nil, err
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}

	r := Vector{
		Conn: vectorpb.NewExternalInterfaceClient(c.Conn()),
		Cfg:  cfg,
	}

	return &r, nil
}

// NewEP returns either a vector struct for escape pod vector, or an error on failure
func NewEP(serial string) (*Vector, error) {
	if serial == "" {
		log.Fatal("please use the -serial argument and set it to your robots serial number")
		return nil, fmt.Errorf("Configuration options missing")
	}

	cfg := options{}

	homedir, _ := os.UserHomeDir()
	dirname := filepath.Join(homedir, ".anki_vector", "sdk_config.ini")

	if initData, _ := ini.Load(dirname); initData != nil {
		sec, _ := initData.GetSection(serial)
		sec.MapTo(&cfg)
	} else {
		return nil, fmt.Errorf("INI file missing")
	}

	cfg.SerialNo = serial
	cfg.Target = fmt.Sprintf("%s:443", cfg.Target)

	println(cfg.SerialNo)
	println(cfg.Target)
	println(cfg.Token)
	println(cfg.CertPath)
	println(cfg.RobotName)

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

	return New(
		WithTarget(cfg.Target),
		WithToken(cfg.Token),
	)
}
