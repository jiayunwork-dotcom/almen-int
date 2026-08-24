// Package coverage implements the coverage law and the saturation rule of the
// Almen process.
//
// The coverage fraction follows the standard exponential accumulation law
//
//	C(t) = 1 - exp(-lambda * t)
//
// where lambda is the coverage rate constant of the peening setup and t is
// the peening time. Given a target coverage the package recovers the implied
// time, evaluates how much of the fully peened arc height a partial coverage
// reaches, and decides whether the process is saturated.
//
// Saturation is pinned to a single rule: doubling the peening time must not
// raise the arc height by more than ten percent. The comparison is made on the
// arc-height gain, so the decision depends on the coverage and the gain curve
// but is independent of the shot velocity, the strip size and the elastic
// modulus. This separation is intentional: coverage enters the final report
// only through the gain and the saturation flag, never through the bending
// moment itself.
package coverage
