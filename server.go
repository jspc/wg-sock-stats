package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jspc/wg-sock-stats/stats"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	stats.UnimplementedStatsServer
	dumper Dumper
	config Config
	ipdb   IPDB
}

func New(c Config, i IPDB) Server {
	return Server{
		stats.UnimplementedStatsServer{},
		WGDump,
		c,
		i,
	}
}

func (s Server) Get(ctx context.Context, in *stats.Statistics) (*stats.Statistics, error) {
	startTime := time.Now()

	output, err := s.dumper()
	if err != nil {
		log.Print(err)

		return nil, status.Error(codes.Unavailable, "unable to access stats")
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))

	line := 0
	sent := 0.0
	received := 0.0

	for scanner.Scan() {
		t := scanner.Text()
		f := strings.Fields(t)

		switch line {
		case 0:
			in.Address = fmt.Sprintf(":%s", f[3])
			in.PublicKey = f[2]
			in.Datetime = timestamppb.New(startTime)

		default:
			peerSent := toFloat(f[6])
			sent += peerSent

			peerReceived := toFloat(f[7])
			received += peerReceived

			p := &stats.Peer{
				PublicKey:     f[1],
				Address:       f[3],
				AllowedIPs:    f[4],
				Handshake:     timestamppb.New(toTime(f[5])),
				SentBytes:     peerSent,
				Sent:          humanize.Bytes(uint64(peerSent)),
				ReceivedBytes: peerReceived,
				Received:      humanize.Bytes(uint64(peerReceived)),
			}

			if s.config.CheckPTR {
				p.Addr = s.parsePTR(f[3])
			}

			if s.config.MapOwners {
				owner, ok := s.config.Owners[f[1]]
				if ok {
					p.Owner = &stats.Owner{
						Name:  owner.Name,
						Email: owner.Email,
					}
				}
			}

			in.Peer = append(in.Peer, p)
		}

		line++
	}

	in.SentBytes = sent
	in.Sent = humanize.Bytes(uint64(sent))

	in.ReceivedBytes = received
	in.Received = humanize.Bytes(uint64(received))

	return in, nil
}

// Given an IP address, lookup PTR records, geolocations
// and anything else useful for stats
func (s Server) parsePTR(in string) (a *stats.Address) {
	a = new(stats.Address)

	// If we can't process an IP address then there's no point continuing
	tcpAddr, err := net.ResolveTCPAddr("", in)
	if err != nil {
		return
	}

	a.Addr = tcpAddr.IP.String()

	ptr, _ := net.LookupAddr(a.Addr)
	switch len(ptr) {
	case 0:
	default:
		a.Ptr = ptr[0]
	}

	record, err := s.ipdb.City(tcpAddr.IP)
	if err != nil {
		return
	}

	a.Long = record.Location.Longitude
	a.Lat = record.Location.Latitude
	a.City = record.City.Names["en"]
	a.Country = record.Country.Names["en"]

	ispRecord, err := s.ipdb.ASN(tcpAddr.IP)
	if err != nil {
		return
	}

	a.Isp = ispRecord.AutonomousSystemOrganization

	return
}

// toTime absorbs errors in conversion; if we can't figure out a
// timestamp then fine
func toTime(s string) (t time.Time) {
	dT, err := strconv.Atoi(s)
	if err != nil {
		return
	}

	return time.Unix(int64(dT), 0).UTC()
}

// toFloat turns a string to a float, discarding errors
func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)

	return f
}
