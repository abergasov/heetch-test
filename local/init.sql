create table zombies
(
    id         uuid
        constraint zombies_pk
            primary key,
    updated_at timestamp,
    point      point,
    status     varchar
);