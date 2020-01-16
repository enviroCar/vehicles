SET lc_time = "de_DE";
SET DateStyle = "German";

DROP TABLE IF EXISTS manufacturers;
DROP TABLE IF EXISTS power_sources;
DROP TABLE IF EXISTS vehicles;

CREATE TABLE vehicles (
  manufacturer_id char(4) NOT NULL,
  id char(3) NOT NULL,
  manufacturer text,
  trade_name text,
  commercial_name text,
  allotment_date date NOT NULL,
  category varchar(3) NOT NULL,
  bodywork varchar(4),
  power_source_id int NOT NULL,
  power int NOT NULL,
  engine_capacity int,
  axles int,
  powered_axles int,
  seats int,
  maximum_mass int,
  CONSTRAINT vehicles_pkey PRIMARY KEY (manufacturer_id, id)
);

CREATE TABLE manufacturers (
  id char(4) PRIMARY KEY,
  name text
);

CREATE TABLE power_sources (
  id int PRIMARY KEY, 
  short_name text, 
  description text
);


COPY vehicles FROM '/data/vehicles.csv'  WITH (FORMAT csv, DELIMITER ',', QUOTE '"', HEADER);
COPY power_sources FROM '/data/power_sources.csv'  WITH (FORMAT csv, DELIMITER ',', QUOTE '"', HEADER);


INSERT INTO manufacturers(id, name)
  SELECT DISTINCT
    manufacturer_id AS id,
    manufacturer AS name
  FROM (
    SELECT 
      manufacturer_id, 
      max(allotment_date) as allotment_date
    FROM vehicles
    GROUP BY manufacturer_id
  ) AS dates
  JOIN vehicles AS v USING (manufacturer_id)
  WHERE v.allotment_date = dates.allotment_date;


ALTER TABLE vehicles DROP COLUMN manufacturer;
ALTER TABLE vehicles ADD FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id);
ALTER TABLE vehicles ADD FOREIGN KEY (power_source_id) REFERENCES power_sources(id);