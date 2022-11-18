CREATE TABLE balances (
                          id BIGSERIAL unique,
                          user_id BIGINT unique NOT NULL,
                          amount BIGINT NOT NULL
);

CREATE TABLE transaction_type (
                                  id SERIAL unique,
                                  type varchar NOT NULL
);

INSERT INTO transaction_type(type) VALUES ('Зачисление');
INSERT INTO transaction_type(type) VALUES ('Резервация средств');
INSERT INTO transaction_type(type) VALUES ('Плата за услугу');

CREATE TABLE transaction_history (
                                     id BIGSERIAL unique,
                                     balance_id BIGINT NOT NULL,
                                     transaction_type_id int NOT NULL,
                                     service_id BIGINT,
                                     amount BIGINT NOT NULL,
                                     created_at date,
                                     FOREIGN KEY (balance_id)
                                         REFERENCES balances(id),
                                     FOREIGN KEY (transaction_type_id)
                                         REFERENCES transaction_type(id)
);

CREATE TABLE reservation (
                             id BIGSERIAL unique,
                             user_id BIGINT NOT NULL,
                             service_id BIGINT NOT NULL,
                             order_id BIGINT unique NOT NULL,
                             amount BIGINT NOT NULL,
                             created_at date,
                             updated_at date,
                             FOREIGN KEY (user_id)
                                 REFERENCES balances(user_id)
);

CREATE TABLE statuses (
                          id SERIAL unique,
                          status varchar NOT NULL
);

INSERT INTO statuses(status) VALUES ('В процессе');
INSERT INTO statuses(status) VALUES ('Отменен');
INSERT INTO statuses(status) VALUES ('Принят');

CREATE TABLE reservation_status (
                                    id BIGSERIAL unique,
                                    reservation_id BIGINT NOT NULL,
                                    status_id int NOT NULL,
                                    FOREIGN KEY(reservation_id)
                                        REFERENCES reservation(id),
                                    FOREIGN KEY(status_id)
                                        REFERENCES statuses(id)
);