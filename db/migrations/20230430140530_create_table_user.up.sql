CREATE TABLE IF NOT EXISTS public.users
(
    id SERIAL NOT NULL,
    name character varying(255) NOT NULL,
    xid UUID NOT NULL,
    is_active bit(1) NOT NULL DEFAULT (1)::bit(1),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp with time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp with time zone,
    deleted_by character varying(255) DEFAULT NULL::character varying
);

-- Initial seeding

INSERT INTO public.users
    (name, xid, is_active)
VALUES
    ('user-testing', 'ea0212d3-abd6-406f-8c67-868e814a2436', '1');