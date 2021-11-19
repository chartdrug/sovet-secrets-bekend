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

       --fold: {wrist: 2, forearm: 0, shoulder_front: 0, chest: 11, xiphoid: 0, belly: 12, anterrior_iliac: 0,…}
       fold_anterrior_iliac decimal,
       fold_back decimal,
       fold_belly decimal, -- 12
       fold_chest decimal, -- 11
       fold_forearm decimal,
       fold_hip_front decimal, -- 13
       fold_hip_inside decimal,
       fold_hip_rear decimal,
       fold_hip_side decimal,
       fold_scapula decimal,
       fold_shin decimal,
       fold_shoulder_front decimal,
       fold_shoulder_rear decimal,
       fold_waist_side decimal,
       fold_wrist decimal,
       fold_xiphoid decimal,

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
CREATE INDEX idx_injection_owner ON injection(owner);

CREATE TABLE injection_dose
(
    id uuid NOT NULL,
    id_injection uuid NOT NULL,
    dose double precision,
    drug varchar(37),
    volume double precision,
    solvent character(1) COLLATE pg_catalog."default",
    points integer,
    CONSTRAINT injection_dose_pkey PRIMARY KEY (id)
);

CREATE TABLE concentration (
       owner uuid NOT NULL,
       Id_injection uuid NOT NULL,
       drug varchar(37),
       dt BIGINT,
       C double precision,
       CC double precision,
       CCT double precision,
       CT double precision
);
CREATE INDEX idx_concentration_owner ON concentration(owner);