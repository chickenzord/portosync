package main

import (
	"strings"

	"github.com/chickenzord/portosync/internal/server"
)

// parseKseiAccountsWithName parses a string of KSEI accounts in the format
// "name:username:password,name2:username2:password2" and returns a map of
// account names to Account structs containing username and password.
func parseKseiAccountsWithName(s string) map[string]server.Account {
	accounts := make(map[string]server.Account)

	entries := strings.SplitSeq(s, ",")
	for entry := range entries {
		parts := strings.SplitN(entry, ":", 3)
		if len(parts) == 3 {
			name := strings.TrimSpace(parts[0])
			username := strings.TrimSpace(parts[1])
			password := strings.TrimSpace(parts[2])

			if name != "" {
				accounts[name] = server.Account{
					Username: username,
					Password: password,
				}
			}
		}
	}

	return accounts
}
