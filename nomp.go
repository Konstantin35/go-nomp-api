package nomp

import (
	"github.com/dghubble/sling"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"log"
	"strings"
)

type NompClient struct {
	sling      *sling.Sling
	httpClient *nompHttpClient
}

// nomp dont send the api response with content type
// we fix this: set content type to json
type nompHttpClient struct {
	client    *http.Client
	debug     bool
	useragent string
}

func (d nompHttpClient) Do(req *http.Request) (*http.Response, error) {
	if d.debug {
		d.dumpRequest(req)
	}
	if d.useragent != "" {
		req.Header.Set("User-Agent", d.useragent)
	}
	client := func() (*http.Client) {
		if d.client != nil {
			return d.client
		} else {
			return http.DefaultClient
		}
	}()
	if client.Transport != nil {
		if transport, ok := client.Transport.(*http.Transport); ok {
			if transport.TLSClientConfig != nil {
				transport.TLSClientConfig.InsecureSkipVerify = true;
			} else {
				transport.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			}
		}
	} else {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			if transport.TLSClientConfig != nil {
				transport.TLSClientConfig.InsecureSkipVerify = true;
			} else {
				transport.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			}
		}
	}
	resp, err := client.Do(req)
	if d.debug {
		d.dumpResponse(resp)
	}
	if err == nil {
		contenttype := resp.Header.Get("Content-Type");
		if len(contenttype) == 0 || strings.HasPrefix(contenttype, "text/html") {
			resp.Header.Set("Content-Type", "application/json")
		}
	}
	return resp, err
}

func (d nompHttpClient) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (d nompHttpClient) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

func NewNompClient(client *http.Client, BaseURL string, UserAgent string) *NompClient {
	httpClient := &nompHttpClient{client:client, useragent:UserAgent}
	return &NompClient{
		httpClient: httpClient,
		sling: sling.New().Doer(httpClient).Base(BaseURL).Path("api/"),
	}
}

func (client NompClient) SetDebug(debug bool) {
	client.httpClient.debug = debug
}
