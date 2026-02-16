CREATE TABLE public.users (
    id              UUID PRIMARY KEY DEFAULT
                    gen_random_uuid(),
    username        TEXT NOT NULL,
    email           CHARACTER(255) NOT NULL,
    password        CHARACTER(255) NOT NULL,
    profile_image   TEXT NOT NULL,
    role            CHARACTER(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL
); 