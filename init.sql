
CREATE TABLE IF NOT EXISTS orders (
                                      order_uid TEXT PRIMARY KEY,
                                      track_number TEXT NOT NULL,
                                      entry TEXT NOT NULL,
                                      locale TEXT NOT NULL,
                                      internal_signature TEXT,
                                      customer_id TEXT NOT NULL,
                                      delivery_service TEXT NOT NULL,
                                      shardkey TEXT NOT NULL,
                                      sm_id INT NOT NULL,
                                      date_created TIMESTAMPTZ NOT NULL,
                                      oof_shard TEXT NOT NULL
);


CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);


CREATE TABLE IF NOT EXISTS deliveries (
                                          order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL
    );


CREATE TABLE IF NOT EXISTS payments (
                                        order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    provider TEXT NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
    );


CREATE TABLE IF NOT EXISTS items (
                                     id SERIAL PRIMARY KEY,
                                     order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INT NOT NULL,
    track_number TEXT NOT NULL,
    price INT NOT NULL,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    sale INT NOT NULL,
    size TEXT NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand TEXT NOT NULL,
    status INT NOT NULL
    );


CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);



INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES
    ('test-uid-1', 'TRK12345', 'WBIL', 'ru', '', 'cust-001', 'DHL', '1', 1, '2024-05-05T10:00:00Z', '1'),
    ('test-uid-2', 'TRK67890', 'WBIL', 'ru', '', 'cust-002', 'CDEK', '2', 2, '2024-06-15T14:30:00Z', '2'),
    ('test-uid-3', 'TRK55555', 'WBIL', 'ru', '', 'cust-003', 'Boxberry', '3', 3, '2024-07-01T09:15:00Z', '3')
    ON CONFLICT DO NOTHING;

INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES
                                                                                       ('test-uid-1', 'Иван Иванов', '+79991234567', '101000', 'Москва', 'ул. Пушкина, д. 1', 'Московская область', 'ivan@example.com'),
                                                                                       ('test-uid-2', 'Петр Петров', '+79998887766', '630000', 'Новосибирск', 'ул. Ленина, д. 15', 'Новосибирская область', 'petr@example.com'),
                                                                                       ('test-uid-3', 'Сергей Сергеев', '+79997776655', '190000', 'Санкт-Петербург', 'Невский проспект, д. 25', 'Ленинградская область', 'sergey@example.com')
    ON CONFLICT DO NOTHING;

INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES
                                                                                                                                                    ('test-uid-1', 'trx-001', '', 'RUB', 'wbpay', 1000, 1700000000, 'Sberbank', 200, 800, 0),
                                                                                                                                                    ('test-uid-2', 'trx-002', '', 'RUB', 'tinkoff', 2500, 1700000500, 'Tinkoff', 300, 2200, 0),
                                                                                                                                                    ('test-uid-3', 'trx-003', '', 'RUB', 'yoomoney', 500, 1700001000, 'VTB', 150, 350, 0)
    ON CONFLICT DO NOTHING;

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES
                                                                                                                          ('test-uid-1', 123456, 'TRK12345', 800, 'rid-001', 'Кроссовки', 0, '42', 800, 555555, 'Nike', 202),
                                                                                                                          ('test-uid-2', 654321, 'TRK67890', 1100, 'rid-002', 'Куртка зимняя', 10, 'M', 990, 444444, 'Columbia', 202),
                                                                                                                          ('test-uid-2', 654322, 'TRK67890', 1100, 'rid-003', 'Шапка вязаная', 0, 'L', 1100, 333333, 'Adidas', 202),
                                                                                                                          ('test-uid-3', 777777, 'TRK55555', 350, 'rid-004', 'Футболка', 0, 'XL', 350, 222222, 'Puma', 202)
    ON CONFLICT DO NOTHING;
