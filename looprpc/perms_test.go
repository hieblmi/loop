package looprpc

import (
	"testing"

	"gopkg.in/macaroon-bakery.v2/bakery"
)

func TestInstantOutRequiresLoopOutPermission(t *testing.T) {
	requiredPerms, ok := RequiredPermissions["/looprpc.SwapClient/InstantOut"]
	if !ok {
		t.Fatalf("InstantOut permission entry missing")
	}

	assertPermission := func(want bakery.Op) {
		t.Helper()

		for _, perm := range requiredPerms {
			if perm == want {
				return
			}
		}

		t.Fatalf("InstantOut permission entry missing %v", want)
	}

	assertPermission(bakery.Op{
		Entity: "swap",
		Action: "execute",
	})
	assertPermission(bakery.Op{
		Entity: "loop",
		Action: "out",
	})
}
