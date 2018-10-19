
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import {push} from 'react-router-redux'

import CommunityLabel
  from 'components/routeservers/communities/label'
import {makeReadableCommunity}
  from 'components/routeservers/communities/utils'

import {makeLinkProps, cloneFilters} from './state'

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

    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              <select className="form-control"
                      onChange={(e) => this.props.onChange(e.target.value)}
                      value={appliedFilter.value}>
                <option value="none" className="options-title">Show results from RS...</option>
                {optionsAvailable}
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
      return <option key={filter.value} value={filter.value}>
          {filter.name}, AS{filter.value} ({filter.cardinality})
        </option>;
    });

    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              <select className="form-control"
                      onChange={(e) => this.props.onChange(e.target.value)}
                      value={appliedFilter.value}> 
                <option className="options-title"
                        value="none">Show only results from AS...</option>
                {optionsAvailable}
              </select>
            </td>
          </tr>
        </tbody>
      </table>
    );
  }
}


class _CommunitiesSelect extends React.Component {
  propagateChange(value) {
    // Decode value
    const [group, community] = value.split(",", 2);
    const filterValue = community.split(":"); // spew. 

    this.props.onChange(group, filterValue);
  }

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

    const makeOption = (group, name, filter, cls) => {
      const value = `${group},${filter.value.join(":")}`; // yikes.
      return (
        <option key={filter.value} value={value} className={cls}>
          {filter.name} {name} ({filter.cardinality})
        </option>
      );
    }

    const communitiesOptions = communitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}`;
      return makeOption(FILTER_GROUP_COMMUNITIES, name, filter, cls);
    });

    const extCommunitiesOptions = extCommunitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}` +
        `select-bgp_community-2-${filter.value[2]}`;
      return makeOption(FILTER_GROUP_EXT_COMMUNITIES, name, filter, cls);
    });

    const largeCommunitiesOptions = largeCommunitiesAvailable.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      const cls = `select-bgp-community-0-${filter.value[0]} ` +
        `select-bgp-community-1-${filter.value[1]}` +
        `select-bgp_community-2-${filter.value[2]}`;
      return makeOption(FILTER_GROUP_LARGE_COMMUNITIES, name, filter, cls);
    });

    // Render list of applied communities
    const makeCommunity = (group, name, filter) => (
      <tr key={filter.value}>
        <td className="select-container">
          <CommunityLabel community={filter.value} />
        </td>
        <td>
          <button className="btn btn-remove"
                  onClick={() => this.props.onRemove(group, filter.value)}>
            <i className="fa fa-times" />
          </button>
        </td>
      </tr>
    );

    const appliedCommunities = this.props.applied.communities.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      return makeCommunity(FILTER_GROUP_COMMUNITIES, name, filter);
    }); 

    const appliedExtCommunities = this.props.applied.ext.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      return makeCommunity(FILTER_GROUP_EXT_COMMUNITIES, name, filter);
    }); 

    const appliedLargeCommunities = this.props.applied.large.map((filter) => {
      const name = makeReadableCommunity(this.props.communities, filter.value);
      return makeCommunity(FILTER_GROUP_LARGE_COMMUNITIES, name, filter);
    }); 

    return (
      <table className="select-ctrl">
        <tbody>
          {appliedCommunities}
          {appliedExtCommunities}
          {appliedLargeCommunities}
          <tr>
            <td className="select-container" colSpan="2">
              <select value="none"
                      onChange={(e) => this.propagateChange(e.target.value)}
                      className="form-control">
                <option value="none" className="options-title">
                  Select communities to match...
                </option>
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
  addFilter(group, value) {
    let nextFilters = _applyFilterValue(
      this.props.applied, group, value 
    );
    this.props.dispatch(push(
      makeLinkProps(Object.assign({}, this.props.link, {
        filtersApplied: nextFilters,
      }))
    ));
  }

  removeFilter(group, sourceId) {
    let nextFilters = _removeFilterValue(
      this.props.applied, group, sourceId
    );

    this.props.dispatch(push(
      makeLinkProps(Object.assign({}, this.props.link, {
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

