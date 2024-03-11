CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name varchar(255) ,
                                     email VARCHAR(255) ,
                                     password VARCHAR(255),
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) ,
                                        price NUMERIC(10, 2),
                                        description TEXT ,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      customerName varchar(255) not null,
                                      totalAmount int,
                                      user_id INTEGER REFERENCES users(id),
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      product_id serial references products(id)
);

INSERT INTO users (email, password) VALUES
                                        ('user1@example.com', 'password1'),
                                        ('user2@example.com', 'password2'),
                                        ('user3@example.com', 'password3');

INSERT INTO products (name, price, description) VALUES
                                                    ('Product 1', 10.99, 'Description for product 1'),
                                                    ('Product 2', 19.99, 'Description for product 2'),
                                                    ('Product 3', 5.99, 'Description for product 3');

INSERT INTO orders (customerName, totalAmount, user_id, product_id) VALUES
                                                                        ('John Doe', 30, 1, 1),
                                                                        ('Jane Smith', 20, 2, 2),
                                                                        ('Alice Johnson', 10, 3, 3);
