FROM ruby:2.5-alpine3.8 as gems
RUN gem install evt-message_store-postgres-database -v 1.0.0 --install-dir /gems

FROM postgres:10.6-alpine

# Run a postgres command that does nothing so initdb will trigger.
RUN docker-entrypoint.sh postgres --version

ENTRYPOINT docker-entrypoint.sh postgres
