CREATE TABLE IF NOT EXISTS public.wallet_transactions
(
    id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    status character varying(50) NOT NULL,
    reference_id UUID NOT NULL UNIQUE,
    amount NUMERIC(15,3) NOT NULL,
    description TEXT NOT NULL,
    is_active bit(1) NOT NULL DEFAULT '1',
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp with time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp with time zone,
    deleted_by character varying(255) DEFAULT NULL::character varying
);