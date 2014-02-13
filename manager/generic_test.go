package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func setup_om_with_ab() (genericManager, model.Timestamp, model.Timestamp) {
	m := newGenericManager()

	A := model.Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := model.Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.register(A)
	m.register(B)

	return m, A, B
}

func TestgenericManagerGetItems(t *testing.T) {
	m, A, B := setup_om_with_ab()

	items := m.GetItems()
	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	for _, ts := range []model.Timestamp{A, B} {
		key := ts.GetKey()
		item, ok := items[key]
		if !ok {
			t.Fatalf("Missing item %s", key)
		}
		if !ts.Eq(item) {
			t.Fatalf("TS %#v does not equal %#v", ts, item)
		}
	}
}

func TestgenericManagerGetByKey(t *testing.T) {
	m, A, B := setup_om_with_ab()

	for _, ts := range []model.Timestamp{A, B} {
		key := ts.GetKey()
		item, ok := m.GetByKey(key)
		if !ok {
			t.Fatalf("Missing item %s", key)
		}
		if !ts.Eq(item) {
			t.Fatalf("TS %#v does not equal %#v", ts, item)
		}
	}
}

func TestgenericManagerRegister(t *testing.T) {
	m, A, B := setup_om_with_ab()

	group := m.GetGroup("5")
	if !(group.Contains(A) && group.Contains(B)) {
		t.Fatal("group was missing timestamps")
	}
}

func TestgenericManagerUnregister(t *testing.T) {
	m, A, B := setup_om_with_ab()

	m.unregister(A)

	group := m.GetGroup("5")
	if group.Contains(A) || !group.Contains(B) || m.Length() != 1 {
		t.Fatal("Failed to unregister A correctly")
	}
}

func TestManagableSetContains(t *testing.T) {
	ms := make(ManageableSet)
	A := model.Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := model.Timestamp{
		QHash: "abc",
		Time:  5,
	}

	if ms.Contains(A) {
		t.Fatal("ms shouldn't contain A")
	}

	ms[A.QHash] = B
	if ms.Contains(A) {
		t.Fatal("ms shouldn't contain A, but has key for A")
	}

	ms[A.QHash] = A
	if !ms.Contains(A) {
		t.Fatal("ms should contain A")
	}
}

func TestgenericManagerContains(t *testing.T) {
	m := newGenericManager()
	A := model.Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := model.Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.register(A)

	if !m.Contains(A) {
		t.Fatal("m should contain A")
	}
	if m.Contains(B) {
		t.Fatal("m should not contain B")
	}
}

func TestgenericManagerGetGroup(t *testing.T) {
	m := newGenericManager()
	group := m.GetGroup("5")

	if len(group) != 0 {
		t.Fatal("group should have been empty")
	}

	ts := model.Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	m.register(&ts)

	group = m.GetGroup("20")
	if len(group) != 0 {
		t.Fatal("group should still have been empty")
	}
	group = m.GetGroup("5")
	if len(group) != 1 {
		t.Fatal("group should not have been empty")
	}

	if group[ts.QHash] != &ts {
		t.Fatal("group contents should have contained &ts")
	}

	if len(m.by_group) != 2 {
		t.Fatalf(
			"Not caching groups - expected 2 groups, got %d: %v",
			len(m.by_group),
			m.by_group,
		)
	}
}