CREATE EXTENSION IF NOT EXISTS dblink;

DO
$$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE  rolname = '{{ ghost_postgres_user }}'
    ) THEN
        RAISE NOTICE 'User "{{ ghost_postgres_user }}" already exists. Skipping.';
    ELSE
        CREATE USER {{ ghost_postgres_user }} CREATEDB LOGIN ENCRYPTED PASSWORD '{{ ghost_postgres_pass }}';
    END IF;
END
$$;


DO
$$
BEGIN
    IF EXISTS (
      SELECT FROM pg_catalog.pg_database
      WHERE datname = '{{ ghost_postgres_db }}'
    ) THEN
        RAISE NOTICE 'Database "{{ ghost_postgres_db }}" already exists. Skipping.';
    ELSE
        PERFORM dblink_exec(
            'dbname=' || current_database(),
            'CREATE DATABASE {{ ghost_postgres_db }} OWNER {{ ghost_postgres_user }}'
        );
    END IF;
END
$$;
