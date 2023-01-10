create table if not exists log_invoice(
  debt_id int not null,
  paid_amount decimal(15, 2) not null,
  paid_at datetime not null,
  paid_by varchar(255) not null,
  status varchar(9) default 'PENDING' not null,

  primary key (debt_id)
);
