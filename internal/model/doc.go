// Package model implements the pinned Almen bending model used by almen-int.
//
// A shot-peened strip carries a thin residual compressive layer on the peened
// face. The layer exerts a bending moment on the strip cross-section, the
// strip bends into a circular arc, and the sagitta of that arc is the Almen
// arc height. The model keeps three quantities pinned so that every report is
// reproducible:
//
//   - the residual stress follows a velocity power law with exponent 2,
//     matching the kinetic energy carried by a shot at a given velocity;
//   - the second moment of area is the rectangular value w*t^3/12, which makes
//     the arc height fall as the strip is thickened;
//   - the arc height is the small-deflection sagitta kappa*L^2/8.
//
// Coverage and saturation are deliberately not part of this package: the
// bending model only describes the fully peened plateau, and the coverage
// package decides how much of that plateau a given coverage reaches.
package model
