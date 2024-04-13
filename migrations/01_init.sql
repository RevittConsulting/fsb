create extension if not exists "uuid-ossp";

drop schema if exists fsb cascade;
create schema fsb;

create table if not exists fsb.invoices (
    id uuid primary key default uuid_generate_v4(),
    party_id uuid not null,
    amount decimal(10,2) not null,
    meta jsonb default null,
    created_at timestamp without time zone not null default (now() at time zone 'utc'),
    deleted_at timestamp without time zone default null
);

create table if not exists fsb.invoice_payments (
    id uuid primary key default uuid_generate_v4(),
    invoice_id uuid not null references fsb.invoices(id),
    meta jsonb default null,
    success boolean default null,
    created_at timestamp without time zone not null default (now() at time zone 'utc'),
    deleted_at timestamp without time zone default null
);

---- create above / drop below ----

drop table if exists fsb.invoice_payments;
drop table if exists fsb.invoices;