CREATE EXTENSION timescaledb;
SELECT create_hypertable('conditions', 'time');
