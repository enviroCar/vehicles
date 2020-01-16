FROM postgres:alpine


ENV POSTGRES_DB vehicles

COPY vehicles.csv power_sources.csv /data/
COPY schema.sql /docker-entrypoint-initdb.d/