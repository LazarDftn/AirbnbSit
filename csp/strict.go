package csp

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/csp/internalunsafecsp"
	"strings"
)

type StrictPolicy struct {
	NoStrictDynamic bool
	UnsafeEval      bool
	BaseURI         string
	ReportURI       string
	Hashes          []string
}

func (s StrictPolicy) Serialize(nonce string, _ safehttp.InterceptorConfig) string {
	var b strings.Builder
	b.WriteString("object-src 'none'; script-src 'unsafe-inline' 'nonce-")
	b.WriteString(nonce)
	b.WriteByte('\'')

	if !s.NoStrictDynamic {
		b.WriteString(" 'strict-dynamic' https: http:")
	}

	if s.UnsafeEval {
		b.WriteString(" 'unsafe-eval'")
	}

	for _, h := range s.Hashes {
		b.WriteString(" '")
		b.WriteString(h)
		b.WriteByte('\'')
	}

	b.WriteString("; base-uri ")
	if s.BaseURI == "" {
		b.WriteString("'none'")
	} else {
		b.WriteString(s.BaseURI)
	}

	if s.ReportURI != "" {
		b.WriteString("; report-uri ")
		b.WriteString(s.ReportURI)
	}

	return b.String()
}
func (StrictPolicy) Match(cfg safehttp.InterceptorConfig) bool {
	_, ok := cfg.(internalunsafecsp.DisableStrict)
	return ok
}
func (StrictPolicy) Overridden(cfg safehttp.InterceptorConfig) (disabled, reportOnly bool) {
	disable := cfg.(internalunsafecsp.DisableStrict)
	return disable.SkipReports, true
}
