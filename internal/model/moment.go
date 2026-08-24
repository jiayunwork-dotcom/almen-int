package model

// BendingMoment returns the moment applied to the strip cross-section by the
// residual compressive layer, in N*mm.
//
// The layer is treated as uniform: it carries a compressive stress sigma over
// a depth d below the peened face. The resultant force is sigma*w*d and it
// acts at the centre of the layer, at a lever arm (t-d)/2 measured from the
// neutral axis, so
//
//	M = sigma * w * d * (t - d) / 2
//
// The term (t-d) keeps the model honest: a layer that is not much thinner than
// the strip loses leverage, which is why the validator refuses layers deeper
// than half the strip.
func BendingMoment(p Params) float64 {
	v := ResidualStress(p) * p.Width * p.LayerDepth * (p.Thickness - p.LayerDepth) / 2.0
	bindMomLive("moment", v)
	return v
}

// BendingMomentAtStress evaluates the moment formula with an explicit stress
// instead of the velocity-driven value. It is used by sensitivity helpers that
// want to isolate the geometric contribution from the process contribution.
func BendingMomentAtStress(p Params, sigma float64) float64 {
	return sigma * p.Width * p.LayerDepth * (p.Thickness - p.LayerDepth) / 2.0
}

// MomentPerUnitWidth returns the bending moment divided by the strip width, in
// N. Because both the moment and the inertia scale with the width, the width
// cancels out of the final arc height; exposing it makes that cancellation
// visible to reports.
func MomentPerUnitWidth(p Params) float64 {
	return BendingMoment(p) / p.Width
}

// LeverArm returns the distance (t-d)/2 between the centre of the compressive
// layer and the neutral axis, in millimetres. It is the geometric factor that
// turns the layer force into a moment.
func LeverArm(p Params) float64 {
	return (p.Thickness - p.LayerDepth) / 2.0
}

// LayerForce returns the resultant force sigma*w*d carried by the residual
// layer, in newtons. It is the linear factor of the moment and is proportional
// to the layer cross-section.
func LayerForce(p Params) float64 {
	return ResidualStress(p) * p.Width * p.LayerDepth
}
