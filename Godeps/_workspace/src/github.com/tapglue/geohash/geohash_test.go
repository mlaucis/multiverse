/**
 * This code is licensed under MIT license.
 * Please see LICENSE.md file for full license.
 */

package geohash_test

import (
	"fmt"
	"math"
	"testing"

	. "github.com/tapglue/geohash"
)

func TestEncodeInt(t *testing.T) {
	var expected uint64 = 4064984913515641
	encodedVal := EncodeInt(37.8324, 112.5584, 52)
	if encodedVal != expected {
		t.Logf("Expected: %d Got: %d\n", encodedVal, expected)
		t.Fail()
	}
}

func TestDecodeInt(t *testing.T) {
	lat, lon, _, _ := DecodeInt(4064984913515641, 52)
	if math.Abs(37.8324-lat) >= 0.0001 {
		t.Logf("37.8324 - %.5f was >= 0.0001\n", lat)
		t.Fail()
	}
	if math.Abs(112.5584-lon) >= 0.0001 {
		t.Logf("112.5584 - %.5f was >= 0.0001\n", lon)
		t.Fail()
	}
}

func TestNeighborsInt(t *testing.T) {
	neighbors := EncodeNeighborsInt(1702789509, 32)
	neighborsTest := []uint64{
		1702789520,
		1702789522,
		1702789511,
		1702789510,
		1702789508,
		1702789422,
		1702789423,
		1702789434,
	}

	for key := range neighbors {
		if neighborsTest[key] != neighbors[key] {
			t.Logf("neighbor %d mismatch expected: %d got: %d", key, neighborsTest[key], neighbors[key])
			t.Fail()
		}
	}

	singleNeighbor := EncodeNeighborInt(1702789509, 1, 0, 52)
	if neighbors[0] != singleNeighbor {
		t.Logf("neighbor mismatch expected: %d got: %d", neighbors[0], singleNeighbor)
		t.Fail()
	}

	neighbors = EncodeNeighborsInt(27898503327470, 52)
	neighborsTest = []uint64{
		27898503327471,
		27898503349317,
		27898503349316,
		27898503349313,
		27898503327467,
		27898503327465,
		27898503327468,
		27898503327469,
	}

	for key := range neighbors {
		if neighborsTest[key] != neighbors[key] {
			t.Logf("neighbor %d mismatch expected: %d got: %d", key, neighborsTest[key], neighbors[key])
			t.Fail()
		}
	}

	singleNeighbor = EncodeNeighborInt(27898503327470, -1, -1, 46)
	if neighbors[5] != singleNeighbor {
		t.Logf("neighbor mismatch expected: %d got: %d", neighbors[0], singleNeighbor)
		t.Fail()
	}
}

func TestDistanceBetweenHashPoints(t *testing.T) {
	start := EncodeInt(50.066389, -5.714722, 52)
	end := EncodeInt(58.643889, -3.070000, 52)

	distance := fmt.Sprintf("%.4f", DistanceBetweenHashPoints(start, end, 52))

	if "982.4086" != distance {
		t.Logf("distance is different expected: %.4f got: %s", 982.4086, distance)
		t.Fail()
	}
}
