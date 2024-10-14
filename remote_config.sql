create table config
(
    id      integer            not null
        primary key,
    name    TEXT               not null,
    value   TEXT               not null,
    message TEXT               not null,
    secret  TEXT    default '' not null,
    enable  integer default 0  not null
);

create index key_name
    on config (name);

create table config_history
(
    id          INTEGER not null
        primary key autoincrement,
    config_id   INTEGER not null,
    old_value   TEXT    not null,
    new_value   TEXT    not null,
    enable      integer not null,
    create_time integer not null
);

create index key_config_id
    on config_history (config_id);

