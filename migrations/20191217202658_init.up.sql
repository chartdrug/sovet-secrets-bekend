CREATE TABLE album
(
    id         VARCHAR PRIMARY KEY,
    name       VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE users (
   id uuid NOT NULL,
   login character varying(64) NOT NULL,
   passwd character varying(128) NOT NULL,
   email character varying(64) NOT NULL,
   date_registered timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
   date_lastlogin timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,


   sex character varying(1) NOT NULL,
   birthday DATE NOT NULL,

   height decimal not null,
   weight decimal not null
);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TABLE antro (
       id uuid NOT NULL,
       owner uuid NOT NULL,
       dt date DEFAULT ('now'::text)::date NOT NULL,
       -- general jsonb,
       general_age int,
       general_hip decimal,
       general_height decimal,
       general_leglen decimal,
       general_weight decimal,
       general_handlen decimal,
       general_shoulders decimal,

       --girth jsonb,
       --fold jsonb,
       notes text,
       basic boolean DEFAULT false NOT NULL,
       --result jsonb
       result_fat decimal,
       result_nofat decimal,
       result_energy decimal
);
ALTER TABLE ONLY antro
    ADD CONSTRAINT antro_pkey PRIMARY KEY (id);

CREATE TABLE injection (
       id UUID PRIMARY KEY,
       owner UUID NOT NULL,
       dt DATETIME NOT NULL DEFAULT CURRENT_DATE,
       course UUID DEFAULT NULL,
       what CHAR(1) NOT NULL DEFAULT '?',
       dose DOUBLE PRECISION[],
       drug UUID[],
       volume DOUBLE PRECISION[],
       solvent CHAR(1)[],
       points INTEGER[],
       zerodt DATETIME[],
       hashid UUID DEFAULT NULL
);
CREATE INDEX idx_injection_owner ON injection(owner);