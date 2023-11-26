package csp

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/csp/internalunsafecsp"
	"strings"
)

type TrustedTypesPolicy struct {
	ReportURI string
}

func (t TrustedTypesPolicy) Serialize(nonce string, _ safehttp.InterceptorConfig) string {
	var b strings.Builder
	b.WriteString("require-trusted-types-for 'script'")

	if t.ReportURI != "" {
		b.WriteString("; report-uri ")
		b.WriteString(t.ReportURI)
	}

	return b.String()
}
func (TrustedTypesPolicy) Match(cfg safehttp.InterceptorConfig) bool {
	_, ok := cfg.(internalunsafecsp.DisableTrustedTypes)
	return ok
}
func (TrustedTypesPolicy) Overridden(cfg safehttp.InterceptorConfig) (disabled, reportOnly bool) {
	disable := cfg.(internalunsafecsp.DisableTrustedTypes)
	return disable.SkipReports, true
}
