CREATE TABLE IF NOT EXISTS People (
  id SERIAL PRIMARY KEY,
  passport_serie char(4),
  passport_number char(6),
  name char(255),
  surname char(255),
  patronymic char(255),
  address char(255)
)