package authproxy

import (
	"fmt"
	"testing"
)

func TestRingPCache(t *testing.T) {
	var cache = NewRingPCache(1000, 60)
	if len(cache.Entries) != 1000 {
		t.Error(fmt.Sprintf("Cache not initialized to expected size %v != %v", 1000, len(cache.Entries)))
	}
	entry := &cache.Entries[0]
	if _, valid := cache.Lookup(entry.Token, entry.Action, entry.Resource); valid {
		t.Error("empty cache should not have valid lookup")
	}
	for i := 0; i < len(cache.Entries); i = i + 1 {
		entry = &cache.Entries[i]
		if entry.CreateTime != 0 {
			t.Error("empty cache not initialized as expected")
		}
	}
	// add a couple thousand entries and try different lookups
	for i := 0; i < 2000; i = i + 1 {
		token := fmt.Sprintf("token%v", i)
		action := "action-" + token
		resource := "resource-" + token
		if _, ok := cache.Lookup(token, action, resource); ok {
			t.Error("Lookup should fail before insert!")
		}
		cache.Add(token, action, resource, false)
		tail := (i + 1) % len(cache.Entries)
		if tail != cache.Tail {
			t.Error(fmt.Sprintf("ring cache tail not where expected %v != %v", tail, cache.Tail))
			break
		}
		if result, ok := cache.Lookup(token, action, resource); !(ok && !result) {
			t.Error(fmt.Sprintf("lookup should be ok %v, but not give access %v", ok, result))
			break
		}

		// add same key again, different access value
		cache.Add(token, action, resource, true)
		if tail != cache.Tail {
			t.Error(fmt.Sprintf("ring cache tail on 2nd same-key add not where expected %v != %v", tail, cache.Tail))
			break
		}
		if result, ok := cache.Lookup(token, action, resource); !(ok && result) {
			t.Error(fmt.Sprintf("lookup 2 should be ok %v, and give access %v", ok, result))
			break
		}

		// check old entries
		if i > 999 {
			j := i - 999
			token := fmt.Sprintf("token%v", j)
			action := "action-" + token
			resource := "resource-" + token

			if result, ok := cache.Lookup(token, action, resource); !(ok && result) {
				t.Error(fmt.Sprintf("lookup old should be ok %v, and give access %v", ok, result))
				break
			}
		}
	}
}
