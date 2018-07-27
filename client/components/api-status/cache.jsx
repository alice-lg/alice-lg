
import React from 'react'
import {connect} from 'react-redux'

import {parseServerTime} from 'components/datetime/parse'

import moment from 'moment'

/*
 * Calculate age (generated_at), and set from_cache_status
 */
export const apiCacheStatus = function(apiStatus) {
  if (apiStatus == {}) {
    return null;
  }

  const cacheStatus = apiStatus["cache_status"] || {};
  const cachedAt = cacheStatus.cached_at;
  if (!cachedAt) {
    return null;
  }

  const fromCache = apiStatus.result_from_cache;
  const ttl = parseServerTime(apiStatus.ttl);
  const generatedAt = parseServerTime(cachedAt);
  const age = ttl.diff(generatedAt); // ms

  return {fromCache, age, generatedAt, ttl};
};

