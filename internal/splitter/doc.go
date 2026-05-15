// Package splitter partitions a differ.Result into named buckets,
// making it easier to analyse large .env files by grouping related
// keys together.
//
// Two built-in strategies are provided:
//
//   - StrategyPrefix — groups keys by their underscore-delimited
//     prefix (e.g. "DB_HOST" and "DB_PORT" both land in the "DB" bucket).
//
//   - StrategyAlpha — groups keys by their first character, producing
//     simple A-Z buckets.
//
// Custom bucket limits can be applied via Options.MaxBuckets so that
// reports stay concise even when many prefixes are present.
package splitter
