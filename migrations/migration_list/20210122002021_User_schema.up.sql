
CREATE TABLE IF NOT EXISTS public.users
(
    uuid    TEXT    NOT NULL    UNIQUE  PRIMARY KEY,
    name    TEXT    NOT NULL,
    age     INT     NOT NULL,

    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL    DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE 
)
WITH (
    OIDS = FALSE
);
