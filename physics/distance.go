package physics

import "math"

const (
	twoPi float64 = 2 * math.Pi
)

// DistCalc represents a wrapper for distance calculation operations
type DistCalc interface {
	Execute(pos Vec2, angle float64) (float64, []float64)
}

type distCalcImpl struct {
	coords   []Vec2
	nextSegs []Seg
	prevSegs []Seg
	len      int
	halfDims Vec2
	sensors  []Seg
}

type distCalcInstance struct {
	centre          Vec2
	front           Seg
	shape           Rect
	frontProjection Rect
}

//  3    0
//  | -> |      Car corners, -> is forward
//  2    1
func newDistCalcInstance(halfDims, position Vec2, angle float64) *distCalcInstance {
	shape := Rect{
		Vec2{-halfDims[0], halfDims[1]}.Rotate(angle).Add(position),
		Vec2{halfDims[0], halfDims[1]}.Rotate(angle).Add(position),
		Vec2{halfDims[0], -halfDims[1]}.Rotate(angle).Add(position),
		Vec2{-halfDims[0], -halfDims[1]}.Rotate(angle).Add(position),
	}

	leftSideProjectionDirection := shape[0].Subtract(shape[3]).Normalise()
	rightSideProjectionDirection := shape[1].Subtract(shape[2]).Normalise()

	frontProjection := Rect{
		shape[0].Add(leftSideProjectionDirection.ScalarMultiple(projectionLimit)),
		shape[1].Add(rightSideProjectionDirection.ScalarMultiple(projectionLimit)),
		shape[1],
		shape[0],
	}

	return &distCalcInstance{
		centre:          position,
		front:           Seg{shape[0], shape[1]},
		shape:           shape,
		frontProjection: frontProjection,
	}
}

// faster p % q where -q < p < q
func qMod(p, q int) int {
	if p >= q {
		return p - q
	}
	if p < 0 {
		return p + q
	}
	return p
}

// NewDistCalc creates a new distance calculator
func NewDistCalc(dims Vec2, paths []Path, sensorRange float64) DistCalc {
	outerLen := len(paths)

	nextSegs := make([][]Seg, outerLen)
	prevSegs := make([][]Seg, outerLen)

	total := 0

	for outerIdx, path := range paths {
		innerLength := len(path)

		// Form corners for a path, have coord in the path, then edges at nextSegs[idx]
		// and prevSegs[idx]
		nextSegs[outerIdx] = make([]Seg, innerLength)
		prevSegs[outerIdx] = make([]Seg, innerLength)

		total += innerLength

		for innerIdx, coord := range path {
			nextSegs[outerIdx][innerIdx][0] = coord
			nextSegs[outerIdx][innerIdx][1] = path[qMod(innerIdx+1, innerLength)]
			prevSegs[outerIdx][innerIdx][0] = coord
			nextSegs[outerIdx][innerIdx][1] = path[qMod(innerIdx-1, innerLength)]
		}
	}

	resultNextSegs := make([]Seg, total)
	resultPrevSegs := make([]Seg, total)
	coords := make([]Vec2, total)

	cursor := 0

	for outerIdx, path := range paths {
		for innerIdx := range path {
			coords[cursor+innerIdx] = paths[outerIdx][innerIdx]
			resultNextSegs[cursor+innerIdx] = nextSegs[outerIdx][innerIdx]
			resultPrevSegs[cursor+innerIdx] = prevSegs[outerIdx][innerIdx]
		}
		cursor += len(path)
	}

	halfDims := dims.ScalarMultiple(0.5)

	halfDimsMagnitude := halfDims.Magnitude()
	halfDimsRatio := sensorRange / halfDimsMagnitude
	magnitudeAdjustedHalfDims := halfDims.ScalarMultiple(halfDimsRatio)

	sensors := []Seg{
		// Left Side
		Seg{Vec2{-halfDims[0], 0}, Vec2{-(halfDims[0] + sensorRange), 0}},
		// Left corner
		Seg{Vec2{-halfDims[0], halfDims[1]}, Vec2{-magnitudeAdjustedHalfDims[0], magnitudeAdjustedHalfDims[1]}},
		// Center
		Seg{Vec2{0, halfDims[1]}, Vec2{0, halfDims[1] + sensorRange}},
		// Right corner
		Seg{halfDims, magnitudeAdjustedHalfDims},
		// Right side
		Seg{Vec2{halfDims[0], 0}, Vec2{halfDims[0] + sensorRange, 0}},
	}

	return &distCalcImpl{
		len:      cursor,
		coords:   coords,
		nextSegs: resultNextSegs,
		prevSegs: resultPrevSegs,
		halfDims: halfDims,
		sensors:  sensors,
	}
}

func updateMin(min *float64, p Vec2, q *Vec2) {
	if q == nil {
		return
	}

	v := math.Abs(p.Subtract(*q).MagnitudeSquared())
	*min = math.Min(*min, v)
}

func (d *distCalcImpl) executeAgainstProjection(c *distCalcInstance) float64 {
	// TODO: optimise, can prune on various levels e.g velocity range
	leftProjectionSide := Seg{c.frontProjection[0], c.frontProjection[3]}
	rightProjectionSide := Seg{c.frontProjection[1], c.frontProjection[2]}

	minimum := distanceCap

	for idx, coord := range d.coords {
		if c.frontProjection.Contains(coord) {
			prevIntersect := Seg{c.centre, d.prevSegs[idx][1]}.Intersection(c.front)
			centreIntersect := Seg{c.centre, coord}.Intersection(c.front)
			nextIntersect := Seg{c.centre, d.nextSegs[idx][1]}.Intersection(c.front)

			prevDist := math.Abs(prevIntersect.Subtract(d.prevSegs[idx][1]).MagnitudeSquared())
			midDist := math.Abs(centreIntersect.Subtract(coord).MagnitudeSquared())
			nextDist := math.Abs(nextIntersect.Subtract(d.nextSegs[idx][1]).MagnitudeSquared())

			// TODO: investigate fast way to do point -> line shortest path rather than approximation here
			// (for cases 2/3)

			if midDist <= prevDist && midDist <= nextDist {
				// Target is minimum
				minimum = math.Min(minimum, midDist)
				continue
			} else if prevDist < nextDist && prevDist < midDist {
				// Previous is minimum, edge check
				leftPrevIntersect := d.prevSegs[idx].Intersection(leftProjectionSide)
				rightPrevIntersect := d.prevSegs[idx].Intersection(rightProjectionSide)

				updateMin(&minimum, c.shape[0], leftPrevIntersect)
				updateMin(&minimum, c.shape[1], rightPrevIntersect)
			} else {
				// Next is minimum, edge check
				leftNextIntersect := d.nextSegs[idx].Intersection(leftProjectionSide)
				rightNextIntersect := d.nextSegs[idx].Intersection(rightProjectionSide)

				updateMin(&minimum, c.shape[0], leftNextIntersect)
				updateMin(&minimum, c.shape[1], rightNextIntersect)
			}
		} else {
			// Edge (nxt)
			leftIntersect := d.nextSegs[idx].Intersection(leftProjectionSide)
			rightIntersect := d.nextSegs[idx].Intersection(rightProjectionSide)

			updateMin(&minimum, c.shape[0], leftIntersect)
			updateMin(&minimum, c.shape[1], rightIntersect)
		}
	}

	return math.Sqrt(minimum)
}

func normaliseAngle(a float64) float64 {
	return a - (twoPi * math.Floor(a/twoPi))
}

func (d *distCalcImpl) Execute(position Vec2, angle float64) (float64, []float64) {
	correctedAngle := normaliseAngle(angle)
	instance := newDistCalcInstance(d.halfDims, position, correctedAngle)

	collisionDistance := d.executeAgainstProjection(instance)

	sensorAngles := make([]float64, len(d.sensors))

	for idx, sensor := range d.sensors {
		sensorPoint0 := sensor[0].Rotate(correctedAngle).Add(position)
		sensorPoint1 := sensor[1].Rotate(correctedAngle).Add(position)
		sensorSegment := Seg{sensorPoint0, sensorPoint1}

		minimum := distanceCap

		for _, seg := range d.nextSegs {
			sensorIntersect := sensorSegment.Intersection(seg)
			updateMin(&minimum, sensor[0], sensorIntersect)
		}

		sensorAngles[idx] = minimum
	}

	return collisionDistance, sensorAngles
}
