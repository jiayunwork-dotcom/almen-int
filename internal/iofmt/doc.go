// Package iofmt loads Almen case files, validates the document shape, and
// renders the computed result as plain text or as JSON.
//
// A case file is a single JSON document with four nested objects: the shot
// stream, the strip geometry, the residual layer, and the process settings,
// plus a reference velocity that anchors the pinned power law. Every numeric
// field is required; a missing or unparseable field produces a list of issues
// that the CLI prints on stderr. The renderer never computes anything on its
// own - it only formats values produced by the model, coverage and grade
// packages.
package iofmt
