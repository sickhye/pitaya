// Copyright (c) nano Author and TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pitaya

import (
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya/acceptor"
	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/constants"
	e "github.com/topfreegames/pitaya/errors"
	"github.com/topfreegames/pitaya/helpers"
	"github.com/topfreegames/pitaya/internal/codec"
	"github.com/topfreegames/pitaya/internal/message"
	"github.com/topfreegames/pitaya/route"
	"github.com/topfreegames/pitaya/router"
	"github.com/topfreegames/pitaya/serialize/json"
	"github.com/topfreegames/pitaya/session"
	"github.com/topfreegames/pitaya/timer"
)

var (
	tables = []struct {
		isFrontend     bool
		serverType     string
		serverMode     ServerMode
		serverMetadata map[string]string
		cfg            *viper.Viper
	}{
		{true, "sv1", Cluster, map[string]string{"name": "bla"}, viper.New()},
		{false, "sv2", Standalone, map[string]string{}, viper.New()},
	}
	typeOfetcdSD        reflect.Type
	typeOfNatsRPCServer reflect.Type
	typeOfNatsRPCClient reflect.Type
)

func TestMain(m *testing.M) {
	setup()
	exit := m.Run()
	os.Exit(exit)
}

func setup() {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())

	etcdSD, _ := cluster.NewEtcdServiceDiscovery(app.config, app.server)
	typeOfetcdSD = reflect.TypeOf(etcdSD)

	natsRPCServer, _ := cluster.NewNatsRPCServer(app.config, app.server)
	typeOfNatsRPCServer = reflect.TypeOf(natsRPCServer)

	natsRPCClient, _ := cluster.NewNatsRPCClient(app.config, app.server)
	typeOfNatsRPCClient = reflect.TypeOf(natsRPCClient)
}

func initApp() {
	app = &App{
		server: &cluster.Server{
			ID:       uuid.New().String(),
			Type:     "game",
			Metadata: map[string]string{},
			Frontend: true,
		},
		debug:         false,
		startAt:       time.Now(),
		dieChan:       make(chan bool),
		acceptors:     []acceptor.Acceptor{},
		packetDecoder: codec.NewPomeloPacketDecoder(),
		packetEncoder: codec.NewPomeloPacketEncoder(),
		serverMode:    Standalone,
		serializer:    json.NewSerializer(),
		configured:    false,
		running:       false,
		router:        router.New(),
	}
}

func TestConfigure(t *testing.T) {
	for _, table := range tables {
		t.Run(table.serverType, func(t *testing.T) {
			initApp()
			Configure(table.isFrontend, table.serverType, table.serverMode, table.serverMetadata, table.cfg)
			assert.Equal(t, table.isFrontend, app.server.Frontend)
			assert.Equal(t, table.serverType, app.server.Type)
			assert.Equal(t, table.serverMode, app.serverMode)
			assert.Equal(t, table.serverMetadata, app.server.Metadata)
			assert.Equal(t, true, app.configured)
		})
	}
}

func TestAddAcceptor(t *testing.T) {
	acc := acceptor.NewTCPAcceptor("0.0.0.0:0")
	for _, table := range tables {
		t.Run(table.serverType, func(t *testing.T) {
			initApp()
			Configure(table.isFrontend, table.serverType, table.serverMode, table.serverMetadata, table.cfg)
			AddAcceptor(acc)
			if table.isFrontend {
				assert.Equal(t, acc, app.acceptors[0])
			} else {
				assert.Equal(t, 0, len(app.acceptors))
			}
		})
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	assert.Equal(t, true, app.debug)
	SetDebug(false)
	assert.Equal(t, false, app.debug)
}

func TestSetPacketDecoder(t *testing.T) {
	d := codec.NewPomeloPacketDecoder()
	SetPacketDecoder(d)
	assert.Equal(t, d, app.packetDecoder)
}

func TestSetPacketEncoder(t *testing.T) {
	e := codec.NewPomeloPacketEncoder()
	SetPacketEncoder(e)
	assert.Equal(t, e, app.packetEncoder)
}

func TestSetHeartbeatInterval(t *testing.T) {
	inter := 35 * time.Millisecond
	SetHeartbeatTime(inter)
	assert.Equal(t, inter, app.heartbeat)
}

func TestSetRPCServer(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	r, err := cluster.NewNatsRPCServer(app.config, app.server)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	SetRPCServer(r)
	assert.Equal(t, r, app.rpcServer)
}

func TestSetRPCClient(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	r, err := cluster.NewNatsRPCClient(app.config, app.server)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	SetRPCClient(r)
	assert.Equal(t, r, app.rpcClient)
}

func TestSetServiceDiscovery(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	r, err := cluster.NewEtcdServiceDiscovery(app.config, app.server)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	SetServiceDiscoveryClient(r)
	assert.Equal(t, r, app.serviceDiscovery)
}

func TestSetSerializer(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	r := json.NewSerializer()
	assert.NotNil(t, r)
	SetSerializer(r)
	assert.Equal(t, r, app.serializer)
}

func TestInitSysRemotes(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	initSysRemotes()
	assert.NotNil(t, remoteComp[0])
}

func TestSetDictionary(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())

	dict := map[string]uint16{"someroute": 12}
	err := SetDictionary(dict)
	assert.NoError(t, err)
	assert.Equal(t, dict, message.GetDictionary())

	app.running = true
	err = SetDictionary(dict)
	assert.EqualError(t, constants.ErrChangeDictionaryWhileRunning, err.Error())
}

func TestAddRoute(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	app.router = nil
	err := AddRoute("somesv", func(session *session.Session, route *route.Route, servers map[string]*cluster.Server) (*cluster.Server, error) {
		return nil, nil
	})
	assert.EqualError(t, constants.ErrRouterNotInitialized, err.Error())

	app.router = router.New()
	err = AddRoute("somesv", func(session *session.Session, route *route.Route, servers map[string]*cluster.Server) (*cluster.Server, error) {
		return nil, nil
	})
	assert.NoError(t, err)

	app.running = true
	err = AddRoute("somesv", func(session *session.Session, route *route.Route, servers map[string]*cluster.Server) (*cluster.Server, error) {
		return nil, nil
	})
	assert.EqualError(t, constants.ErrChangeRouteWhileRunning, err.Error())
}

func TestShutdown(t *testing.T) {
	initApp()
	go func() {
		Shutdown()
	}()
	<-app.dieChan
}

func TestStartDefaultSD(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	startDefaultSD()
	assert.NotNil(t, app.serviceDiscovery)
	assert.Equal(t, typeOfetcdSD, reflect.TypeOf(app.serviceDiscovery))
}

func TestStartDefaultRPCServer(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	startDefaultRPCServer()
	assert.NotNil(t, app.rpcServer)
	assert.Equal(t, typeOfNatsRPCServer, reflect.TypeOf(app.rpcServer))
}

func TestStartDefaultRPCClient(t *testing.T) {
	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, viper.New())
	startDefaultRPCClient()
	assert.NotNil(t, app.rpcClient)
	assert.Equal(t, typeOfNatsRPCClient, reflect.TypeOf(app.rpcClient))
}

func TestStartAndListenStandalone(t *testing.T) {
	initApp()
	Configure(true, "testtype", Standalone, map[string]string{}, viper.New())

	acc := acceptor.NewTCPAcceptor("0.0.0.0:0")
	AddAcceptor(acc)

	go func() {
		Start()
	}()
	helpers.ShouldEventuallyReturn(t, func() bool {
		return app.running
	}, true)

	assert.NotNil(t, handlerService)
	assert.NotNil(t, timer.GlobalTicker)
	// should be listening
	assert.NotEmpty(t, acc.GetAddr())
	helpers.ShouldEventuallyReturn(t, func() error {
		n, err := net.Dial("tcp", acc.GetAddr())
		defer n.Close()
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)
}

func ConfigureClusterApp() {

}

func TestStartAndListenCluster(t *testing.T) {
	es, cli := helpers.GetTestEtcd(t)
	defer es.Terminate(t)

	ns := helpers.GetTestNatsServer(t)
	nsAddr := ns.Addr().String()

	cfg := viper.New()
	cfg.Set("pitaya.cluster.rpc.client.nats.connect", fmt.Sprintf("nats://%s", nsAddr))
	cfg.Set("pitaya.cluster.rpc.server.nats.connect", fmt.Sprintf("nats://%s", nsAddr))

	initApp()
	Configure(true, "testtype", Cluster, map[string]string{}, cfg)

	etcdSD, err := cluster.NewEtcdServiceDiscovery(app.config, app.server, cli)
	assert.NoError(t, err)
	SetServiceDiscoveryClient(etcdSD)

	acc := acceptor.NewTCPAcceptor("0.0.0.0:0")
	AddAcceptor(acc)

	go func() {
		Start()
	}()
	helpers.ShouldEventuallyReturn(t, func() bool {
		return app.running
	}, true)

	assert.NotNil(t, handlerService)
	assert.NotNil(t, timer.GlobalTicker)
	// should be listening
	assert.NotEmpty(t, acc.GetAddr())
	helpers.ShouldEventuallyReturn(t, func() error {
		n, err := net.Dial("tcp", acc.GetAddr())
		defer n.Close()
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)
}

func TestError(t *testing.T) {
	t.Parallel()

	tables := []struct {
		name     string
		err      error
		code     string
		metadata map[string]string
	}{
		{"nil_metadata", errors.New(uuid.New().String()), uuid.New().String(), nil},
		{"empty_metadata", errors.New(uuid.New().String()), uuid.New().String(), map[string]string{}},
		{"non_empty_metadata", errors.New(uuid.New().String()), uuid.New().String(), map[string]string{"key": uuid.New().String()}},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			var err *e.Error
			if table.metadata != nil {
				err = Error(table.err, table.code, table.metadata)
			} else {
				err = Error(table.err, table.code)
			}
			assert.NotNil(t, err)
			assert.Equal(t, table.code, err.Code)
			assert.Equal(t, table.err.Error(), err.Message)
			assert.Equal(t, table.metadata, err.Metadata)
		})
	}
}
