create table if not exists users
(
    id         int auto_increment   primary key,
    name       varchar(50)          null,
    email      varchar(100)         not null,
    password   varchar(250)         not null
);