package gobgp

import (
	api "github.com/alice-lg/alice-lg/backend/api"
	"github.com/alice-lg/alice-lg/backend/caches"
	gobgpapi "github.com/osrg/gobgp/api"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"

	"context"
	"fmt"
	"io"
	"log"
	"time"
)

type GoBGP struct {
	config Config
	client gobgpapi.GobgpApiClient

	// Caches: Neighbors
	neighborsCache *caches.NeighborsCache

	// Caches: Routes
	routesRequiredCache    *caches.RoutesCache
	routesReceivedCache    *caches.RoutesCache
	routesFilteredCache    *caches.RoutesCache
	routesNotExportedCache *caches.RoutesCache
}

func NewGoBGP(config Config) *GoBGP {

	dialOpts := make([]grpc.DialOption, 0)
	if config.Insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(config.TLSCert, config.TLSCommonName)
		if err != nil {
			log.Fatalf("could not load tls cert: %s", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.Dial(config.Host, dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := gobgpapi.NewGobgpApiClient(conn)

	// Cache settings:
	// TODO: Maybe read from config file
	neighborsCacheDisable := false

	routesCacheDisabled := false
	routesCacheMaxSize := 128

	// Initialize caches
	neighborsCache := caches.NewNeighborsCache(neighborsCacheDisable)
	routesRequiredCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesReceivedCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesFilteredCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesNotExportedCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)

	return &GoBGP{
		config: config,
		client: client,

		neighborsCache: neighborsCache,

		routesRequiredCache:    routesRequiredCache,
		routesReceivedCache:    routesReceivedCache,
		routesFilteredCache:    routesFilteredCache,
		routesNotExportedCache: routesNotExportedCache,
	}
}

func (gobgp *GoBGP) ExpireCaches() int {
	count := gobgp.routesRequiredCache.Expire()
	count += gobgp.routesNotExportedCache.Expire()

	return count
}

func (gobgp *GoBGP) NeighboursStatus() (*api.NeighboursStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response := api.NeighboursStatusResponse{}
	response.Neighbours = make(api.NeighboursStatus, 0)

	resp, err := gobgp.client.ListPeer(ctx, &gobgpapi.ListPeerRequest{})
	if err != nil {
		return nil, err
	}
	for {
		_resp, err := resp.Recv()
		if err == io.EOF {
			break
		}

		ns := api.NeighbourStatus{}
		ns.Id = PeerHash(_resp.Peer)

		switch _resp.Peer.State.SessionState {
		case gobgpapi.PeerState_ESTABLISHED:
			ns.State = "up"
		default:
			ns.State = "down"
		}

		if _resp.Peer.Timers.State.Uptime != nil {
			ns.Since = time.Now().Sub(time.Unix(_resp.Peer.Timers.State.Uptime.Seconds, int64(_resp.Peer.Timers.State.Uptime.Nanos)))
		}

	}
	return &response, nil
}

func (gobgp *GoBGP) Status() (*api.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := gobgp.client.GetBgp(ctx, &gobgpapi.GetBgpRequest{})
	if err != nil {
		return nil, err
	}

	response := api.StatusResponse{}
	response.Status.RouterId = resp.Global.RouterId
	response.Status.LastReboot = time.Unix(resp.Global.StartedAt.Seconds, int64(resp.Global.StartedAt.Nanos))
	response.Status.LastReconfig = time.Unix(resp.Global.ReconfiguredAt.Seconds, int64(resp.Global.ReconfiguredAt.Nanos))
	response.Status.Version = resp.Global.Version.Version
	response.Status.Message = "Daemon is up and running"
	response.Status.Backend = "gobgp"
	return &response, nil
}

func (gobgp *GoBGP) Neighbours() (*api.NeighboursResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response := api.NeighboursResponse{}
	response.Neighbours = make(api.Neighbours, 0)

	resp, err := gobgp.client.ListPeer(ctx, &gobgpapi.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil, err
	}
	for {
		_resp, err := resp.Recv()
		if err == io.EOF {
			break
		}

		neigh := api.Neighbour{}
		if _resp != nil {

			neigh.Address = _resp.Peer.State.NeighborAddress
			neigh.Asn = int(_resp.Peer.State.PeerAs)
			switch _resp.Peer.State.SessionState {
			case gobgpapi.PeerState_ESTABLISHED:
				neigh.State = "up"
			default:
				neigh.State = "down"
			}
			neigh.Description = _resp.Peer.Conf.Description

			neigh.Id = PeerHash(_resp.Peer)

			response.Neighbours = append(response.Neighbours, &neigh)
			for _, afiSafi := range _resp.Peer.AfiSafis {
				neigh.RoutesReceived += int(afiSafi.State.Received)
				neigh.RoutesExported += int(afiSafi.State.Advertised)
				neigh.RoutesAccepted += int(afiSafi.State.Accepted)
				neigh.RoutesFiltered += (neigh.RoutesReceived - neigh.RoutesAccepted)
			}

			if _resp.Peer.Timers.State.Uptime != nil {
				neigh.Uptime = time.Now().Sub(time.Unix(_resp.Peer.Timers.State.Uptime.Seconds, int64(_resp.Peer.Timers.State.Uptime.Nanos)))
			}
		}

	}

	return &response, nil
}

// Get neighbors from neighbors summary
func (gobgp *GoBGP) summaryNeighbors() (*api.NeighboursResponse, error) {

	return nil, fmt.Errorf("Not implemented summaryNeighbors")
}

// Get neighbors from protocols
func (gobgp *GoBGP) bgpProtocolsNeighbors() (*api.NeighboursResponse, error) {
	return nil, fmt.Errorf("Not implemented protocols")
}

// Get filtered and exported routes
func (gobgp *GoBGP) Routes(neighbourId string) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

/*
RoutesRequired is a specialized request to fetch:

 - RoutesExported and
 - RoutesFiltered

from Birdwatcher. As the not exported routes can be very many
these are optional and can be loaded on demand using the
RoutesNotExported() API.

A route deduplication is applied.
*/

func (gobgp *GoBGP) getRoutes(neighbourId string) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

func (gobgp *GoBGP) RoutesRequired(neighbourId string) (*api.RoutesResponse, error) {
	return gobgp.getRoutes(neighbourId)
}

// Get all received routes
func (gobgp *GoBGP) RoutesReceived(neighbourId string) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	routes.Filtered = nil
	return &routes, nil
}

// Get all filtered routes
func (gobgp *GoBGP) RoutesFiltered(neighbourId string) (*api.RoutesResponse, error) {
	routes, err := gobgp.getRoutes(neighbourId)
	if err != nil {
		log.Print(err)
	}
	routes.Imported = nil
	return routes, err
}

// Get all not exported routes
func (gobgp *GoBGP) RoutesNotExported(neighbourId string) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}
	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(neigh, gobgpapi.TableType_ADJ_OUT, &routes)
	if err != nil {
		return nil, err
	}
	routes.NotExported = routes.Filtered
	return &routes, nil
}

// Make routes lookup
func (gobgp *GoBGP) LookupPrefix(prefix string) (*api.RoutesLookupResponse, error) {
	return nil, fmt.Errorf("Not implemented LookupPrefix")
}

/*
AllRoutes:
	Here a routes dump (filtered, received) is returned, which is used to learn all prefixes to build up a local store for searching.
*/
func (gobgp *GoBGP) AllRoutes() (*api.RoutesResponse, error) {
	routes := NewRoutesResponse()
	peers, err := gobgp.GetNeighbours()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		err = gobgp.GetRoutes(peer, gobgpapi.TableType_ADJ_IN, &routes)
		if err != nil {
			log.Print(err)
		}
	}

	return &routes, nil
}
