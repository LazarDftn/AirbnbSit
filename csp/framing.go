package csp

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/framing/internalunsafeframing"
	"strings"
)

type FramingPolicy struct {
	ReportURI string
}

func (f FramingPolicy) Serialize(nonce string, cfg safehttp.InterceptorConfig) string {
	var b strings.Builder

	var allow []string
	if a, ok := cfg.(internalunsafeframing.AllowList); ok {
		allow = a.Hostnames
	}
	b.WriteString(frameAncestors(allow))
	b.WriteString(report(f.ReportURI))

	return strings.TrimSpace(b.String())
}

// Match matches strict policies overrides.
func (FramingPolicy) Match(cfg safehttp.InterceptorConfig) bool {
	switch cfg.(type) {
	case internalunsafeframing.Disable, internalunsafeframing.AllowList:
		return true
	}
	return false
}

func (FramingPolicy) Overridden(cfg safehttp.InterceptorConfig) (disabled, reportOnly bool) {
	switch c := cfg.(type) {
	case internalunsafeframing.Disable:
		return c.SkipReports, true
	case internalunsafeframing.AllowList:
		return false, c.ReportOnly
	}
	// This should not happen.
	return false, false
}

func frameAncestors(sources []string) string {
	var b strings.Builder
	b.WriteString("frame-ancestors 'self'")
	for _, s := range sources {
		b.WriteString(" ")
		b.WriteString(s)
	}
	b.WriteString("; ")
	return b.String()
}
