begin;

create type order_status as enum ('reserved', 'paid', 'cancelled', 'failed', 'refunded');
create type reservation_status as enum('active', 'expired', 'cancelled', 'confirmed');
create type payment_status as enum('pending', 'completed', 'failed', 'refunded');

create table users (
    id uuid primary key,
    full_name varchar(128) not null
);

create table products (
    id uuid primary key, 
    name varchar(128) unique not null,
    quantity integer not null,
    price integer not null
);

create table orders (
    id uuid primary key,
    user_id uuid not null references users(id),
    status order_status not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table order_items (
    id uuid primary key,
    order_id uuid not null references orders(id),
    product_id uuid not null references products(id),
    quantity integer not null,
    constraint order_items_order_product_uq unique (order_id, product_id)
);

create table order_reservations (
    id uuid primary key,
    order_id uuid not null unique references orders(id),
    status reservation_status not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table item_reservations (
    id uuid primary key, 
    item_id uuid unique not null references order_items(id),
    order_reservation_id uuid unique not null references order_reservations(id)
);

create table payments (
    payment_id uuid primary key,
    user_id uuid not null references users(id),
    order_id uuid not null references orders(id),
    amount integer not null,
    status payment_status,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

commit;