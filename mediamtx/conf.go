// Package mediamtx runs MediaMTX, the RTMP server that takes the feed in.
package mediamtx

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/carol-sh/nakanohito/config"
)

// WriteConf writes a MediaMTX config that takes RTMP only: the ingest user
// may publish to its path, and localhost may read it and use the API.
func WriteConf(path string, c *config.Config) error {
	q := func(s string) string {
		b, _ := json.Marshal(s) // a JSON string is a YAML string
		return string(b)
	}
	r := strings.NewReplacer(
		"$rtmp", q(c.MediaMTX.RTMP),
		"$api", q(c.MediaMTX.API),
		"$path", q(c.Ingest.Path),
		"$user", q(c.Ingest.User),
		"$pass", q(c.Ingest.Password),
	)
	conf := r.Replace(`# Written by nakanohito from config.json. Edits here are lost.
rtmp: yes
rtmpAddress: $rtmp
rtsp: no
hls: no
webrtc: no
srt: no
moq: no

api: yes
apiAddress: $api

authMethod: internal
authInternalUsers:
  - user: $user
    pass: $pass
    ips: []
    permissions:
      - action: publish
        path: $path
  - user: any
    pass:
    ips: ["127.0.0.1", "::1"]
    permissions:
      - action: read
        path: $path
      - action: api

paths:
  $path: {}
`)
	if err := os.WriteFile(path, []byte(conf), 0o600); err != nil {
		return err
	}
	return os.Chmod(path, 0o600)
}
