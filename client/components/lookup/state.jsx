

/*
 * Manage state
 */

export function lookupStateUrlEncode(state) {
    const pageImported = state.pageImported;
    const pageFiltered = state.pageFiltered;
    const filters = state.filtersApplied;
    const q = `q=${query}`;
    const f = filtersUrlEncode(filters); 

    let p = "";
    if (pageFiltered > 0) {
      p += `&page_filtered=${pageFiltered}`;
    }
    if (pageImported > 0) {
      p += `&page_imported=${pageImported}`;
    }

    return `${q}${f}${p}`;
}

