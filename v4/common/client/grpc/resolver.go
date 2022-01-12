package grpc

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	pb "github.com/pydio/cells/v4/common/proto/registry"
	"github.com/pydio/cells/v4/common/registry"
)

const (
	defaultPort = "8001"
)

var (
	errMissingAddr = errors.New("cells resolver: missing address")

	errAddrMisMatch = errors.New("cells resolver: invalid uri")

	regex, _ = regexp.Compile("^([A-z0-9.]*?)(:[0-9]{1,5})?\\/([A-z_]*)$")
)

func init() {
	// resolver.Register(NewBuilder())
}

type cellsBuilder struct {
	reg registry.Registry
}

type cellsResolver struct {
	reg                  registry.Registry
	address              string
	cc                   resolver.ClientConn
	name                 string
	m                    map[string][]string
	updatedState         chan struct{}
	updatedStateTimer    *time.Timer
	disableServiceConfig bool
}

func NewBuilder(reg registry.Registry) resolver.Builder {
	return &cellsBuilder{reg: reg}
}

func (b *cellsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	_, _, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		return nil, err
	}

	cr := &cellsResolver{
		reg:                  b.reg,
		name:                 name,
		cc:                   cc,
		m:                    map[string][]string{},
		disableServiceConfig: opts.DisableServiceConfig,
		updatedStateTimer:    time.NewTimer(100 * time.Millisecond),
	}

	// fmt.Println("Building much ? ")
	// debug.PrintStack()

	go cr.updateState()
	go cr.watch()

	services, err := b.reg.List(registry.WithType(pb.ItemType_SERVICE))
	if err != nil {
		return nil, err
	}

	for _, s := range services {
		for _, n := range s.(registry.Service).Nodes() {
			for _, a := range n.Address() {
				cr.m[a] = append(cr.m[a], s.Name())
			}
		}
	}

	cr.sendState()

	return cr, nil
}

func (cr *cellsResolver) watch() {
	w, err := cr.reg.Watch(registry.WithType(pb.ItemType_SERVICE))
	if err != nil {
		return
	}

	fmt.Println("Initial Registry watch done")
	for {
		r, err := w.Next()
		if err != nil {
			return
		}

		var s registry.Service
		if r.Item().As(&s) && (r.Action() == "create" || r.Action() == "update") {
			for _, n := range s.Nodes() {
				for _, a := range n.Address() {
					cr.m[a] = append(cr.m[a], s.Name())
				}
				// cr.m[n.Address()[0]] = append(cr.m[n.Address()[0]], s.Name())
			}

			cr.updatedStateTimer.Reset(1 * time.Second)
		}
	}
}

func (cr *cellsResolver) updateState() {
	for {
		select {
		case <-cr.updatedStateTimer.C:
			cr.sendState()
		}
	}
}

func (cr *cellsResolver) sendState() {
	var addresses []resolver.Address
	for k, v := range cr.m {
		addresses = append(addresses, resolver.Address{
			Addr:       k,
			ServerName: "main",
			Attributes: attributes.New("services", v),
		})
	}
	if len(addresses) == 0 {
		// dont' bother sending yet
		return
	}

	if err := cr.cc.UpdateState(resolver.State{
		Addresses:     addresses,
		ServiceConfig: cr.cc.ParseServiceConfig(`{"loadBalancingPolicy": "lb"}`),
	}); err != nil {
		fmt.Println("And the error is ? ", err)
	}
}

func (b *cellsBuilder) Scheme() string {
	return "cells"
}

func (cr *cellsResolver) ResolveNow(opt resolver.ResolveNowOptions) {
}

func (cr *cellsResolver) Close() {
}

func parseTarget(target string) (host, port, name string, err error) {
	if target == "" {
		return "", "", "", errMissingAddr
	}

	if !regex.MatchString(target) {
		return "", "", "", errAddrMisMatch
	}

	groups := regex.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]

	return host, port, name, nil
}
