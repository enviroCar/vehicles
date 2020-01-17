FROM postgres:alpine


ENV POSTGRES_DB vehicles

COPY vehicles.csv power_sources.csv /data/
COPY schema.sql /docker-entrypoint-initdb.d/

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD \
  pg_isready -q -h localhost -d ${POSTGRES_DB} -U ${POSTGRES_USER}
