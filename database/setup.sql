drop table if exists posts;

create table posts (
  id      serial primary key,
  content text,
  author  varchar(255)
);

drop table if exists pet;

create table pet (
  id      serial primary key,
  name    text,
  owner   varchar(255)
);