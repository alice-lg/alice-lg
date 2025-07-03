package gobgp

import (
	"fmt"

	gobgpapi "github.com/osrg/gobgp/v3/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/caches"

	"context"
	"io"
	"log"
	"time"
)

// GoBGP is a source for Alice.
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

// NewGoBGP creates a new GoBGP source instance
func NewGoBGP(config Config) *GoBGP {
	dialOpts := make([]grpc.DialOption, 0)
	if config.Insecure {
		// dialOpts = append(dialOpts, grpc.WithInsecure())
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds, err := credentials.NewClientTLSFromFile(
			config.TLSCert,
			config.TLSCommonName)
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

// ExpireCaches clears all local caches
func (gobgp *GoBGP) ExpireCaches() int {
	count := gobgp.routesRequiredCache.Expire()
	count += gobgp.routesNotExportedCache.Expire()
	return count
}

// NeighborsStatus retrieves all status information
// for all peers on the RS.
func (gobgp *GoBGP) NeighborsStatus(
	ctx context.Context,
) (*api.NeighborsStatusResponse, error) {
	ctx, cancel := context.WithTimeout(
		ctx, time.Second*time.Duration(gobgp.config.ProcessingTimeout))
	defer cancel()

	response := api.NeighborsStatusResponse{}
	response.Neighbors = make(api.NeighborsStatus, 0)

	resp, err := gobgp.client.ListPeer(ctx, &gobgpapi.ListPeerRequest{})
	if err != nil {
		return nil, err
	}
	for {
		_resp, err := resp.Recv()
		if err == io.EOF {
			break
		}

		ns := api.NeighborStatus{}
		ns.ID = PeerHash(_resp.Peer)

		switch _resp.Peer.State.SessionState {
		case gobgpapi.PeerState_ESTABLISHED:
			ns.State = "up"
		default:
			ns.State = "down"
		}

		if _resp.Peer.Timers.State.Uptime != nil {
			ns.Since = time.Since(time.Unix(
				_resp.Peer.Timers.State.Uptime.Seconds,
				int64(_resp.Peer.Timers.State.Uptime.Nanos)))
		}
		response.Neighbors = append(response.Neighbors, &ns)

	}
	return &response, nil
}

// Status retrieves the routers status
func (gobgp *GoBGP) Status(
	ctx context.Context,
) (*api.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(
		ctx, time.Second*time.Duration(gobgp.config.ProcessingTimeout))
	defer cancel()

	resp, err := gobgp.client.GetBgp(ctx, &gobgpapi.GetBgpRequest{})
	if err != nil {
		return nil, err
	}

	response := api.StatusResponse{}
	response.Status.RouterID = resp.Global.GetRouterId()
	response.Status.Backend = "gobgp"
	response.Status.ServerTime = time.Now()
	// TODO: What get last reconfig time from gobgp?
	response.Status.LastReconfig = time.Now()
	response.Status.Message = fmt.Sprintf("AS%d", resp.Global.GetAsn())
	// TODO: What get version gobgp?
	response.Status.Version = ""
	return &response, nil
}

// Neighbors retrieves a list of neighbors
func (gobgp *GoBGP) Neighbors(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	ctx, cancel := context.WithTimeout(
		ctx, time.Second*time.Duration(gobgp.config.ProcessingTimeout))
	defer cancel()

	response := api.NeighborsResponse{}
	response.Neighbors = make(api.Neighbors, 0)

	resp, err := gobgp.client.ListPeer(ctx, &gobgpapi.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil, err
	}
	for {
		_resp, err := resp.Recv()
		if err == io.EOF {
			break
		}

		neigh := api.Neighbor{}

		neigh.Address = _resp.Peer.State.NeighborAddress
		neigh.ASN = int(_resp.Peer.State.PeerAsn)
		switch _resp.Peer.State.SessionState {
		case gobgpapi.PeerState_ESTABLISHED:
			neigh.State = "up"
		default:
			neigh.State = "down"
		}
		neigh.Description = _resp.Peer.Conf.Description

		neigh.ID = PeerHash(_resp.Peer)
		neigh.RouteServerID = gobgp.config.ID

		response.Neighbors = append(response.Neighbors, &neigh)
		for _, afiSafi := range _resp.Peer.AfiSafis {
			neigh.RoutesReceived += int(afiSafi.State.Received)
			neigh.RoutesExported += int(afiSafi.State.Advertised)
			neigh.RoutesAccepted += int(afiSafi.State.Accepted)
			neigh.RoutesFiltered += (neigh.RoutesReceived - neigh.RoutesAccepted)
		}

		if _resp.Peer.Timers.State.Uptime != nil {
			neigh.Uptime = time.Since(time.Unix(
				_resp.Peer.Timers.State.Uptime.Seconds,
				int64(_resp.Peer.Timers.State.Uptime.Nanos)))
		}

	}

	return &response, nil
}

// NeighborsSummary is an alias of Neighbors for now
func (gobgp *GoBGP) NeighborsSummary(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	return gobgp.Neighbors(ctx)
}

// Routes retrieves filtered and exported routes
func (gobgp *GoBGP) Routes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbor(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(ctx, neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

func (gobgp *GoBGP) getRoutes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbor(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(ctx, neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

// RoutesRequired is a specialized request to fetch:
//
// - RoutesExported and
// - RoutesFiltered
//
// from Birdwatcher. As the not exported routes can be very many
// these are optional and can be loaded on demand using the
// RoutesNotExported() API.
//
// A route deduplication is applied.
func (gobgp *GoBGP) RoutesRequired(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	return gobgp.getRoutes(ctx, neighborID)
}

// RoutesReceived gets all received routes
func (gobgp *GoBGP) RoutesReceived(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbor(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(ctx, neigh, gobgpapi.TableType_ADJ_IN, &routes)
	if err != nil {
		return nil, err
	}
	routes.Filtered = nil
	return &routes, nil
}

// RoutesFiltered gets all filtered routes
func (gobgp *GoBGP) RoutesFiltered(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	routes, err := gobgp.getRoutes(ctx, neighborID)
	if err != nil {
		log.Print(err)
	}
	routes.Imported = nil
	return routes, err
}

// RoutesNotExported gets all not exported routes
func (gobgp *GoBGP) RoutesNotExported(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbor(ctx, neighborID)
	if err != nil {
		return nil, err
	}
	routes := NewRoutesResponse()
	err = gobgp.GetRoutes(ctx, neigh, gobgpapi.TableType_ADJ_OUT, &routes)
	if err != nil {
		return nil, err
	}
	routes.NotExported = routes.Filtered
	return &routes, nil
}

// AllRoutes returns a routes dump (filtered, received),
// which is used to learn all prefixes to build
// up a local store for searching.
func (gobgp *GoBGP) AllRoutes(
	ctx context.Context,
) (*api.RoutesResponse, error) {
	routes := NewRoutesResponse()
	peers, err := gobgp.GetNeighbors(ctx)
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		err = gobgp.GetRoutes(ctx, peer, gobgpapi.TableType_ADJ_IN, &routes)
		if err != nil {
			log.Print(err)
		}
	}
	return &routes, nil
}
