CREATE TABLE IF NOT EXISTS widgets (
    id uuid PRIMARY KEY,
    owner_id text NOT NULL,
    organisation_id text NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS widgets_org_created_idx
    ON widgets (organisation_id, created_at DESC, id);

CREATE TABLE IF NOT EXISTS example_outbox (
    id text PRIMARY KEY,
    topic text NOT NULL,
    payload jsonb NOT NULL,
    attempts integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at timestamptz,
    claimed_by text,
    claimed_until timestamptz,
    claim_token text
);

CREATE INDEX IF NOT EXISTS example_outbox_pending_idx
    ON example_outbox (created_at)
    WHERE sent_at IS NULL;

