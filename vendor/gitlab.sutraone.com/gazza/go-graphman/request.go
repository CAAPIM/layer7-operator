package graphman

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type CustomTransport struct {
	username string
	password string
	encpass  string
	r        http.RoundTripper
}

func (ct *CustomTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(ct.username, ct.password)
	r.Header.Set("encpass", ct.encpass)
	return ct.r.RoundTrip(r)
}

func gqlClient(username string, password string, target string, encpass string) graphql.Client {
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: &CustomTransport{username: username, password: password, encpass: encpass, r: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
	}

	return graphql.NewClient(target, httpClient)
}

func GqlSummary() {
	//client := gqlClient()
	// resp, err := summary(context.Background(), client)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(resp.ClusterProperties)

}
