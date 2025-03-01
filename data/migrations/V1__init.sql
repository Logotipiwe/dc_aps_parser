CREATE TABLE IF NOT EXISTS parsers
(
    id       bigint            NOT NULL,
    link     character varying NOT NULL,
    username character varying NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions
(
    chat_id           bigint  NOT NULL,
    permitted_aps_num integer NOT NULL
);

CREATE TABLE IF NOT EXISTS results_history
(
    id uuid NOT NULL default gen_random_uuid(),
    ap_id bigint NOT NULL unique,
    created_at timestamp NOT NULL default current_timestamp,
    data jsonb
);