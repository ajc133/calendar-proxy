ARG RUST_VERSION=1.72.0
ARG APP_NAME=calendar-proxy
FROM docker.io/rust:${RUST_VERSION}-slim-bookworm AS build
ARG APP_NAME
WORKDIR /app

RUN apt-get update && apt-get -y install pkg-config libssl-dev && rm -rf /var/lib/apt/lists/*

# From https://docs.docker.com/language/rust/develop/
# Build the application.
# Leverage a cache mount to /usr/local/cargo/registry/
# for downloaded dependencies and a cache mount to /app/target/ for 
# compiled dependencies which will speed up subsequent builds.
# Leverage a bind mount to the src directory to avoid having to copy the
# source code into the container. Once built, copy the executable to an
# output directory before the cache mounted /app/target is unmounted.
RUN --mount=type=bind,source=src,target=src \
    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
    --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/usr/local/cargo/registry/ \
cargo build --locked --release && cp ./target/release/$APP_NAME /bin/server


FROM docker.io/debian:bookworm-slim AS final
WORKDIR /app

RUN apt-get update && apt-get -y install pkg-config libssl-dev ca-certificates&& rm -rf /var/lib/apt/lists/*
# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
RUN mkdir -p /app/data && chown appuser:appuser /app/data
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

# Expose the port that the application listens on.
EXPOSE 3000

ENTRYPOINT ["/bin/server"]
