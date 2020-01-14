package userip

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const userIPKey key = 0

func FromRequest(r *http.Request) (net.IP, error) {
	ip, _, e := net.SplitHostPort(r.RemoteAddr)

	if e != nil {
		return nil, fmt.Errorf("userip: %q is not IP:Port", r.RemoteAddr)
	}
	netIp, _, e := net.ParseCIDR(ip)

	return netIp, e
}

func NewContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

func FromContext(ctx context.Context) (net.IP, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	userIP, ok := ctx.Value(userIPKey).(net.IP)
	return userIP, ok
}
