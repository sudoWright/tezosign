create table vestings
(
	vst_id serial not null
		constraint vesting_pk
			primary key,
    vst_name varchar not null,
    vst_address varchar(36) not null,
	ctr_id int
		constraint vestings_contracts_ctr_id_fk
			references contracts
);

create unique index vestings_ctr_id_vst_address_uindex
	on vestings (ctr_id, vst_address);

