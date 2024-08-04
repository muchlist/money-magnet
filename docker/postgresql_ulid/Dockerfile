ARG postgresql_major=16
ARG postgresql_release=${postgresql_major}.2

ARG pgx_ulid_release=0.1.5

####################
# Postgres
####################
FROM postgres:${postgresql_release} as base

# Redeclare args for use in subsequent stages
ARG TARGETARCH
ARG postgresql_major

####################
# Extension: pgx_ulid
####################
FROM base as pgx_ulid

# Download package archive
ARG pgx_ulid_release
ADD "https://github.com/pksunkara/pgx_ulid/releases/download/v${pgx_ulid_release}/pgx_ulid-v${pgx_ulid_release}-pg${postgresql_major}-${TARGETARCH}-linux-gnu.deb" \
    /tmp/pgx_ulid.deb

####################
# Collect extension packages
####################
FROM scratch as extensions
COPY --from=pgx_ulid /tmp/*.deb /tmp/

####################
# Build final image
####################
FROM base as production

# Setup extensions
COPY --from=extensions /tmp /tmp

RUN apt-get update && apt-get install -y --no-install-recommends \
    /tmp/*.deb \
    && rm -rf /var/lib/apt/lists/* /tmp/*