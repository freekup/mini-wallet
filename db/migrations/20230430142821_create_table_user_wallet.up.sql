CREATE TABLE IF NOT EXISTS public.user_wallets
(
    id UUID NOT NULL,
    user_xid UUID NOT NULL UNIQUE,
    current_balance NUMERIC(10,3) NOT NULL,
    is_enabled bit(1) NOT NULL DEFAULT '0',
    enabled_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp with time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp with time zone,
    deleted_by character varying(255) DEFAULT NULL::character varying
);