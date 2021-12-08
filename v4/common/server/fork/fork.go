package fork

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/server"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	s *ForkServer
}

func NewServer(ctx context.Context) server.Server {
	ctx, cancel := context.WithCancel(ctx)

	return server.NewServer(ctx, &Server{
		ctx:    ctx,
		cancel: cancel,
		s:      &ForkServer{},
	})
}

func (s *Server) Serve() error {
	cmd := exec.CommandContext(s.ctx, os.Args[0], buildForkStartParams(s.s.name)...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	scannerOut := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerOut.Scan() {
			text := strings.TrimRight(scannerOut.Text(), "\n")
			log.Logger(s.ctx).Info(text)
		}
	}()

	go func() {
		for scannerErr.Scan() {
			text := strings.TrimRight(scannerOut.Text(), "\n")
			log.Logger(s.ctx).Error(text)
		}
	}()

	go func() {
		if err := cmd.Start(); err != nil {
			return
		}

		if err := cmd.Wait(); err != nil {
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.cancel()

	return nil
}

func (s *Server) Name() string {
	return "testfork"
}

func (s *Server) Metadata() map[string]string {
	return map[string]string{}
}

func (s *Server) Address() []string {
	return []string{}
}

func (s *Server) Endpoints() []string {
	return []string{}
}

func (s *Server) As(i interface{}) bool {
	v, ok := i.(**ForkServer)
	if !ok {
		return false
	}

	*v = s.s

	return true
}

type ForkServer struct {
	name string
}

func (f *ForkServer) RegisterForkParam(name string) {
	f.name = name
}

func buildForkStartParams(serviceName string) []string {

	r := fmt.Sprintf("grpc://%s", viper.GetString("grpc.address"))
	b := fmt.Sprintf("grpc://%s", viper.GetString("grpc.address"))

	params := []string{
		"start",
		"--fork",
		"--grpc.address", ":0",
		"--http.address", ":0",
		"--config", viper.GetString("config"),
		"--registry", r,
		"--broker", b,
	}
	if viper.GetBool("enable_metrics") {
		params = append(params, "--enable_metrics")
	}
	if viper.GetBool("enable_pprof") {
		params = append(params, "--enable_pprof")
	}
	// Use regexp to specify that we want to start that specific service
	params = append(params, "^"+serviceName+"$")
	return params
}
