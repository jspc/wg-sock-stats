package main

import (
	"flag"
	"log"
	"net"
	"path/filepath"

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

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
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
