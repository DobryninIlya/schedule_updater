CREATE TABLE IF NOT EXISTS notify_clients
(
    id integer NOT NULL,
    vk_id bigint,
    tg_id bigint,
    schedule_change boolean,
    CONSTRAINT notify_clients_pkey PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS notify_clients
    OWNER to botkai;