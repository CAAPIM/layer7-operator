package util

import (
	"strings"
	"testing"
)

func TestConvertOpaqueMapToGraphmanBundle(t *testing.T) {
	secrets := []GraphmanSecret{{Name: "test1", Secret: "secret1"}, {Name: "test2", Secret: "secret2"}}
	bundleBytes, err := ConvertOpaqueMapToGraphmanBundle(secrets, []string{})
	if err != nil {
		t.Errorf("Error getting secret bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "secret1") {
		t.Errorf("bundle %s should contain %s", bundle, "secret1")
	}
	if !strings.Contains(bundle, "secret2") {
		t.Errorf("bundle %s should contain %s", bundle, "secret2")
	}

}

/*func TestConvertX509ToGraphmanBundle(t *testing.T) {
	key := GraphmanKey{
		Name: "test",
		Crt:  "-----BEGIN CERTIFICATE-----MIIC6jCCAdKgAwIBAgIGAYnTwy1CMA0GCSqGSIb3DQEBCwUAMDYxNDAyBgNVBAMMK2FNZTk4OXNSdk9YTlgyckgtYmZjVHUyZVB5NFJhVXhOOHpCenJsSzFpcFkwHhcNMjMwODA4MDYwODUxWhcNMjQwNjAzMDYwODUxWjA2MTQwMgYDVQQDDCthTWU5ODlzUnZPWE5YMnJILWJmY1R1MmVQeTRSYVV4Tjh6QnpybEsxaXBZMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0D5Q71trDwX4BnFR4jnxEBsohmD/R9CU19eWGN5iKRHVxQwsM9tR1569NZkUqAXRexo/6RHp3IT8fXdai2+227i3tpSt6hKoiVNMFznktnQ5WTZpBIa6D7iwdGZCDoY17juBrNr5ko/WWsIIVD06Z14BIs6wcqyM2QHPaKLQAgt+ZJfRps3vWCmoBHxLRhuQrcxpPhYnb/ZFsNW6fq5aJA7TG5fU7PKo69DVWUVga65ysTEEb79c7ytHHUrdEE5oR2dFmemN6yev36I92oSFqb5sBKkn2lim9VCTY6ZGZitF3XbUSqfJGkDxHIANLRi+trPdI71RKTWBtHMTkIsFhwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBJ7ShdGRSKeMVPnmb9NnQX9aZlU5Sphbb5UkgTCdd5y+8k/QpKgk+BG4u5P3wN359X1HgpDQGh3OfboMhZMJY2VnQ3qK7W0r8au6IQ5mFtlrUukBWjxAJc/1rbdBD2TlCHdEBpqgg2s7fEgeu+6NRIeJFYDLXOiQaZES01WMxL9CZDfxijwJSO6ZSSEMlDQ0K0UY3p/B0V0rSvXTrJIPE8boDzksL/0GiRBFOc0tQhqtq33h7pnKW70CjDiM7ib2fuZLLtLse+jrbZiJ79bINRmB+kd5HNJtI5xTTwXvf+sfs2v81Wdmpzdv3aKIVcnDk63+lVVh9+114QifWNNeuy-----END CERTIFICATE-----",
		Key:  "-----BEGIN PRIVATE KEY-----MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQPlDvW2sPBfgGcVHiOfEQGyiGYP9H0JTX15YY3mIpEdXFDCwz21HXnr01mRSoBdF7Gj/pEenchPx9d1qLb7bbuLe2lK3qEqiJU0wXOeS2dDlZNmkEhroPuLB0ZkIOhjXuO4Gs2vmSj9ZawghUPTpnXgEizrByrIzZAc9ootACC35kl9Gmze9YKagEfEtGG5CtzGk+Fidv9kWw1bp+rlokDtMbl9Ts8qjr0NVZRWBrrnKxMQRvv1zvK0cdSt0QTmhHZ0WZ6Y3rJ6/foj3ahIWpvmwEqSfaWKb1UJNjpkZmK0XddtRKp8kaQPEcgA0tGL62s90jvVEpNYG0cxOQiwWHAgMBAAECggEAbtF6yKXhpzEJ+IO9i6JCAswxGLHtqA3755E2sy1FF44CMMZ1j3Mbbp9vGWLJd1EBVX12nVWHGm863pnxeVqN+Qen3GXq1zHutoW5bHLGn8Hh8vPdlycLROqIHKl+ZbROZuUL8Szmu3QIImw3enzK489G03sisyPYIHOyKIDcKPl2OUEGeznXXVgzPI5LC5vXeazy7nI9ykzaBdlxf2bnykIR9RShRWxaJQ5/xqZ8hnywaGexXzv2Vpo88a3KRL9f9CyhCRR3Su0cCVua3rwj/Ijl3hx5oBPzXn0nZLeCZ8NhpC89BnQ5T0u15wKOVA9axKMFfE4PpZ/MXjaOJlOT4QKBgQD7foeM9lAhzUIW1k9qzgcR38q7BFLn8oK7FKRv1PIL0GGL+w6G00QYJ9xAQBJkZH2bzdij8dhflmbh/C3Of6ptxDuHFOoCGZArFWqbQSiMB3fb1AwGMc9hATeJ+d0Jx5/GvgWrzheVFCFXvAK9dkMIIKWYZ7xM0gA+lZrfCZuICwKBgQDT+WsGjzGIbrERirXhA1GhysYSrLP6fMGy0Ko0qlHgrwkAI1YF8eK2zqNdApqTpFoonxvwSAPICZNsh2rGup1tBjSaN3HCJC6PZonuzbApuPcfp+QZbuOBEakgp05HV2tG7cGoV48eW+FBrzuXXZgItZsZMly2VJHSxRhqVvJZ9QKBgQCwPfhyIX8AYR5qcJ9RArbToNgqfRo4b6uLvSiLMli5TLu/ZB3HADCdGPnxkLUS45Ve5T9njKkMO5M31QioyLC/oZ/xxwdCl3V/q898o4ntr6IgTJZslOV2Xmr0Z0SugNWIakwBHTlLgMLo/9mPulu5S1+g0TmVQClpsl/I46u6BwKBgQDSmK9bMfKdQJQNAImYhyqYGpRVQ14AU+hBVoxzjG+SUXQYvgKeH2YGByBIrOiUHKoyR3mDbJjNKa5dGeDcldUH1y11tfYAUuArOk15gsMtgIWM3smA9ylyNvCX74CW4mRDcL2BGZSoLdKK5qTGyobcyEjSbLWttDG4fHa4V6+p7QKBgGWLpT8TmN7qif0L07BiPTkqptqSvZRXP4lLh6wkAZxrOxzkJbSQjM+rW/cmuLWphyyMbru3xqrORzMjuc7t/3FKwbd97ZlYba8tvThNFPTA0cVAjplBIpEIcfnpBd0oPklxsj28fENz6dQ+3a1ZzSDF65kRE9/R/fgR3fw9LikM-----END PRIVATE KEY-----",
		Port: "8443"}
	keys := []GraphmanKey{key}

	bundleBytes, err := ConvertX509ToGraphmanBundle(keys)
	if err != nil {
		t.Errorf("Error getting key bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "test") {
		t.Errorf("bundle %s should contain key %s", bundle, "test")
	}
}*/
