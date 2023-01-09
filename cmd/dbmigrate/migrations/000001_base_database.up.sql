create table if not exists invoices(
  debt_id int not null,
  debt_amount decimal(15, 2) not null,
  debt_due_date date not null,
  email varchar(255) not null,
  government_id varchar(20) not null,
  name varchar(255) not null,

  primary key (debt_id)
)
