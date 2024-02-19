// Package sources provides the base interface for all
// route server data source implementations.
package sources

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Source Errors
var (
	// SourceNotFound indicates that a source could
	// not be resolved by an identifier.
	ErrSourceNotFound = errors.New("route server unknown")

	// ErrSourceBusy is returned when a refresh is
	// already in progress.
	ErrSourceBusy = errors.New("source is busy")
)

// Source is a generic datasource for alice.
// All route server Source adapters implement this interface.
type Source interface {
	ExpireCaches() int

	Status(context.Context) (*api.StatusResponse, error)
	Neighbors(context.Context) (*api.NeighborsResponse, error)
	NeighborsSummary(context.Context) (*api.NeighborsResponse, error)
	NeighborsStatus(context.Context) (*api.NeighborsStatusResponse, error)
	Routes(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesReceived(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesFiltered(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesNotExported(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	AllRoutes(context.Context) (*api.RoutesResponse, error)
}

func hiddenNeighborsExcludeLists(hiddenNeighbors []string) ([]*net.IPNet, []net.IP, []*regexp.Regexp, error) {
	excludeCIDRs := make([]*net.IPNet, 0, len(hiddenNeighbors))
	excludeIPs := make([]net.IP, 0, len(hiddenNeighbors))
	excludePatterns := make([]*regexp.Regexp, 0, len(hiddenNeighbors))
	for _, hiddenNeighbor := range hiddenNeighbors {
		if _, hiddenNet, err := net.ParseCIDR(hiddenNeighbor); err == nil {
			excludeCIDRs = append(excludeCIDRs, hiddenNet)
		} else if ip := net.ParseIP(hiddenNeighbor); ip != nil {
			excludeIPs = append(excludeIPs, ip)
		} else {
			pattern, err := regexp.Compile(hiddenNeighbor)
			if err != nil {
				return nil, nil, nil, err
			}
			excludePatterns = append(excludePatterns, pattern)
		}
	}
	return excludeCIDRs, excludeIPs, excludePatterns, nil
}

func FilterHiddenNeighbors(neighbors api.Neighbors, hiddenNeighbors []string) (api.Neighbors, error) {
	if len(hiddenNeighbors) > 0 {
		filteredNeighbors := make(api.Neighbors, 0, len(neighbors))
		excludeCIDRs, excludeIPs, excludePatterns, err := hiddenNeighborsExcludeLists(hiddenNeighbors)
		if err != nil {
			return nil, err
		}
	neighbors:
		for _, neighbor := range neighbors {
			neighborIP := net.ParseIP(neighbor.Address)
			if neighborIP == nil {
				return nil, fmt.Errorf("Neighbor ID '%s' is not parseable as an IP", neighborIP)
			}
			if neighborIP != nil {
				for _, ip := range excludeIPs {
					if neighborIP.Equal(ip) {
						continue neighbors
					}
				}
				for _, cidr := range excludeCIDRs {
					if cidr.Contains(neighborIP) {
						continue neighbors
					}
				}
			}
			for _, pattern := range excludePatterns {
				if pattern.MatchString(neighbor.Address) {
					continue neighbors
				}
			}
			filteredNeighbors = append(filteredNeighbors, neighbor)
		}
		return filteredNeighbors, nil
	} else {
		return neighbors, nil
	}
}

func FilterHiddenNeighborsStatus(neighborsStatus api.NeighborsStatus, hiddenNeighbors []string) (api.NeighborsStatus, error) {
	if len(hiddenNeighbors) > 0 {
		filteredNeighborsStatus := make(api.NeighborsStatus, 0, len(neighborsStatus))
		excludeCIDRs, excludeIPs, excludePatterns, err := hiddenNeighborsExcludeLists(hiddenNeighbors)
		if err != nil {
			return neighborsStatus, err
		}
		neighbors:
		for _, neighborStatus := range neighborsStatus {
			neighborIDAsIP := net.ParseIP(neighborStatus.ID)
			if neighborIDAsIP != nil {
				for _, ip := range excludeIPs {
					if ip.Equal(neighborIDAsIP) {
						continue neighbors
					}
				}
				for _, cidr := range excludeCIDRs {
					if cidr.Contains(neighborIDAsIP) {
						continue neighbors
					}
				}
			}
			for _, pattern := range excludePatterns {
				if pattern.MatchString(neighborStatus.ID) {
					continue neighbors
				}
			}
			filteredNeighborsStatus = append(filteredNeighborsStatus, neighborStatus)
		}
		return filteredNeighborsStatus, nil
	} else {
		return neighborsStatus, nil
	}
}
