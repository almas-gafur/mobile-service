ALTER TABLE devices ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE devices ALTER COLUMN imei DROP NOT NULL;
ALTER TABLE devices ALTER COLUMN brand SET NOT NULL;
ALTER TABLE devices ALTER COLUMN model SET NOT NULL;
ALTER TABLE devices ALTER COLUMN imei TYPE VARCHAR(15);
ALTER TABLE devices ALTER COLUMN brand TYPE VARCHAR(50);
ALTER TABLE devices ALTER COLUMN model TYPE VARCHAR(50);
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'devices' AND column_name = 'client_phone'
    ) THEN
        EXECUTE 'ALTER TABLE devices ALTER COLUMN client_phone DROP NOT NULL';
    END IF;
END $$;

ALTER TABLE repair_tickets ADD COLUMN IF NOT EXISTS client_name VARCHAR(100) NOT NULL DEFAULT 'Клиент';
ALTER TABLE repair_tickets ADD COLUMN IF NOT EXISTS client_phone VARCHAR(20);
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'devices' AND column_name = 'client_phone'
    ) THEN
        EXECUTE 'UPDATE repair_tickets t
                 SET client_phone = COALESCE(t.client_phone, d.client_phone, ''+70000000000'')
                 FROM devices d
                 WHERE d.id = t.device_id';
    END IF;
END $$;
UPDATE repair_tickets SET client_phone = '+70000000000' WHERE client_phone IS NULL;
ALTER TABLE repair_tickets ALTER COLUMN client_phone SET NOT NULL;

ALTER TABLE repair_tickets ADD COLUMN IF NOT EXISTS defect_description TEXT NOT NULL DEFAULT '';
UPDATE repair_tickets
SET defect_description = COALESCE(defect_description, '')
WHERE defect_description IS NULL;

ALTER TABLE repair_tickets ADD COLUMN IF NOT EXISTS rating INT;
ALTER TABLE repair_tickets ADD COLUMN IF NOT EXISTS review_text TEXT;

ALTER TABLE repair_tickets DROP CONSTRAINT IF EXISTS repair_tickets_status_check;
UPDATE repair_tickets
SET status = CASE status
    WHEN 'accepted' THEN 'Принято'
    WHEN 'in_progress' THEN 'В работе'
    WHEN 'done' THEN 'Готово'
    WHEN 'issued' THEN 'Выдано'
    ELSE status
END;

ALTER TABLE repair_tickets ALTER COLUMN short_hash DROP NOT NULL;
ALTER TABLE repair_tickets ALTER COLUMN short_hash TYPE VARCHAR(8);
UPDATE repair_tickets
SET short_hash = substr(md5(random()::text || id::text), 1, 8)
WHERE short_hash IS NOT NULL AND short_hash !~ '^[A-Za-z0-9]{8}$';

ALTER TABLE repair_tickets ALTER COLUMN price TYPE INT USING round(price)::INT;
ALTER TABLE repair_tickets ALTER COLUMN price SET DEFAULT 0;
ALTER TABLE repair_tickets ALTER COLUMN warranty_days SET DEFAULT 0;
ALTER TABLE repair_tickets ALTER COLUMN water_damage SET DEFAULT FALSE;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'repair_tickets_status_check'
    ) THEN
        ALTER TABLE repair_tickets
            ADD CONSTRAINT repair_tickets_status_check
            CHECK (status IN ('Заявка', 'Принято', 'В работе', 'Готово', 'Выдано'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'repair_tickets_short_hash_check'
    ) THEN
        ALTER TABLE repair_tickets
            ADD CONSTRAINT repair_tickets_short_hash_check
            CHECK (short_hash IS NULL OR short_hash ~ '^[A-Za-z0-9]{8}$');
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'repair_tickets_rating_check'
    ) THEN
        ALTER TABLE repair_tickets
            ADD CONSTRAINT repair_tickets_rating_check
            CHECK (rating BETWEEN 1 AND 5);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_tickets_hash ON repair_tickets(short_hash);
CREATE INDEX IF NOT EXISTS idx_repair_tickets_workshop_status ON repair_tickets(workshop_id, status);
CREATE INDEX IF NOT EXISTS idx_repair_tickets_created_at ON repair_tickets(created_at DESC);
