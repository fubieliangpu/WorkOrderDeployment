package mtools

import (
	"log"
	"strconv"
	"strings"
)

func GetGatewayByIpStr(ipaddrs ...string) (gips []string) {
	gips = make([]string, 0, len(ipaddrs))
	for _, ipaddr := range ipaddrs {
		convertIp := strings.Split(ipaddr, ".")
		res, err := strconv.Atoi(convertIp[3])
		if err != nil {
			log.Fatal(err)
		}
		res++
		convertIp[3] = strconv.Itoa(res)
		gips = append(gips, strings.Join(convertIp, "."))
	}

	return gips
}
