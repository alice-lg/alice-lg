
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import {push} from 'react-router-redux'

import {cloneFilters} from 'components/filters/state'

import {FILTER_GROUP_SOURCES,
        FILTER_GROUP_ASNS,
        FILTER_GROUP_COMMUNITIES,
        FILTER_GROUP_EXT_COMMUNITIES,
        FILTER_GROUP_LARGE_COMMUNITIES}
  from './groups'

import {RouteserversSelect,
        PeersFilterSelect,
        CommunitiesSelect}
 from './widgets'

/*
 * Helper: Add and remove filter
 */
function _applyFilterValue(filters, group, value) {
  let nextFilters = cloneFilters(filters);
  nextFilters[group].filters.push({
    value: value,
  });
  return nextFilters;
}

function _removeFilterValue(filters, group, value) {
  const svalue = value.toString();
  let nextFilters = cloneFilters(filters);
  let groupFilters = nextFilters[group].filters;
  nextFilters[group].filters = _.filter(groupFilters, (f) => {
    return f.value.toString() !== svalue;
  });
  return nextFilters;
}

class FiltersEditor extends React.Component {
  addFilter(group, value) {
    let nextFilters = _applyFilterValue(
      this.props.applied, group, value
    );
    this.props.dispatch(push(
      this.props.makeLinkProps(Object.assign({}, this.props.link, {
        filtersApplied: nextFilters,
      }))
    ));
  }

  removeFilter(group, sourceId) {
    let nextFilters = _removeFilterValue(
      this.props.applied, group, sourceId
    );

    this.props.dispatch(push(
      this.props.makeLinkProps(Object.assign({}, this.props.link, {
        filtersApplied: nextFilters,
      }))
    ));
  }

  render() {
    if (!this.props.hasRoutes) {
      return null;
    }
    return (
      <div className="card lookup-filters-editor">
        <h2>Route server</h2>
        <RouteserversSelect onChange={(value) => this.addFilter(FILTER_GROUP_SOURCES, value)}
                            onRemove={(value) => this.removeFilter(FILTER_GROUP_SOURCES, value)}
                            available={this.props.availableSources}
                            applied={this.props.appliedSources} />

        <h2>Neighbor</h2>
        <PeersFilterSelect onChange={(value) => this.addFilter(FILTER_GROUP_ASNS, value)}
                           onRemove={(value) => this.removeFilter(FILTER_GROUP_ASNS, value)}
                           available={this.props.availableAsns}
                           applied={this.props.appliedAsns} />

        <h2>Communities</h2>
        <CommunitiesSelect onChange={(group, value) => this.addFilter(group, value)}
                           onRemove={(group, value) => this.removeFilter(group, value)}
                           available={this.props.availableCommunities}
                           applied={this.props.appliedCommunities} />

      </div>
    );
  }
}

export default connect(
  (state) => ({
    isLoading: state.lookup.isLoading,
    hasRoutes: state.lookup.routesFiltered.length > 0 ||
               state.lookup.routesImported.length > 0,

    link: {
      pageReceived:   0, // Reset pagination on filter change
      pageFiltered:   0,
      query:          state.lookup.query,
      filtersApplied: state.lookup.filtersApplied,
      routing:        state.routing.locationBeforeTransitions,
    },

    available: state.lookup.filtersAvailable,
    applied: state.lookup.filtersApplied,

    availableSources: state.lookup.filtersAvailable[FILTER_GROUP_SOURCES].filters,
    appliedSources:   state.lookup.filtersApplied[FILTER_GROUP_SOURCES].filters,

    availableAsns: state.lookup.filtersAvailable[FILTER_GROUP_ASNS].filters,
    appliedAsns:   state.lookup.filtersApplied[FILTER_GROUP_ASNS].filters,

    availableCommunities: {
      communities: state.lookup.filtersAvailable[FILTER_GROUP_COMMUNITIES].filters,
      ext:         state.lookup.filtersAvailable[FILTER_GROUP_EXT_COMMUNITIES].filters,
      large:       state.lookup.filtersAvailable[FILTER_GROUP_LARGE_COMMUNITIES].filters,
    },
    appliedCommunities: {
      communities: state.lookup.filtersApplied[FILTER_GROUP_COMMUNITIES].filters,
      ext:         state.lookup.filtersApplied[FILTER_GROUP_EXT_COMMUNITIES].filters,
      large:       state.lookup.filtersApplied[FILTER_GROUP_LARGE_COMMUNITIES].filters,
    },

  })

)(FiltersEditor);

