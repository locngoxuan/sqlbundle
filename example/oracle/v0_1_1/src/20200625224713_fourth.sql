--+up BEGIN
CREATE TABLE fourth (
  ID NUMBER(19) NOT NULL,
  NAME varchar2(100),
  PRIMARY KEY(ID)
)
--+up END

--+down BEGIN
DROP TABLE fourth
--+down END