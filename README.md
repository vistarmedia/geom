# Geometry Engine

Provides thread and memory safe access for Go programs to the
[libgeos](https://trac.osgeo.org/geos/) engine. Requires libgeos 3.5.0 or
greater. All exported package functions and objects can freely be used across
goroutines and will be managed by the GC. This library deals solely with planar
geometry and is not concerned with projections or coordinate systems.

The main entry point for constructing objects from this package is through the
geom/context package.
