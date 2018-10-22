
import {humanizedJoin} from 'components/utils/text'
import {intersect, resolve} from 'components/utils/lists'

const filterable = [
  "gateway", "network"
];

export function filterableColumns(columns, order) {
  return resolve(columns, intersect(order, filterable));
}

export function filterableColumnsText(columns, order) {
  return humanizedJoin(filterableColumns(columns, order), "or");
}


