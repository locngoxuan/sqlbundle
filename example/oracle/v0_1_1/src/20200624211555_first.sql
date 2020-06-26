--+up BEGIN
CREATE SEQUENCE oracle_first_seq START WITH 1 increment by 1;

CREATE TABLE oracle_first (
  ID NUMBER(19, 0) DEFAULT oracle_first_seq.nextval NOT NULL,
  NAME varchar2(100),
  PRIMARY KEY(ID)
);
--+up END

--+down BEGIN
DROP SEQUENCE oracle_first_seq;
DROP TABLE oracle_first;
--+down END