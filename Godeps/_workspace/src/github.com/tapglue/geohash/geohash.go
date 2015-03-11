/**
 * This code is licensed under MIT license.
 * Please see LICENSE.md file for full license.
 */

// Package geohash implements geo hashing encoding and decoding functions
package geohash

import "math"

// EncodeInt encodes the latitude and longitude into a uint64 hash
func EncodeInt(latitude, longitude float64, bitDepth uint8) uint64 {
	var (
		bitsTotal    uint8   = 0
		maxLat       float64 = 90
		minLat       float64 = -90
		maxLon       float64 = 180
		minLon       float64 = -180
		mid          float64 = 0
		combinedBits uint64  = 0
	)

	for bitsTotal < bitDepth {
		combinedBits *= 2
		if bitsTotal%2 == 0 {
			mid = (maxLon + minLon) / 2
			if longitude > mid {
				combinedBits += 1
				minLon = mid
			} else {
				maxLon = mid
			}
		} else {
			mid = (maxLat + minLat) / 2
			if latitude > mid {
				combinedBits += 1
				minLat = mid
			} else {
				maxLat = mid
			}
		}
		bitsTotal++
	}

	return combinedBits
}

// DecodeInt decodes a hash procduced by EncodeInt into the original latitude and longitude
func DecodeInt(position uint64, bitDepth uint8) (float64, float64, float64, float64) {
	rLat, rLon, rLatErr, rLonErr := decodeBboxInt(position, bitDepth)
	lat := (rLat + rLatErr) / 2
	lon := (rLon + rLonErr) / 2
	latErr := rLatErr - lat
	lonErr := rLonErr - lon
	return lat, lon, latErr, lonErr
}

// EncodeNeighborInt returns the neighbor of a location from a certain direction
func EncodeNeighborInt(location uint64, l, r float64, bitDepth uint8) uint64 {
	lat, lon, latErr, lonErr := DecodeInt(location, bitDepth)
	neighborLat := lat + l*latErr*2
	neighborLon := lon + r*lonErr*2

	return EncodeInt(neighborLat, neighborLon, bitDepth)
}

// EncodeNeighborsInt returns the encoded neighbors of a certain location
func EncodeNeighborsInt(location uint64, bitDepth uint8) []uint64 {
	lat, lon, latErr, lonErr := DecodeInt(location, bitDepth)
	latErr *= 2
	lonErr *= 2

	encodeNeighborInt := func(l, r float64) uint64 {
		neighborLat := lat + l*latErr
		neighborLon := lon + r*lonErr
		return EncodeInt(neighborLat, neighborLon, bitDepth)
	}

	return []uint64{
		encodeNeighborInt(1, 0),
		encodeNeighborInt(1, 1),
		encodeNeighborInt(0, 1),
		encodeNeighborInt(-1, 1),
		encodeNeighborInt(-1, 0),
		encodeNeighborInt(-1, -1),
		encodeNeighborInt(0, -1),
		encodeNeighborInt(1, -1),
	}
}

// DistanceBetweenPoints computes the distance between two geo hashes
func DistanceBetweenHashPoints(start, end uint64, bitDepth uint8) float64 {
	startLat, startLon, _, _ := DecodeInt(start, bitDepth)
	endLat, endLon, _, _ := DecodeInt(end, bitDepth)

	return DistanceBetweenPoints(startLat, startLon, endLat, endLon)
}

// DistanceBetweenPoints computes the distance between two lat & lon coordinates
func DistanceBetweenPoints(startLat, startLon, endLat, endLon float64) float64 {
	dLat := rad(endLat - startLat)
	dLon := rad(endLon - startLon)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rad(startLat))*math.Cos(rad(startLon))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	return 6371 * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func rad(d float64) float64 {
	return d * math.Pi / 180
}

func deg(r float64) float64 {
	return r / math.Pi / 180
}

func decodeBboxInt(locationHash uint64, bitDepth uint8) (float64, float64, float64, float64) {
	var (
		maxLat float64 = 90
		minLat float64 = -90
		maxLon float64 = 180
		minLon float64 = -180
		step           = bitDepth / 2
		i      uint8   = 0
	)

	for i = 0; i < step; i++ {
		lonBit := getBit(locationHash, ((step-i)*2)-1)
		latBit := getBit(locationHash, ((step-i)*2)-2)

		if latBit == 0 {
			maxLat = (maxLat + minLat) / 2
		} else {
			minLat = (maxLat + minLat) / 2
		}

		if lonBit == 0 {
			maxLon = (maxLon + minLon) / 2
		} else {
			minLon = (maxLon + minLon) / 2
		}
	}

	return minLat, minLon, maxLat, maxLon
}

func getBit(bits uint64, position uint8) uint64 {
	return (bits / uint64(math.Pow(2, float64(position)))) & 0x01
}
