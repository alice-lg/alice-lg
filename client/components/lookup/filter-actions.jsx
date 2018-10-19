
export const APPLY_FILTER = "@lookup/APPLY_FILTER";

export function applyFilterValue(group, value) {
  return {
    type: APPLY_FILTER,
    payload: {
      group: group,
      filter: {
        value: value,
      },
    },
  };
}

export function applyFilter(group, filter) {
  return {
    type: APPLY_FILTER,
    payload: {
      group: group,
      filter: filter,
    }
  };
};

