package api

import (
	"strconv"
	"strings"
)

// FilterQueryParser parses a filter value into a search filter
type FilterQueryParser func(value string) (*SearchFilter, error)

func parseQueryValueList(parser FilterQueryParser, value string) ([]*SearchFilter, error) {
	components := strings.Split(value, ",")
	result := make([]*SearchFilter, 0, len(components))

	for _, component := range components {
		filter, err := parser(strings.TrimSpace(component))
		if err != nil {
			return result, err
		}
		result = append(result, filter)
	}

	return result, nil
}

func parseIntValue(value string) (*SearchFilter, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	return &SearchFilter{
		Value: v,
	}, nil
}

func parseStringValue(value string) (*SearchFilter, error) {
	return &SearchFilter{
		Value: value,
	}, nil
}

func parseCommunityValue(value string) (*SearchFilter, error) {
	components := strings.Split(value, ":")
	community := make(Community, len(components))

	for i, c := range components {
		v, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		community[i] = v
	}

	return &SearchFilter{
		Name:  community.String(),
		Value: community,
	}, nil
}

func parseExtCommunityValue(value string) (*SearchFilter, error) {
	components := strings.Split(value, ":")
	community := make(ExtCommunity, len(components))

	for i, c := range components {
		community[i] = c
	}

	return &SearchFilter{
		Name:  community.String(),
		Value: community,
	}, nil
}
