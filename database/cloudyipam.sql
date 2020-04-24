-- Extensions

-- Required for UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables

CREATE TABLE IF NOT EXISTS zone (
  id            UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          VARCHAR       UNIQUE NOT NULL,
  range         CIDR          UNIQUE NOT NULL,
  prefixlen     INT           NOT NULL
);

CREATE TABLE IF NOT EXISTS subnet (
  id            UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
  zone          UUID          REFERENCES zone(id),
  range         CIDR          NOT NULL,
  available     BOOLEAN       NOT NULL,
  usage         VARCHAR       NOT NULL
);

-- API

CREATE OR REPLACE FUNCTION create_zone(_name VARCHAR, _range CIDR, _prefixlen INT)
  RETURNS UUID
  SECURITY DEFINER
AS $$
  DECLARE
    newid UUID;
  BEGIN
    INSERT INTO zone (name,range,prefixlen) VALUES (_name, _range, _prefixlen) RETURNING id INTO newid;
    CALL populate_zone(newid);
    RETURN newid;
  END;
$$ LANGUAGE plpgsql;

-- not that useful, but necessary for Terraform provider implementation
CREATE OR REPLACE FUNCTION read_zone(_id UUID)
  RETURNS SETOF zone
  SECURITY DEFINER
AS $$
  BEGIN
    RETURN QUERY SELECT * FROM zone WHERE id = _id;
  END
$$ LANGUAGE plpgsql;

-- requires all subnets to be destroyed first. This is enforced via referential
-- integrity. As this is intended to be an API used by applications (and
-- particularly a Terraform provider), it can be considered a bug if zone
-- teardown is not done in the appropriate order:
--
-- 1. deallocate all allocated subnets in the zone. After this step, all
--    subnets in the zone should be available
-- 2. destroy the zone
--
-- Deleting the available subnets is performed here rather than in a separate
-- procedure because this way it will always be transaction-safe and the system
-- should never be able to get into an inconsistent state.
CREATE OR REPLACE PROCEDURE destroy_zone(_id UUID)
  SECURITY DEFINER
AS $$
  BEGIN
    DELETE FROM subnet WHERE zone = _id AND available = TRUE;
    DELETE FROM zone WHERE id = _id;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE populate_zone(_id UUID)
  SECURITY DEFINER
AS $$
  DECLARE
    zonecidr  CIDR;
    plen      INT;
    incr      INT;
    nincr     INT;
  BEGIN
    SELECT range, prefixlen FROM zone WHERE id = _id INTO zonecidr, plen;
    IF NOT FOUND THEN
      RAISE EXCEPTION 'zone % not found.', _id;
    END IF;
    IF masklen(zonecidr) > plen THEN
      RAISE EXCEPTION 'zone % (%) not large enough for /% networks.', _id, text(zonecidr), plen;
    END IF;
    incr := floor(2^(32-plen))::integer;
    nincr := floor(2^(32-masklen(zonecidr)))::integer / incr - 1;
    RAISE NOTICE 'incr = %, nincr = %', incr, nincr;
    FOR i IN 0 .. nincr LOOP
      INSERT INTO subnet (zone, range, available, usage) VALUES (_id, set_masklen(zonecidr + i*incr, plen), TRUE, 'available');
    END LOOP;
  END;
$$ LANGUAGE plpgsql;

-- primarily intended for CLI usage
CREATE OR REPLACE FUNCTION list_zones()
  RETURNS SETOF zone
  SECURITY DEFINER
AS $$
  BEGIN
    RETURN QUERY SELECT * FROM zone ORDER BY name;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION allocate_subnet(_id UUID, _usage VARCHAR)
  RETURNS SETOF subnet
  SECURITY DEFINER
AS $$
  BEGIN
    RETURN QUERY UPDATE subnet SET available = FALSE, usage = _usage
      WHERE range = ( SELECT range FROM subnet WHERE available = TRUE AND zone = _id LIMIT 1 )
      RETURNING *;
  END;
$$ LANGUAGE plpgsql;

-- not that useful, but necessary for Terraform provider implementation
CREATE OR REPLACE FUNCTION read_subnet(_id UUID)
  RETURNS SETOF subnet
  SECURITY DEFINER
AS $$
  BEGIN
    RETURN QUERY SELECT * FROM subnet WHERE id = _id;
    IF NOT FOUND THEN
      RAISE EXCEPTION 'Subnet % not found', _id;
    END IF;
    RETURN;
  END;
$$ LANGUAGE plpgsql;

-- primarily intended for CLI usage
CREATE OR REPLACE FUNCTION list_subnets()
  RETURNS SETOF subnet
  SECURITY DEFINER
AS $$
  BEGIN
    RETURN QUERY SELECT * FROM subnet ORDER BY zone, range;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE free_subnet(_id UUID)
  SECURITY DEFINER
AS $$
  BEGIN
    UPDATE subnet SET available = TRUE, usage = 'available' WHERE id = _id AND available = FALSE;
    IF NOT FOUND THEN
      RAISE EXCEPTION 'subnet % not allocated.', _id;
    END IF;
  END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE destroy_subnet(_id UUID)
  SECURITY DEFINER
AS $$
  BEGIN
    DELETE FROM subnet WHERE id = _id AND available = TRUE;
    IF NOT FOUND THEN
      RAISE EXCEPTION 'subnet % not found.', _id;
    END IF;
  END
$$ LANGUAGE plpgsql;

CREATE ROLE cloudyipam_client
  NOSUPERUSER
  NOCREATEDB
  NOINHERIT
  NOREPLICATION
  NOBYPASSRLS
  NOLOGIN;

GRANT EXECUTE ON FUNCTION create_zone TO cloudyipam_client;
GRANT EXECUTE ON PROCEDURE populate_zone TO cloudyipam_client;
GRANT EXECUTE ON FUNCTION list_zones TO cloudyipam_client;
GRANT EXECUTE ON PROCEDURE destroy_zone TO cloudyipam_client;
GRANT EXECUTE ON FUNCTION allocate_subnet TO cloudyipam_client;
GRANT EXECUTE ON PROCEDURE free_subnet TO cloudyipam_client;
GRANT EXECUTE ON FUNCTION list_subnets TO cloudyipam_client;
GRANT EXECUTE ON PROCEDURE destroy_subnet TO cloudyipam_client;
