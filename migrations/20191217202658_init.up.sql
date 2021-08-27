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

CREATE TABLE injection
(
    id uuid NOT NULL,
    owner uuid NOT NULL,
    dt timestamp without time zone NOT NULL DEFAULT ('now'::text)::date,
    course uuid,
    what character(1) COLLATE pg_catalog."default" NOT NULL DEFAULT '?'::bpchar,
    /*dose double precision[],
    drug uuid[],
    volume double precision[],
    solvent character(1)[] COLLATE pg_catalog."default",
    points integer[],
    zerodt timestamp without time zone,
    hashid uuid,
    cutoff integer[],*/
    CONSTRAINT injection_pkey PRIMARY KEY (id)
);

CREATE TABLE injection_dose
(
    id uuid NOT NULL,
    id_injection uuid NOT NULL,
    dose double precision,
    drug uuid,
    volume double precision,
    solvent character(1) COLLATE pg_catalog."default",
    points integer,
    CONSTRAINT injection_dose_pkey PRIMARY KEY (id)
);



CREATE INDEX idx_injection_owner ON injection(owner);