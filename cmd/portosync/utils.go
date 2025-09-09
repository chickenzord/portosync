package main

import "strings"

func parseKseiAccountsString(s string) map[string]string {
	accounts := make(map[string]string)

	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			accounts[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return accounts
}
