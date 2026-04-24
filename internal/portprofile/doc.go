// Package portprofile assembles a composite metadata profile for a network
// port by combining information from multiple independent providers: label,
// rank, classification, and scope.
//
// # Overview
//
// A Profiler is constructed with New and then configured with provider
// functions via the With* fluent methods. Calling Build(port) invokes each
// provider and returns a Profile struct containing all metadata plus any
// automatically generated notes (e.g. warnings for critical-rank ports).
//
// # Usage
//
//	profiler := portprofile.New().
//	    WithLabeler(labelMap.Label).
//	    WithRanker(rankMap.Get).
//	    WithClasser(classifier.Classify).
//	    WithScoper(scopeRegistry.Scope)
//
//	profile := profiler.Build(443)
//	fmt.Println(profile.Label, profile.Rank)
//
// # Config
//
// Config.Apply wraps a Profiler so that empty scope or rank fields fall back
// to configurable defaults instead of returning empty strings.
package portprofile
