-- Go fast!
PRAGMA synchronous=OFF;

-- Precompute area.
ALTER TABLE multipolygons ADD COLUMN area REAL;
UPDATE multipolygons SET area = Area(Transform(GEOMETRY, 25832));
CREATE INDEX idx_multipolygons_area ON multipolygons(area);
VACUUM;