package unlock

import "net/http"

// defaultUA is a plausible desktop User-Agent used by every probe.
// Streaming services fingerprint on UA; a missing or suspicious one
// causes extra challenges (captcha, 403) that are noise in a probe.
const defaultUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

// drain reads and discards at most 4 KB of the response body so the
// connection can be reused. Probes don't need the payload beyond what
// their per-target handler already reads.
func drain(resp *http.Response) error {
	buf := make([]byte, 4096)
	_, err := resp.Body.Read(buf)
	if err != nil && err.Error() == "EOF" {
		return nil
	}
	return err
}
