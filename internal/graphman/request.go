/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
package graphman

import (
	"context"
	"crypto/tls"
	"net"
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

func dialTimeout(ctx context.Context, network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 3*time.Second)
}

func (ct *CustomTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(ct.username, ct.password)
	r.Header.Set("l7-passphrase", ct.encpass)
	return ct.r.RoundTrip(r)
}

func gqlClient(username string, password string, target string, encpass string) graphql.Client {
	httpClient := &http.Client{
		Timeout:   time.Second * 60,
		Transport: &CustomTransport{username: username, password: password, encpass: encpass, r: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DialContext: dialTimeout}},
	}
	return graphql.NewClient(target, httpClient)
}
