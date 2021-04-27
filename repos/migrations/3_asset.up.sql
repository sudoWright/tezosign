create table assets
(
	ast_id serial not null
		constraint asset_pk
			primary key,
    ast_name varchar not null,
    ast_contract_type varchar not null,
    ast_address varchar(36) not null,
    ast_dexter_address varchar(36),
    ast_scale int not null,
    ast_ticker varchar not null,
    ast_token_id int,
    ast_is_active bool default TRUE not null,
    ast_last_block_level int,
    ast_updated_at timestamp default now() not null,
	ctr_id int
		constraint assets_contracts_ctr_id_fk
			references contracts
);

create unique index assets_ctr_id_ast_address_ast_token_id_uindex
	on assets (ctr_id, ast_address,ast_token_id);

