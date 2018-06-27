

# Feature: Display Routes Filtered

Merged commits:
    
    610ab3b002b1 - Modify parseNeighbours for Multi-Table style counters.
    22e9f59ff7c5 - Modify client to display Routes Accepted in Neighbours view.
    a3da2600fbfa - Fix values displayed in client.
    d0ca03fc55ad - Fix accepted routes display in case of multiple routers.
    bd0dbd1cf361 - Add suggestions on semantics of received column
    7db8bf28567f - Refactor route counts on protocols page

    Removed Feature: Next Hop Link
    0d9b2ab41231 - This reverts commit 15e728da2c6855a3fad6a22a58dbd6d62456a7cb.
    b7181a4b69b8 - Forgot to clean them up while reverting the nexthop feature

Todo:
    
    This needs to be made configurable!
   
    Use RoutesFiltered / Accepted from Birdwatcher;
    Make configurable using:
        - Caclulated strategy from DECIX
        - Direct values from Birdwatcher


# Feature: Display routes exported

Merged commits
    
    ea76d9390b3 - Add column for exported routes in NeighboursTable. 

Todo:
    
    Configurable columns on Neighbours table


# Feature: UI Tweaks / Improvements
    
Merged commits:

    7902d7c4e815 - Add router-selective display on RoutesPage.
    a7260aad7e1d - Remove filter for nextHop in column 'Description'.
    056637eef7a4 - Change search box hint string in routes view
    f840140dfff9 - Change search box hint string on splash page
    9ee4d08c65cc - Change some GUI text
    ab104cf43f30 - Show modal dialog when clicking on prefix.


# Feature: IRRExplorer Link

Merged Commits:

    a0751d4c0cac - Add Link to IRRExplorer for prefixes.
    f15496f6cfe4 - Add ASN links to IRRExplorer. 
    1d7d376b987d - Change IRRDB links, make export reasons clickable

Todo:

    ASN links to IRRExplorer is a bit hacky, see if
    we can improve this!



