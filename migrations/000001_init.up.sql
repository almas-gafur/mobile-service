CREATE TABLE IF NOT EXISTS workshops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS masters (
    id SERIAL PRIMARY KEY,
    workshop_id INT NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    imei VARCHAR(15) UNIQUE CHECK (imei IS NULL OR imei ~ '^[0-9]{15}$'),
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS repair_tickets (
    id SERIAL PRIMARY KEY,
    short_hash VARCHAR(8) UNIQUE CHECK (short_hash IS NULL OR short_hash ~ '^[A-Za-z0-9]{8}$'),
    workshop_id INT NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    device_id INT NOT NULL REFERENCES devices(id) ON DELETE RESTRICT,
    client_name VARCHAR(100) NOT NULL,
    client_phone VARCHAR(20) NOT NULL CHECK (char_length(client_phone) BETWEEN 7 AND 20),
    status VARCHAR(20) NOT NULL DEFAULT 'Заявка' CHECK (status IN ('Заявка', 'Принято', 'В работе', 'Готово', 'Выдано')),
    defect_description TEXT NOT NULL DEFAULT '',
    water_damage BOOLEAN NOT NULL DEFAULT FALSE,
    warranty_days INT NOT NULL DEFAULT 0 CHECK (warranty_days >= 0 AND warranty_days <= 730),
    price INT NOT NULL DEFAULT 0 CHECK (price >= 0),
    rating INT CHECK (rating BETWEEN 1 AND 5),
    review_text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_hash ON repair_tickets(short_hash);
CREATE INDEX IF NOT EXISTS idx_repair_tickets_workshop_status ON repair_tickets(workshop_id, status);
CREATE INDEX IF NOT EXISTS idx_repair_tickets_created_at ON repair_tickets(created_at DESC);
