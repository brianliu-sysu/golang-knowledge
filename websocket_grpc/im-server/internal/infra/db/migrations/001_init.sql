create table if not exists user (
    id bigint primary key auto_increment,
    uuid varchar(255) not null,
    email varchar(255) not null,
    phone varchar(255) not null,
    name varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);
create index if not exists idx_user_uuid on user (uuid);
create index if not exists idx_user_email on user (email);
create index if not exists idx_user_phone on user (phone);

create table if not exists friend (
    id bigint primary key auto_increment,
    user_id bigint not null,
    friend_id bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);
create index if not exists idx_friend_user_id on friend (user_id);
create index if not exists idx_friend_friend_id on friend (friend_id);

create table if not exists group (
    id bigint primary key auto_increment,
    uuid varchar(255) not null,
    name varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);
create index if not exists idx_group_uuid on group (uuid);

create table if not exists group_member (
    id bigint primary key auto_increment,
    group_id bigint not null,
    user_id bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

create index if not exists idx_group_member_group_id on group_member (group_id);
create index if not exists idx_group_member_user_id on group_member (user_id);

-- need to add a field to indicate the message type, 1 for single message, 2 for group message
create table if not exists message (
    id bigint primary key auto_increment,
    from_user_id bigint not null,
    to_user_id bigint not null,
    to_group_id bigint not null,
    message_type int not null,
    content text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

-- inbox table to store the messages that the user has received
create table if not exists inbox (
    id bigint primary key auto_increment,
    user_id bigint not null,
    message_id bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

create index if not exists idx_inbox_user_id on inbox (user_id);