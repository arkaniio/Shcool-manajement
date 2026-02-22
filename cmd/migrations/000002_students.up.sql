CREATE TABLE public.students (
    id          UUID PRIMARY KEY DEFAULT
                gen_random_uuid(),
    name        VARCHAR(50) NOT NULL,
    class       VARCHAR(50) NOT NULL,
    address     TEXT NOT NULL,
    major       VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);