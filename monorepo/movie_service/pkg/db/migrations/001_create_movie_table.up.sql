create table movie(
    id UUID NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    director VARCHAR NOT NULL,
    year NUMERIC(4) NOT NULL,
    plot VARCHAR NOT NULL,
    state NUMERIC(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

COMMENT ON COLUMN movie.state IS 'state=1 bu faol, state=0 o''chirilgan, state=-1 arxiv';