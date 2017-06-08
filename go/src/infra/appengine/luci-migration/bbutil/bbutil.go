// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package bbutil contains utility functions and constants for Buildbucket.
package bbutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/luci/luci-go/common/api/buildbucket/buildbucket/v1"
)

const (
	// StatusCompleted means a build has reached its terminal state.
	StatusCompleted = "COMPLETED"
	// ResultFailure means a build has failed, with or without an infra-failure.
	ResultFailure = "FAILURE"
	// ResultSuccess means a build succeeded.
	ResultSuccess = "SUCCESS"
	// TagBuildSet is a key of a tag used to group builds together.
	// For tryjobs, it specifies the patchset.
	TagBuildSet = "buildset"
)

// ParseTag parses a buildbucket tag.
//
// If tag does not have ":", the whole tag becomes the key with an empty
// value.
func ParseTag(tag string) (k, v string) {
	parts := strings.SplitN(tag, ":", 2)
	k = parts[0]
	if len(parts) > 1 {
		v = parts[1]
	} else {
		// this tag is invalid. This should not happen in practice.
		// Do not panic because this function is used for externally-supplied
		// data.
	}
	return
}

// FormatTag formats a tag from a key-value pair.
func FormatTag(k, v string) string {
	return k + ":" + v
}

// ParseTimestamp parses a buildbucket timestamp.
func ParseTimestamp(ts int64) time.Time {
	if ts == 0 {
		return time.Time{}
	}
	return time.Unix(ts/1000000, 0)
}

// FormatTimestamp t converts to a buildbucket timestamp.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1000
}

// Duration returns duration from build creation to build compeltion.
func Duration(b *buildbucket.ApiCommonBuildMessage) time.Duration {
	return ParseTimestamp(b.CompletedTs).Sub(ParseTimestamp(b.CreatedTs))
}

// BuildSet returns the value of the "buildset" tag in b, or "" if not found.
func BuildSet(b *buildbucket.ApiCommonBuildMessage) string {
	for _, t := range b.Tags {
		if k, v := ParseTag(t); k == TagBuildSet {
			return v
		}
	}

	return ""
}

// BuildSetURL converts a buildSet to a URL, if possible. Otherwise returns "".
func BuildSetURL(buildSet string) string {
	parts := strings.Split(buildSet, "/")
	if len(parts) >= 5 && parts[0] == "patch" {
		switch parts[1] {
		case "rietveld":
			return fmt.Sprintf("https://%s/%s/#ps%s", parts[2], parts[3], parts[4])
		case "gerrit":
			return fmt.Sprintf("https://%s/c/%s/%s", parts[2], parts[3], parts[4])
		}
	}

	return ""
}
