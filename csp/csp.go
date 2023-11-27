package csp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-safeweb/safehttp/plugins/csp/internalunsafecsp"
	"github.com/google/go-safeweb/safehttp/plugins/htmlinject"

	"github.com/google/go-safeweb/safehttp"
)

const (
	responseHeaderKey           = "Content-Security-Policy"
	responseHeaderReportOnlyKey = responseHeaderKey + "-Report-Only"
)

const nonceSize = 20

func generateNonce() string {
	b := make([]byte, nonceSize)
	_, err := internalunsafecsp.RandReader.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate entropy using crypto/rand/RandReader: %v", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}

type key string

const (
	nonceKey   key = "csp-nonce"
	headersKey key = "csp-headers"
)

// Nonce retrieves the nonce from the given context. If there is no nonce stored
// in the context, an error will be returned.
func Nonce(ctx context.Context) (string, error) {
	v := safehttp.FlightValues(ctx).Get(nonceKey)
	if v == nil {
		return "", errors.New("no nonce in context")
	}
	return v.(string), nil
}

func nonce(r *safehttp.IncomingRequest) string {
	v := safehttp.FlightValues(r.Context()).Get(nonceKey)
	var nonce string
	if v == nil {
		nonce = generateNonce()
		safehttp.FlightValues(r.Context()).Put(nonceKey, nonce)
	} else {
		nonce = v.(string)
	}
	return nonce
}

func claimedHeaders(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) (cspe func([]string), cspro func([]string)) {
	type claimed struct {
		cspe, cspro func([]string)
	}
	v := safehttp.FlightValues(r.Context()).Get(headersKey)
	var c claimed
	if v == nil {
		h := w.Header()
		cspe := h.Claim(responseHeaderKey)
		cspro := h.Claim(responseHeaderReportOnlyKey)
		c = claimed{cspe: cspe, cspro: cspro}
		safehttp.FlightValues(r.Context()).Put(headersKey, c)
	} else {
		c = v.(claimed)
	}
	return c.cspe, c.cspro
}

type Policy interface {
	Serialize(nonce string, cfg safehttp.InterceptorConfig) string
	Match(cfg safehttp.InterceptorConfig) bool
	Overridden(cfg safehttp.InterceptorConfig) (disabled, reportOnly bool)
}

func report(reportURI string) string {
	var b strings.Builder

	if reportURI != "" {
		b.WriteString("report-uri ")
		b.WriteString(reportURI)
		b.WriteString("; ")
	}

	return b.String()
}

type Interceptor struct {
	Policy     Policy
	ReportOnly bool
}

var _ safehttp.Interceptor = Interceptor{}

func Default(reportURI string) []Interceptor {
	return []Interceptor{
		{Policy: StrictPolicy{ReportURI: reportURI}},
		{Policy: TrustedTypesPolicy{ReportURI: reportURI}},
	}
}

func (it Interceptor) processOverride(cfg safehttp.InterceptorConfig, nonce string) (enf, ro string) {
	disabled, reportOnly := false, false
	if it.Policy.Match(cfg) {
		disabled, reportOnly = it.Policy.Overridden(cfg)
	}
	if disabled {
		return "", ""
	}
	p := it.Policy.Serialize(nonce, cfg)
	if reportOnly || it.ReportOnly {
		return "", p
	}
	return p, ""
}

func (it Interceptor) Before(w safehttp.ResponseWriter, r *safehttp.IncomingRequest, cfg safehttp.InterceptorConfig) safehttp.Result {
	nonce := nonce(r)
	enf, ro := it.processOverride(cfg, nonce)
	setCSP, setCSPReportOnly := claimedHeaders(w, r)
	if enf != "" {
		prev := w.Header().Values(responseHeaderKey)
		setCSP(append(prev, enf))
	}
	if ro != "" {
		prev := w.Header().Values(responseHeaderReportOnlyKey)
		setCSPReportOnly(append(prev, ro))
	}
	return safehttp.NotWritten()
}
func (it Interceptor) Commit(w safehttp.ResponseHeadersWriter, r *safehttp.IncomingRequest, resp safehttp.Response, cfg safehttp.InterceptorConfig) {
	tmplResp, ok := resp.(*safehttp.TemplateResponse)
	if !ok {
		return
	}

	nonce, err := Nonce(r.Context())
	if err != nil {
		// The nonce should have been added in the Before stage and, if that is
		// not the case, a server misconfiguration occurred.
		panic("no CSP nonce")
	}

	if tmplResp.FuncMap == nil {
		tmplResp.FuncMap = map[string]interface{}{}
	}
	tmplResp.FuncMap[htmlinject.CSPNoncesDefaultFuncName] = func() string { return nonce }
}

func (it Interceptor) Match(cfg safehttp.InterceptorConfig) bool {
	return it.Policy.Match(cfg)
}
