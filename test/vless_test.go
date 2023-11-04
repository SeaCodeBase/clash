package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/require"

	"github.com/ClashCore/clash/adapter/outbound"
	C "github.com/ClashCore/clash/constant"
)

func TestClash_Vless(t *testing.T) {
	configPath := C.Path.Resolve("vless.json")

	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds:        []string{fmt.Sprintf("%s:/etc/v2ray/config.json", configPath)},
	}

	id, err := startContainer(cfg, hostCfg, "vless")
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:   "vless",
		Server: localIP.String(),
		Port:   10002,
		UUID:   "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher: "auto",
		UDP:    true,
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessTLS(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-tls.json")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/fullchain.pem", C.Path.Resolve("example.org.pem")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/privkey.pem", C.Path.Resolve("example.org-key.pem")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-tls")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:           "vless",
		Server:         localIP.String(),
		Port:           10002,
		UUID:           "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:         "auto",
		TLS:            true,
		SkipCertVerify: true,
		ServerName:     "example.org",
		UDP:            true,
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessHTTP2(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-http2.json")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/fullchain.pem", C.Path.Resolve("example.org.pem")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/privkey.pem", C.Path.Resolve("example.org-key.pem")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-http2")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:           "vless",
		Server:         localIP.String(),
		Port:           10002,
		UUID:           "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:         "auto",
		Network:        "h2",
		TLS:            true,
		SkipCertVerify: true,
		ServerName:     "example.org",
		UDP:            true,
		HTTP2Opts: outbound.HTTP2Options{
			Host: []string{"example.org"},
			Path: "/test",
		},
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessHTTP(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-http.json")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-http")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:    "vless",
		Server:  localIP.String(),
		Port:    10002,
		UUID:    "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:  "auto",
		Network: "http",
		UDP:     true,
		HTTPOpts: outbound.HTTPOptions{
			Method: "GET",
			Path:   []string{"/"},
			Headers: map[string][]string{
				"Host": {"www.amazon.com"},
				"User-Agent": {
					"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36 Edg/84.0.522.49",
				},
				"Accept-Encoding": {
					"gzip, deflate",
				},
				"Connection": {
					"keep-alive",
				},
				"Pragma": {"no-cache"},
			},
		},
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessWebsocket(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-ws.json")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-ws")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:    "vless",
		Server:  localIP.String(),
		Port:    10002,
		UUID:    "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:  "auto",
		Network: "ws",
		UDP:     true,
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessWebsocketTLS(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-ws-tls.json")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/fullchain.pem", C.Path.Resolve("example.org.pem")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/privkey.pem", C.Path.Resolve("example.org-key.pem")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-ws-tls")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:           "vless",
		Server:         localIP.String(),
		Port:           10002,
		UUID:           "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:         "auto",
		Network:        "ws",
		TLS:            true,
		SkipCertVerify: true,
		UDP:            true,
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessWebsocketTLSZero(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-ws-tls-zero.json")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/fullchain.pem", C.Path.Resolve("example.org.pem")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/privkey.pem", C.Path.Resolve("example.org-key.pem")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-ws-tls-zero")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:           "vless",
		Server:         localIP.String(),
		Port:           10002,
		UUID:           "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:         "zero",
		Network:        "ws",
		TLS:            true,
		SkipCertVerify: true,
		UDP:            true,
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessGrpc(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-grpc.json")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/fullchain.pem", C.Path.Resolve("example.org.pem")),
			fmt.Sprintf("%s:/etc/ssl/v2ray/privkey.pem", C.Path.Resolve("example.org-key.pem")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-grpc")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:           "vless",
		Server:         localIP.String(),
		Port:           10002,
		UUID:           "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:         "auto",
		Network:        "grpc",
		TLS:            true,
		SkipCertVerify: true,
		UDP:            true,
		ServerName:     "example.org",
		GrpcOpts: outbound.GrpcOptions{
			GrpcServiceName: "example!",
		},
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessWebsocket0RTT(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/v2ray/config.json", C.Path.Resolve("vless-ws-0rtt.json")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-ws-0rtt")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:       "vless",
		Server:     localIP.String(),
		Port:       10002,
		UUID:       "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:     "auto",
		Network:    "ws",
		UDP:        true,
		ServerName: "example.org",
		WSOpts: outbound.WSOptions{
			MaxEarlyData:        2048,
			EarlyDataHeaderName: "Sec-WebSocket-Protocol",
		},
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func TestClash_VlessWebsocketXray0RTT(t *testing.T) {
	cfg := &container.Config{
		Image:        ImageXray,
		ExposedPorts: defaultExposedPorts,
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds: []string{
			fmt.Sprintf("%s:/etc/xray/config.json", C.Path.Resolve("vless-ws-0rtt.json")),
		},
	}

	id, err := startContainer(cfg, hostCfg, "vless-xray-ws-0rtt")
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:       "vless",
		Server:     localIP.String(),
		Port:       10002,
		UUID:       "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:     "auto",
		Network:    "ws",
		UDP:        true,
		ServerName: "example.org",
		WSOpts: outbound.WSOptions{
			Path: "/?ed=2048",
		},
	})
	require.NoError(t, err)

	time.Sleep(waitTime)
	testSuit(t, proxy)
}

func Benchmark_Vless(b *testing.B) {
	configPath := C.Path.Resolve("vless.json")

	cfg := &container.Config{
		Image:        ImageVmess,
		ExposedPorts: defaultExposedPorts,
		Entrypoint:   []string{"/usr/bin/v2ray"},
		Cmd:          []string{"run", "-c", "/etc/v2ray/config.json"},
	}
	hostCfg := &container.HostConfig{
		PortBindings: defaultPortBindings,
		Binds:        []string{fmt.Sprintf("%s:/etc/v2ray/config.json", configPath)},
	}

	id, err := startContainer(cfg, hostCfg, "vless-bench")
	require.NoError(b, err)

	b.Cleanup(func() {
		cleanContainer(id)
	})

	proxy, err := outbound.NewVless(outbound.VlessOption{
		Name:    "vless",
		Server:  localIP.String(),
		Port:    10002,
		UUID:    "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher:  "auto",
		AlterID: 0,
		UDP:     true,
	})
	require.NoError(b, err)

	time.Sleep(waitTime)
	benchmarkProxy(b, proxy)
}
