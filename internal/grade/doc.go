// Package grade maps an Almen arc height to the recommended Almen strip
// designator (N, A or C) and assembles the final result of a case.
//
// Almen intensity is measured on calibrated strips. The three common strips
// are the thin N strip, the standard A strip and the thick C strip. Each strip
// has a band of arc heights it measures cleanly; the recommendation picks the
// strip whose band contains the computed arc height. The A strip covers the
// 0.10-0.60 mm band that the reference case is designed to land in.
//
// A strength grade is only meaningful when the process has saturated: the
// saturation rule decides how much of the plateau the coverage reaches, and
// until the process is saturated no grade letter is reported. The gating is
// enforced here so that every renderer reads it from one place.
package grade
