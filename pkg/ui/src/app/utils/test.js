
/**
 * Join list of words with ',' and provide a glue
 * for the last element.
 *
 * Example:
 *   humanizedJoin(["foo", "bar", "baz"], "or") ->
 *   "foo, bar or baz"
 */
export function humanizedJoin(list, glue="and") {
  // Doing this the other way round in one step would be nice.
  let [last, ...init] = list.reverse();
  init = init.reverse();
  if (init.length === 0) {
    return last;
  }
  return init.join(", ") + ` ${glue} ${last}`; 
}

