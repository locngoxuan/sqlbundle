--+up BEGIN
CREATE TABLE third (
  ID NUMBER(19) NOT NULL,
  NAME varchar2(100),
  PRIMARY KEY(ID)
)
--+up END

--+down BEGIN
DROP TABLE third
--+down END