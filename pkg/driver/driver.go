package driver

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

const (
	DefaultName = "bsos.ssun3.dev"
)

type Driver struct {
	name     string
	region   string
	endpoint string

	srv *grpc.Server

	ready bool
}

type InputParams struct {
	Name     string
	Endpoint string
	Token    string
	Region   string
}

func NewDriver(params InputParams) *Driver {
	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
		region:   params.Region,
	}
}

// start the gRPC server, like its mentioned in the CSI spec
func (d *Driver) Run() error {
	url, err := url.Parse(d.endpoint)

	if err != nil {
		return fmt.Errorf("parsing the endpoint %s", err.Error())
	}

	if url.Scheme != "unix" {
		return fmt.Errorf("only supported scheme is unix, but providered %s", url.Scheme)
	}
	grpcAddress := path.Join(url.Host, filepath.FromSlash(url.Path))

	if url.Host == "" {
		grpcAddress = filepath.FromSlash(url.Path)
	}

	if err := os.Remove(grpcAddress); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing listen address: %s", err.Error())
	}

	listener, err := net.Listen(url.Scheme, grpcAddress)

	if err != nil {
		return fmt.Errorf("listen failed %s", err.Error())
	}

	d.srv = grpc.NewServer()

	csi.RegisterNodeServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterIdentityServer(d.srv, d)

	d.ready = true

	return d.srv.Serve(listener)
}
