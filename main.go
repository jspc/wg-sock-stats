package main

import (
	"flag"
	"log"
	"net"
	"path/filepath"
	"regexp"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jspc/wg-sock-stats/stats"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configFile = flag.String("c", "wg-sock-stats.toml", "Configuration file for wg-sock-stats (note: missing or empty file starts app with default, privacy focused configuration")
	geoIPDir   = flag.String("g", "/etc/wg-sock-stats/geoip2", "Directory containing at least the GeoIP2 City and ASN databases. If config.CheckPTR is false no lookups are performed and these files need not exist")
	listenAddr = flag.String("a", "unix:///var/run/wg-sock-stats.sock", "Address on which to listen")

	reSchema = regexp.MustCompile("^(?:([a-z0-9]+)://)?(.*)$")
)

func main() {
	var (
		config Config
		ipdb   IPDB
		err    error
	)

	flag.Parse()

	config, err = ParseConfig(*configFile)
	if err != nil {
		log.Printf("Config file %q does not exist, using default, privacy focused config", *configFile)
	}

	if config.CheckPTR {
		ipdb, err = NewIPDB(filepath.Join(*geoIPDir, "GeoLite2-City.mmdb"), filepath.Join(*geoIPDir, "GeoLite2-ASN.mmdb"))
		if err != nil {
			log.Panic(err)
		}
	}

	server := New(config, ipdb)

	proto, addr := SplitSchemaAddr(*listenAddr)
	lis, err := net.Listen(proto, addr)
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)

	reflection.Register(grpcServer)

	stats.RegisterStatsServer(grpcServer, server)

	grpcServer.Serve(lis)
}

// SplitSchemaAddr takes an address from flags and derives
// the listen protocol and address to pass to net.Listen
//
// This allows us to use tcp/ unix sockets
//
// This code comes from: https://github.com/grpc/grpc-go/issues/1846
// and is reused here with gratitude
func SplitSchemaAddr(addr string) (string, string) {
	parts := reSchema.FindStringSubmatch(addr)
	proto, addr := parts[1], parts[2]
	if proto == "" {
		proto = "tcp"
	}
	return proto, addr
}
