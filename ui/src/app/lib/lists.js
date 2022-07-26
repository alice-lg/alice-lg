
/** 
 * Intersect lists: [x | x <- A, x `elem` B]
 */
export function intersect(a, b) {
  let res = [];
  for (const e of a) {
    for (const k of b) {
      if (e===k) {
        res.push(e);
        break;
      }
    }
  }
  return res;
}

/**
 * Resolve list with dict: [dict[x] or x | x <- L]
 */
export function resolve(dict, list) {
  let result = [];
  for (const e of list) {
    result.push(dict[e]||e); 
  }
  return result;
}

