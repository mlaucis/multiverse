package user

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/tapglue/multiverse/platform/metrics"
)

type prepareFunc func(t *testing.T, namespace string) Service

func testList() List {
	us := List{}

	for i := 0; i < 13; i++ {
		u := testUser()

		u.Deleted = true
		u.Enabled = false

		us = append(us, u)
	}

	for i := 0; i < 7; i++ {
		us = append(us, testUser())
	}

	return us
}

func randStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	buf := make([]rune, n)

	for i := range buf {
		buf[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(buf)
}

func testServiceCount(t *testing.T, p prepareFunc) {
	var (
		customID  = randStringRunes(12)
		deleted   = true
		enabled   = true
		namespace = "service_count"
		service   = p(t, namespace)
	)

	count, err := service.Count(namespace, QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	u := testUser()
	u.CustomID = customID

	created, err := service.Put(namespace, u)
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range testList() {
		_, err := service.Put(namespace, u)
		if err != nil {
			t.Fatal(err)
		}
	}

	count, err = service.Count(namespace, QueryOptions{
		CustomIDs: []string{
			customID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		Deleted: &deleted,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 8; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		IDs: []uint64{
			created.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testServiceCreatedByDay(t *testing.T, p prepareFunc) {
	var (
		dates = []string{
			"2016-02-01",
			"2016-02-02",
			"2016-02-03",
			"2016-02-04",
			"2016-02-05",
			"2016-02-06",
			"2016-02-07",
			"2016-02-08",
			"2016-02-09",
		}
		namespace = "service_created_by_day"
		service   = p(t, namespace)
		want      = metrics.Timeseries{}
	)

	ts, err := service.CreatedByDay(namespace, time.Now(), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ts), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, d := range dates {
		count := rand.Intn(12) + 1

		for i := 0; i < count; i++ {
			created, err := time.Parse(metrics.BucketFormat, d)
			if err != nil {
				t.Fatal(err)
			}

			u := testUser()
			u.CreatedAt = created

			_, err = service.Put(namespace, u)
			if err != nil {
				t.Fatal(err)
			}
		}

		want = append(want, metrics.Datapoint{
			Bucket: d,
			Value:  count,
		})
	}

	start, err := time.Parse(metrics.BucketFormat, "2016-02-03")
	if err != nil {
		t.Fatal(err)
	}

	end, err := time.Parse(metrics.BucketFormat, "2016-02-05")
	if err != nil {
		t.Fatal(err)
	}

	ts, err = service.CreatedByDay(namespace, start, end)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ts), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
	if have := ts; reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testServicePut(t *testing.T, p prepareFunc) {
	var (
		enabled   = true
		namespace = "service_put"
		service   = p(t, namespace)
		user      = testUser()
	)

	created, err := service.Put(namespace, user)
	if err != nil {
		t.Fatal(err)
	}
	list, err := service.Query(namespace, QueryOptions{
		Enabled: &enabled,
		IDs: []uint64{
			created.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := list[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	created.Enabled = false

	updated, err := service.Put(namespace, created)
	if err != nil {
		t.Fatal(err)
	}

	list, err = service.Query(namespace, QueryOptions{
		Enabled: &created.Enabled,
		IDs: []uint64{
			created.ID,
		},
	})

	if have, want := len(list), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := list[0], updated; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = service.Put(namespace, &User{})
	if have, want := err, ErrInvalidUser; !IsInvalidUser(err) {
		t.Errorf("have %v, want %v", have, want)
	}

	invalidID := testUser()
	invalidID.ID = 1

	_, err = service.Put(namespace, invalidID)
	if have, want := err, ErrNotFound; !IsNotFound(err) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testServicePutLastRead(t *testing.T, p prepareFunc) {
	var (
		namespace = "service_put_last_read"
		service   = p(t, namespace)
		user      = testUser()
	)

	created, err := service.Put(namespace, user)
	if err != nil {
		t.Fatal(err)
	}

	format := "2006-01-02 15:04:05 UTC"

	ts, err := time.Parse(format, time.Now().Format(format))
	if err != nil {
		t.Fatal(err)
	}

	err = service.PutLastRead(namespace, created.ID, ts)
	if err != nil {
		t.Fatal(err)
	}

	list, err := service.Query(namespace, QueryOptions{
		IDs: []uint64{
			created.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	created.LastRead = ts.UTC()

	if have, want := list[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %v,\nwant %v", have, want)
	}
}

func testServiceQuery(t *testing.T, p prepareFunc) {
	var (
		customID  = randStringRunes(12)
		deleted   = true
		enabled   = true
		namespace = "service_query"
		service   = p(t, namespace)
	)

	list, err := service.Query(namespace, QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	u := testUser()
	u.CustomID = customID

	created, err := service.Put(namespace, u)
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range testList() {
		_, err := service.Put(namespace, u)
		if err != nil {
			t.Fatal(err)
		}
	}

	us, err := service.Query(namespace, QueryOptions{
		CustomIDs: []string{
			customID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(us), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	us, err = service.Query(namespace, QueryOptions{
		Deleted: &deleted,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(us), 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	us, err = service.Query(namespace, QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(us), 8; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	us, err = service.Query(namespace, QueryOptions{
		IDs: []uint64{
			created.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(us), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUser() *User {
	return &User{
		Enabled:  true,
		Password: randStringRunes(8),
		Username: randStringRunes(8),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
