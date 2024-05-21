// Package cachex provides a unified cache interface definition, LRU Cache, File Cache, and Cache Chain.
//
// The provided Cache interface definition and implementations based on generics make it more convenient to use.
//
// Currently, implementations based on this Cache Interface include:
// 1. LRUCache   : In-memory LRU cache on a single machine
// 2. FileStore    : File-based cache on a single machine
// 3. Chain        : Integrates multiple caches together, such as local LRU and Redis multi-level cache, for higher performance and easier use
// 4. NoCache      : Empty cache, always returns no result on query, and always succeeds on write
//
// FetcherOne and FetcherMulti provide unified encapsulation for querying caches and performing origin fetches with cache writebacks.
// Examples are provided below.
package cachex
