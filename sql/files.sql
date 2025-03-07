CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(128) NOT NULL,
    path varchar(250) NOT NULL,
    created_at timestamp default current_timestamp 
);