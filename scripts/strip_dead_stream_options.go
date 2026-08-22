//go:build ignore

// Command strip_dead_stream_options removes Fern's generated SSE
// stream-reconnect options from option/ and core/ after a regen. Nothing in
// this SDK reads them (RequestOptions fields included): the only streaming
// surface is the DCP control WebSocket in drivers/mobile, which carries its
// own reconnect contract. Shipping these options would advertise a
// transparent Last-Event-ID reconnect that silently does nothing.
//
// Run by scripts/regen.sh after every regeneration. Exits non-zero when a
// pattern no longer matches, so a generator-shape change surfaces as a
// failed regen instead of silently resurrecting the dead surface.
package main

import (
	"fmt"
	"os"
	"regexp"
)

// patterns maps a file to the regexps whose matches are deleted from it.
var patterns = map[string][]string{
	"option/request_option.go": {
		`(?s)// WithMaxStreamReconnectAttempts caps.*?func WithMaxStreamReconnectAttempts.*?\n}\n\n`,
		`(?s)// WithoutStreamReconnection disables.*?func WithoutStreamReconnection.*?\n}\n\n`,
	},
	"core/request_option.go": {
		`(?s)// MaxStreamReconnectAttemptsOption implements.*?func \(m \*MaxStreamReconnectAttemptsOption\) applyRequestOptions.*?\n}\n\n`,
		`(?s)// WithoutStreamReconnectionOption implements.*?func \(w \*WithoutStreamReconnectionOption\) applyRequestOptions.*?\n}\n\n`,
		`\tMaxStreamReconnectAttempts uint\n`,
		`\tDisableStreamReconnection  bool\n`,
	},
}

func main() {
	for file, res := range patterns {
		raw, err := os.ReadFile(file)
		if err != nil {
			fail(fmt.Sprintf("read %s: %v", file, err))
		}
		src := string(raw)
		for _, pat := range res {
			re := regexp.MustCompile(pat)
			if !re.MatchString(src) {
				fail(fmt.Sprintf("%s: pattern %q matched nothing — the generator's shape changed; update this strip or re-decide the deletion", file, pat))
			}
			src = re.ReplaceAllString(src, "")
		}
		if err := os.WriteFile(file, []byte(src), 0o644); err != nil {
			fail(fmt.Sprintf("write %s: %v", file, err))
		}
	}
	fmt.Println("stripped dead stream-reconnect options")
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "strip_dead_stream_options:", msg)
	os.Exit(1)
}
