package consistenthash

import (
	"fmt"
	"math"
)

// return ip list length by n
func ipFactory(n int) []string {
	var (
		domainA = 192
		domainB = 168
		domainC = 0
		domainD = 0
	)
	var res []string
	for index := 1; index < n+1; index++ {
		domainD++
		if domainD == 255 {
			domainC++
			domainD = 0

			if domainC == 255 {
				domainB++
				domainC = 0

				if domainB == 255 {
					domainA++
					domainB = 0

					if domainA == 255 {
						panic("ip overflow")
					}
				}
			}
		}

		res = append(res, fmt.Sprintf(`%d.%d.%d.%d`, domainA, domainB, domainC, domainD))
	}
	return res
}

func have(indexs []int, i int) bool {
	for _, index := range indexs {
		if i == index {
			return true
		}
	}
	return false
}

func virtual(ip string, i int) []byte {
	return []byte(fmt.Sprintf("%s#%d", ip, i))
}

func getPb(n node) float32 {
	return float32(n.num) / float32(math.MaxUint32) * 100
}
