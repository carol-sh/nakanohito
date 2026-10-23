// nakanohito takes a remote RTMP feed and puts it on air through the local OBS.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/carol-sh/nakanohito/config"
	"github.com/carol-sh/nakanohito/mediamtx"
)

const dataDir = "data"

func main() {
	c, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Another RTMP server may already hold the ports.
	for _, addr := range []string{c.MediaMTX.RTMP, c.MediaMTX.API} {
		if err := free(addr); err != nil {
			log.Fatal(err)
		}
	}

	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Fatal(err)
	}
	conf := filepath.Join(dataDir, "mediamtx.yml")
	if err := mediamtx.WriteConf(conf, c); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", conf)
}

// free reports whether addr is free. Go listeners set SO_REUSEADDR, so
// listening alone can succeed beside a server on 0.0.0.0; dial first.
func free(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second); err == nil {
		_ = conn.Close()
		return fmt.Errorf("%s is in use", addr)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("%s is in use: %w", addr, err)
	}
	return l.Close()
}
