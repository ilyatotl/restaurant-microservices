CREATE TABLE dishes
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    price       INT          NOT NULL,
    quantity    INT          NOT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP
);

CREATE TABLE orders
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    INT         NOT NULL,
    status     VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE order_dish
(
    id       BIGSERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    dish_id  INT NOT NULL,
    quantity INT NOT NULL,
    price    INT NOT NULL
);