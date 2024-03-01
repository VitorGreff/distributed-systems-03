DROP DATABASE IF EXISTS sd2;
CREATE DATABASE sd2;
\c sd2;

DROP TABLE IF EXISTS products;
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    quantity INT NOT NULL
);

INSERT INTO products (name, price, quantity) 
VALUES
    ('Tênis de Caminhada', 150.00, 6),
    ('Cafeteira Elétrica', 250.00, 5),
    ('Monitor de 24 Polegadas', 300.00, 7),
    ('Laptop Gamer', 1200.00, 3),
    ('Smartphone de Ultima Geração', 800.00, 1),
    ('Cadeira de Escritório', 400.00, 8),
    ('Bicicleta de Montanha', 1000.00, 2),
    ('Fone de Ouvido sem Fio', 150.00, 1),
    ('Câmera Fotográfica DSLR', 1000.00, 4),
    ('Cerveja Artesanal', 20.00, 2);
