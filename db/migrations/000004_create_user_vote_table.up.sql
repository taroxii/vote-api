BEGIN;

create table if not EXISTS user_items (
    id serial PRIMARY KEY,
    correlation_id VARCHAR(255),
    user_id integer not null constraint user_items_user_id_fk references users on update cascade on delete cascade,
    item_id integer not null constraint user_items_item_id_fk references items on update cascade on delete cascade,
    created_at timestamp,
    updated_at timestamp
);
create unique index if not exists user_item_user_id_user_items_item_id_uindex on user_items (user_id, item_id);
COMMIT;

