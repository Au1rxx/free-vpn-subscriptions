// Package node defines the canonical in-memory representation of a proxy node
// used across the Au1rxx tooling suite: this repository aggregates and verifies
// nodes, and sibling tools (e.g. Au1rxx/proxykit) reuse this package to parse,
// convert, and rank the same node structures without duplicating the data
// model.
//
// Stability
//
// This package is public but pre-1.0. The Node struct and protocol constants
// are the stable surface: fields may be added, but existing field names,
// types, and semantics will not change without a deprecation cycle. Helper
// functions in parse.go are less stable and may be renamed.
//
// Do not import internal/* from this repository in external projects — those
// packages remain internal to the aggregator and may change freely. If you
// need functionality that currently lives under internal/*, open an issue
// requesting promotion to pkg/.
package node
