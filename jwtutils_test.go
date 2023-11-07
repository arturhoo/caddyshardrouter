package caddyshardrouter

import "testing"

var tokenStr = "eyJhbGciOiJSUzI1NiJ9.eyJjdXN0b21lciI6IndhbG1hcnQifQ.PZBoUWYkdcjnGO7mccnUYXuHTdrYGLSDR0GkbcwjtoNFG8OU4a-ALGTXgHrIijWerxb4f53Y4XqXtA0xWnkNiht6g1aUFzXwf2_kYPoEg2JRJUHJwR0pDsdSJHWi2pN9gnxTQETUNVOdokptTkCHOcHgJdA4g3Ywy83Sud9x5Apwbe0UZrU7yir7cIEu_HXHoeok2sxMSf1al0Kl6GwlamVB09edgkRFbx9953u-H6KHCC3u_Ku2zlif13JKiawnAqsO8RQtX1NzcWr2jdl0SQvLV0MuvIsQ-yr9w-t2tLKRbwhTYnbaARHHzK2GOtgg_ALcmz562N011P3YPtMqzg"

func TestParseJWT(t *testing.T) {
	claims, err := ParseJWT(tokenStr)
	if err != nil {
		t.Error("got an error:", err)
	}
	customer, _ := claims["customer"].(string)
	expected := "walmart"

	if customer != expected {
		t.Errorf("expected '%s' but got '%s'", expected, customer)
	}
}

func BenchmarkParseJWT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := ParseJWT(tokenStr)
		if err != nil {
			b.Error("got an error", err)
		}
	}
}
