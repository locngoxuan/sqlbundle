--+up BEGIN
CREATE TABLE IF NOT EXISTS third (
  ID serial NOT NULL,
  NAME varchar(100),
  PRIMARY KEY(ID)
);
--+up END

--+down BEGIN
DROP TABLE IF EXISTS third;
--+down END