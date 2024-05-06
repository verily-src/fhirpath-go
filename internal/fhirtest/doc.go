/*
Package fhirtest provides resource test dummies and useful utilities to enable
better testing of the R4 FHIR Protos.

This provides:
  - Pseudo-randomized FHIR resource identity generation
  - Construction utilities for forming new resources at runtime
  - Utilities for emulating Meta updates to FHIR resources
  - Resources are organized by their higher-level interface abstractions (e.g.
    organized by Resource, DomainResource, etc), and are keyed by resource-name.
*/
package fhirtest
