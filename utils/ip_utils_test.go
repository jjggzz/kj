package utils

import (
	"fmt"
	"testing"
)

func Test_ip(t *testing.T) {
	fmt.Println(LocalIpv4())
	fmt.Println(PublicNetworkIp())
}
