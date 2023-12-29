create table if not exists users
(
    id         varchar(100) primary key,
    name       varchar(50)          null,
    email      varchar(100)         not null,
    password   varchar(250)         not null
);

create table if not exists drugs
(
    id           varchar(100) primary key,
    name         varchar(100) not null,
    approved     bool        not null,
    min_dose     int         not null,
    max_dose     int         not null,
    available_at datetime    not null
);


create table if not exists vaccinations
(
    id       varchar(100) primary key,
    name     varchar(100) not null,
    drugs_id varchar(100) not null,
    dose     int          not null,
    date     datetime     not null
);