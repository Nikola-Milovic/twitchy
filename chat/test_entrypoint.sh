#!/bin/sh
# Adapted from Alex Kleissner's post, Running a Phoenix 1.3 project with docker-compose
# https://medium.com/@hex337/running-a-phoenix-1-3-project-with-docker-compose-d82ab55e43cf

set -e

# Ensure the app's dependencies are installed
mix deps.get

# Prepare Dialyzer if the project has Dialyxer set up
if mix help dialyzer >/dev/null 2>&1
then
  echo "\nFound Dialyxer: Setting up PLT..."
  mix do deps.compile, dialyzer --plt
else
  echo "\nNo Dialyxer config: Skipping setup..."
fi

#Analysis style code
# Prepare Credo if the project has Credo start code analyze
if mix help credo >/dev/null 2>&1
then
  echo "\nFound Credo: analyzing..."
  mix credo || true
else
  echo "\nNo Credo config: Skipping code analyze..."
fi

# Wait for Postgres to become available.
until pg_isready -h $POSTGRES_HOST -U "postgres"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "\nPostgres is available: continuing with database setup..."

# Set up the database
MIX_ENV=test mix ecto.reset

echo "\n Testing..."
# Running the tests
mix test