DROP TABLE IF EXISTS time;

create table time
(
    mac_address macaddr      not null,
    seconds      varchar(255) not null,
    router_id   integer      not null
);

create index "time__index-mac_address"
    on time (mac_address, seconds);

create table time_summary
(
    mac_address   macaddr not null,
    seconds       integer not null,
    breaks        json    not null,
    date          date    not null,
    seconds_begin integer not null,
    seconds_end   integer not null
);
