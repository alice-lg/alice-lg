
--
-- ----------------------
-- AliceLG schema v.1.0.0
-- ----------------------
--
-- %% Author:      annika
-- %% Description: Apply alice-lg db schema.
--

-- Clear state
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS neighbors;
DROP TABLE IF EXISTS __meta__;

-- Neighbors
CREATE TABLE neighbors (
    id    VARCHAR(255) NOT NULL PRIMARY KEY,

    -- Indexed attributes
    rs_id VARCHAR(255) NOT NULL,

    -- JSON serialized neighbor
    neighbor     jsonb NOT NULL,

    -- Timestamps
    updated_at  TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_neighbors_rs_id 
          ON neighbors            USING HASH (rs_id);

-- Routes
CREATE TABLE routes (
    id            VARCHAR(255) NOT NULL PRIMARY KEY,
    neighbor_id   VARCHAR(255) NOT NULL
                  REFERENCES   neighbors(id)
                  ON DELETE    CASCADE,

    -- Indexed attributes 
    network       cidr         NOT NULL,
   
    -- JSON serialized route
    route         jsonb        NOT NULL,

    -- Timestamps
    updated_at  TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_routes_network ON routes ( network );


-- The meta table stores information about the schema
-- like when it was migrated and the current revision.
CREATE TABLE __meta__ (
    version     INTEGER   NOT NULL  UNIQUE,
    description TEXT      NOT NULL,
    applied_at  TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO __meta__ (version, description)
     VALUES (1, 'initial schema');

