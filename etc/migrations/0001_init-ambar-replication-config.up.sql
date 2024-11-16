CREATE USER replication REPLICATION LOGIN PASSWORD 'repl-pass';

GRANT CONNECT ON DATABASE "eventstore" TO replication;
GRANT USAGE ON SCHEMA public TO replication;
GRANT SELECT ON TABLE event TO replication;

CREATE PUBLICATION event_publication FOR TABLE public.event
