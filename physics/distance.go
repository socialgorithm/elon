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
	coords  []Vec2
	nsegs   []Seg
	psegs   []Seg
	len     int
	hdims   Vec2
	sensors []Seg
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
func newDistCalcInstance(hdims, pos Vec2, angle float64) *distCalcInstance {
	shape := Rect{
		Vec2{-hdims[0], hdims[1]}.Rotate(angle).Add(pos),
		Vec2{hdims[0], hdims[1]}.Rotate(angle).Add(pos),
		Vec2{hdims[0], -hdims[1]}.Rotate(angle).Add(pos),
		Vec2{-hdims[0], -hdims[1]}.Rotate(angle).Add(pos),
	}

	ld := shape[0].Subtract(shape[3]).Normalise()
	rd := shape[1].Subtract(shape[2]).Normalise()

	frontProjection := Rect{
		shape[0].Add(ld.ScalarMultiple(projectionLimit)),
		shape[1].Add(rd.ScalarMultiple(projectionLimit)),
		shape[1],
		shape[0],
	}

	return &distCalcInstance{
		centre:          pos,
		front:           Seg{shape[0], shape[1]},
		shape:           shape,
		frontProjection: frontProjection,
	}
}

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
	olen := len(paths)

	nsegs := make([][]Seg, olen)
	psegs := make([][]Seg, olen)

	total := 0

	for oidx, path := range paths {
		ilen := len(path)
		nsegs[oidx] = make([]Seg, ilen)
		psegs[oidx] = make([]Seg, ilen)
		total += ilen

		for iidx, t := range path {
			nsegs[oidx][iidx][0] = t
			nsegs[oidx][iidx][1] = path[qMod(iidx+1, ilen)]
			psegs[oidx][iidx][0] = t
			nsegs[oidx][iidx][1] = path[qMod(iidx-1, ilen)]
		}
	}

	cnsegs := make([]Seg, total)
	cpsegs := make([]Seg, total)
	coords := make([]Vec2, total)

	cursor := 0

	for oidx, path := range paths {
		for iidx := range path {
			coords[cursor+iidx] = paths[oidx][iidx]
			cnsegs[cursor+iidx] = nsegs[oidx][iidx]
			cpsegs[cursor+iidx] = psegs[oidx][iidx]
		}
		cursor += len(path)
	}

	hdims := dims.ScalarMultiple(0.5)

	hdm := hdims.Magnitude()
	hdmr := sensorRange / hdm
	hdma := hdims.ScalarMultiple(hdmr)

	sensors := []Seg{
		Seg{Vec2{-hdims[0], 0}, Vec2{-(hdims[0] + sensorRange), 0}}, // L
		Seg{Vec2{-hdims[0], hdims[1]}, Vec2{-hdma[0], hdma[1]}},     // LC
		Seg{Vec2{0, hdims[1]}, Vec2{0, hdims[1] + sensorRange}},     // C
		Seg{hdims, hdma},                                            // RC
		Seg{Vec2{hdims[0], 0}, Vec2{hdims[0] + sensorRange, 0}},     // R
	}

	return &distCalcImpl{
		len:     cursor,
		coords:  coords,
		nsegs:   cnsegs,
		psegs:   cpsegs,
		hdims:   hdims,
		sensors: sensors,
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
	lpr := Seg{c.frontProjection[0], c.frontProjection[3]}
	rpr := Seg{c.frontProjection[1], c.frontProjection[2]}

	min := distanceCap

	for idx, v := range d.coords {
		if c.frontProjection.Contains(v) {
			pi := Seg{c.centre, d.psegs[idx][1]}.Intersection(c.front)
			mi := Seg{c.centre, v}.Intersection(c.front)
			ni := Seg{c.centre, d.nsegs[idx][1]}.Intersection(c.front)

			pd := math.Abs(pi.Subtract(d.psegs[idx][1]).MagnitudeSquared())
			md := math.Abs(mi.Subtract(v).MagnitudeSquared())
			nd := math.Abs(ni.Subtract(d.nsegs[idx][1]).MagnitudeSquared())

			if md <= pd && md <= nd {
				// Target is minimum
				min = math.Min(min, md)
				continue
			} else if pd < nd && pd < md {
				// Previous is minimum, edge check
				lpin := d.psegs[idx].Intersection(lpr)
				rpin := d.psegs[idx].Intersection(rpr)

				updateMin(&min, c.shape[0], lpin)
				updateMin(&min, c.shape[1], rpin)
			} else {
				// Next is minimum, edge check
				lnin := d.nsegs[idx].Intersection(lpr)
				rnin := d.nsegs[idx].Intersection(rpr)

				updateMin(&min, c.shape[0], lnin)
				updateMin(&min, c.shape[1], rnin)
			}
		} else {
			// Edge (nxt)
			li := d.nsegs[idx].Intersection(lpr)
			ri := d.nsegs[idx].Intersection(rpr)

			updateMin(&min, c.shape[0], li)
			updateMin(&min, c.shape[1], ri)
		}
	}

	return math.Sqrt(min)
}

func normaliseAngle(a float64) float64 {
	return a - (twoPi * math.Floor(a/twoPi))
}

func (d *distCalcImpl) Execute(pos Vec2, angle float64) (float64, []float64) {
	na := normaliseAngle(angle)
	i := newDistCalcInstance(d.hdims, pos, na)

	p := d.executeAgainstProjection(i)

	sr := make([]float64, len(d.sensors))

	for idx, s := range d.sensors {
		ns0 := s[0].Rotate(na).Add(pos)
		ns1 := s[1].Rotate(na).Add(pos)
		ns := Seg{ns0, ns1}

		min := distanceCap

		for _, seg := range d.nsegs {
			d := ns.Intersection(seg)
			updateMin(&min, s[0], d)
		}

		sr[idx] = min
	}

	return p, sr
}
