package cidrxpndr

import (
        "net"
)

// Expand takes a CIDR notation string and returns a 
// []string slice of IPs.
func Expand(s string) ([]string, error) {
        ip, ipnet, err := net.ParseCIDR(s)
        if err != nil {
                return nil, err
        }

        // Get network bits.
        nb, _ := ipnet.Mask.Size()

        // Short-circuit a /32.
        if nb == 32 {
                return []string{ip.String()}, nil
        }

        // Get usable hosts count. Init slice with the size;
        // /16 and larger will start to cause runtime slow downs
        // if we used append().
        nHosts := 2<<uint(31-nb) - 2
        ips := make([]string, nHosts)

        // net.IP slice start position.
        p := 15

        // Increment start IP by # hosts and populate slice.
        for n := 0; n < nHosts; n++ {
                ip[15]++
                // Increment the next class if the current position is > 254.
                if ip[15] > 254 {
                        ip[15] = 1
                decPos:
                        p--
                        for n := 15 - p; n > 0; n-- {
                                // Break at Class A limit.
                                if n > 3 {
                                        break
                                }
                                ip[15-n]++
                                if ip[15-n] > 254 {
                                        ip[15-n] = 0
                                        goto decPos
                                }
                        }
                        ips[n] = ip.String()
                } else {
                        ips[n] = ip.String()
                }
        }

        return ips, nil
}
