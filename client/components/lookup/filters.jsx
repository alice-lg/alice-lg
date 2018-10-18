
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import {push} from 'react-router-redux'

import {makeReadableCommunity}
  from 'components/routeservers/communities/utils'

import {makeLinkProps} from './state'

import {FILTER_GROUP_SOURCES,
        FILTER_GROUP_ASNS,
        FILTER_GROUP_COMMUNITIES,
        FILTER_GROUP_EXT_COMMUNITIES,
        FILTER_GROUP_LARGE_COMMUNITIES}
  from './filter-groups'


/*
 * Helper: Add and remove filter
 */
function _applyFilterValue(filters, group, value) {
  let nextFilters = Object.assign([], filters);
  nextFilters[group].filters.push({
    value: value,
  });
  return nextFilters;
}

function _removeFilterValue(filters, group, value) {
  const svalue = value.toString();
  let nextFilters = Object.assign([], filters);
  let groupFilters = nextFilters[group].filters;
  nextFilters[group].filters = _.filter(groupFilters, (f) => {
    return f.value.toString() !== svalue;
  });
  return nextFilters;
}


class RouteserversSelect extends React.Component {
  render() {
    // Sort filters available
    const sortedFiltersAvailable = this.props.available.sort((a, b) => {
      return a.value - b.value;
    });

    // For now we allow only one applied
    const appliedFilter = this.props.applied[0] || {value: undefined};

    if (appliedFilter.value !== undefined) {
      const filter = _.findWhere(sortedFiltersAvailable, {
        value: appliedFilter.value
      });

      // Just render this, with a button for removal
      return (
        <table className="select-ctrl">
          <tbody>
            <tr>
              <td className="select-container">
                {filter.name}
              </td>
              <td>
                <button className="btn btn-remove"
                        onClick={() => this.props.onRemove(filter.value)}>
                  <i className="fa fa-times" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      );
    }

    // Build options
    const optionsAvailable = sortedFiltersAvailable.map((filter) => {
      return <option key={filter.value} value={filter.value}>
          {filter.name} ({filter.cardinality})
        </option>;
    });

    let options = [
      <option key="none">Show results from RS...</option>
    ];
    options = options.concat(optionsAvailable);

    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              <select className="form-control"
                      onChange={(e) => this.props.onChange(e.target.value)}
                      value={appliedFilter.value}>
                {options}
              </select>
            </td>
          </tr>
        </tbody>
      </table>
    );
  }
}


class PeersFilterSelect extends React.Component {
  render() {
    // Sort filters available
    const sortedFiltersAvailable = this.props.available.sort((a, b) => {
      return a.name.localeCompare(b.name);
    });

    // For now we allow only one applied
    const appliedFilter = this.props.applied[0] || {value: undefined};

    if (appliedFilter.value !== undefined) {
      const filter = _.findWhere(sortedFiltersAvailable, {
        value: appliedFilter.value
      });

      // Just render this, with a button for removal
      return (
        <table className="select-ctrl">
          <tbody>
            <tr>
              <td className="select-container">
                {filter.name}
              </td>
              <td>
                <button className="btn btn-remove"
                        onClick={() => this.props.onRemove(filter.value)}>
                  <i className="fa fa-times" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      );
    }

    // Build options
    const optionsAvailable = sortedFiltersAvailable.map((filter) => {
      return <option key={filter.value} value={filter}>
          {filter.name} AS{filter.value} ({filter.cardinality})
        </option>;
    });

    let options = [];
    if (!appliedFilter.value) {
      options = [
        <option key="none">Show only results from AS...</option>
      ];
    }
    options = options.concat(optionsAvailable);

    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              <select className="form-control"
                      onChange={(e) => this.props.onChange(e.target.value)}
                      value={appliedFilter.value}> 
                {options}
              </select>
            </td>
          </tr>
        </tbody>
      </table>
    );
  }
}


class _CommunitiesSelect extends React.Component {
  render() {
    const communitiesAvailable = this.props.available.communities.sort((a, b) => {
      return (a.value[0] - b.value[0]) * 100000 + (a.value[1] - b.value[1]);
    });

    const extCommunitiesAvailable = this.props.available.ext.sort((a, b) => {
      return (a.value[1] - b.value[1]) * 100000 + (a.value[2] - b.value[2]);
    });

    const largeCommunitiesAvailable = this.props.available.large.sort((a, b) => {
      return (a.value[0] - b.value[0]) * 10000000000 +
             (a.value[1] - b.value[1]) * 100000 +
             (a.value[2] - b.value[2]);
    });

    const communitiesOptions = communitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}`;
      return (
        <option key={filter.value} value={filter} className={cls}>
          {filter.name} {name} ({filter.cardinality})
        </option>
      );
    });

    const extCommunitiesOptions = extCommunitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}` +
        `select-bgp_community-2-${filter.value[2]}`;
      return (
        <option key={filter.value} value={filter} className={cls}>
          {filter.name} {name} ({filter.cardinality})
        </option>
      );
    });

    const largeCommunitiesOptions = largeCommunitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}` +
        `select-bgp_community-2-${filter.value[2]}`;
      return (
        <option key={filter.value} value={filter} className={cls}>
          {filter.name} {name} ({filter.cardinality})
        </option>
      );
    });


    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              <select className="form-control">
                {communitiesOptions.length > 0 &&
                  <optgroup label="Communities">
                    {communitiesOptions} 
                  </optgroup>}

                {extCommunitiesOptions.length > 0 &&
                  <optgroup label="Ext. Communities">
                    {extCommunitiesOptions} 
                  </optgroup>}

                {largeCommunitiesOptions.length > 0 &&
                  <optgroup label="Large Communities">
                    {largeCommunitiesOptions}
                  </optgroup>}
              </select>
            </td>
            <td>
              <button className="btn">
                <i className="fa fa-plus"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    );
  }
}

const CommunitiesSelect = connect(
  (state) => ({
    communities: state.config.bgp_communities,
  })
)(_CommunitiesSelect);


class FiltersEditor extends React.Component {
  selectSource(sourceId) {
    let nextFilters = _applyFilterValue(
      this.props.applied, FILTER_GROUP_SOURCES, sourceId
    );

    this.props.dispatch(push(
      /*
      makeLinkProps(Object.assign({}, this.props.link, {
        filtersApplied: nextFilters,
      }))
      */
    ));
  }

  render() {
    if (!this.props.hasRoutes) {
      return null;
    }
    return (
      <div className="card lookup-filters-editor">
        <h2>Route server</h2>
        <RouteserversSelect onChange={(value) => this.selectSource(value)}
                            available={this.props.availableSources}
                            applied={this.props.appliedSources} />

        <h2>Neighbor</h2>
        <PeersFilterSelect available={this.props.availableAsns}
                           applied={this.props.appliedAsns} />

        <h2>Communities</h2>
        <CommunitiesSelect available={this.props.availableCommunities}
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
      pageReceived:   state.lookup.pageReceived,
      pageFiltered:   state.lookup.pageFiltered,
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

