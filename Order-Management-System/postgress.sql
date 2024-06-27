
CREATE TABLE product(
   id VARCHAR (255) PRIMARY KEY,
   name VARCHAR (255) UNIQUE NOT NULL,
   description VARCHAR (255) NOT NULL,
   avail_qty INT NOT NULL,
   reserve_qty INT NOT NULL
);


insert into product(id, name, description, avail_qty, reserve_qty) values('p1', 'cycle1', 'hero cycle', 10, 0);
insert into product(id, name, description, avail_qty, reserve_qty) values('p2', 'bike2', 'hero bike', 10, 0);
insert into product(id, name, description, avail_qty, reserve_qty) values('p3', 'computer3', 'amd processor computer', 10, 0);


COMMIT;

SELECT * FROM product;