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
   birthday DATE NOT NULL/*,

   height decimal not null,
   weight decimal not null*/
  ,type_sports text[]
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

CREATE TABLE IF NOT EXISTS spr (
     name_spr TEXT,
     id INT,
     id_parent INT,
     name_ru TEXT,
     name_eng TEXT,
     describe TEXT
);
INSERT INTO spr VALUES
('sport',12,0,'Скоростно-силовые','',''),
('sport',190,12,'Армрестлинг','',''),
('sport',193,12,'Бодибилдинг','',''),
('sport',191,12,'Пауэрлифтинг','',''),
('sport',192,12,'Тяжелая атлетика','',''),
('sport',194,12,'Фитнес','',''),
('sport',35,12,'Бег на короткие дистанции (спринт)','',''),
('sport',71,12,'Велотрековые гонки (спринт)','',''),
('sport',72,12,'Метание диска','',''),
('sport',78,12,'Плавание (короткие дистанции)','',''),
('sport',74,12,'Прыжки в высоту','',''),
('sport',75,12,'Прыжок в длину','',''),
('sport',73,12,'Толкание ядра','',''),
('sport',76,12,'Тройной прыжок','',''),
('sport',13,0,'Многоборья','',''),
('sport',212,13,'Авиационное многоборье','',''),
('sport',204,13,'Акватлон','',''),
('sport',210,13,'Военное троеборье','',''),
('sport',203,13,'Дуатлон','',''),
('sport',205,13,'Зимний триатлон','',''),
('sport',207,13,'Квадратлон','',''),
('sport',206,13,'Кросс-кантри триатлон','',''),
('sport',209,13,'Многоборье спасателей МЧС','',''),
('sport',213,13,'Морское пятиборье','',''),
('sport',211,13,'Офицерское многоборье','',''),
('sport',208,13,'Пятиборье современное','',''),
('sport',202,13,'Триатлон','',''),
('sport',14,0,'Единоборства','',''),
('sport',81,14,'Айкидо','',''),
('sport',82,14,'Бокс','',''),
('sport',83,14,'Борьба вольная','',''),
('sport',84,14,'Борьба классическая','',''),
('sport',188,14,'Джиу-джитсу','',''),
('sport',85,14,'Дзюдо','',''),
('sport',91,14,'Капоэйра','',''),
('sport',189,14,'ММА','',''),
('sport',86,14,'Самбо','',''),
('sport',87,14,'Тайский бокс','',''),
('sport',88,14,'Тхэквондо','',''),
('sport',90,14,'Ушу','',''),
('sport',89,14,'Фехтование','',''),
('sport',24,0,'Спортивные игры','',''),
('sport',41,24,'Американский футбол','',''),
('sport',42,24,'Бадминтон','',''),
('sport',43,24,'Баскетбол','',''),
('sport',44,24,'Мини-баскетбол','',''),
('sport',45,24,'Бейсбол','',''),
('sport',47,24,'Бильярд','',''),
('sport',49,24,'Велобол','',''),
('sport',50,24,'Велополо','',''),
('sport',51,24,'Водное поло','',''),
('sport',52,24,'Волейбол','',''),
('sport',54,24,'Гандбол большой','',''),
('sport',53,24,'Гандбол зальный (малый)','',''),
('sport',55,24,'Го (японские шашки)','',''),
('sport',56,24,'Гольф','',''),
('sport',57,24,'Мини-гольф','',''),
('sport',58,24,'Конное поло','',''),
('sport',59,24,'Крикет','',''),
('sport',60,24,'Крокет','',''),
('sport',61,24,'Регби','',''),
('sport',63,24,'Теннис настольный','',''),
('sport',62,24,'Теннис','',''),
('sport',64,24,'Футбол','',''),
('sport',65,24,'Хоккей на роликах','',''),
('sport',66,24,'Хоккей на траве','',''),
('sport',67,24,'Хоккей с шайбой','',''),
('sport',68,24,'Шахматы','',''),
('sport',69,24,'Шашки','',''),
('sport',25,0,'Сложнокоординационные','',''),
('sport',195,25,'Акробатика','',''),
('sport',196,25,'Водный слалом','',''),
('sport',197,25,'Гимнастика спортивная','',''),
('sport',198,25,'Гимнастика художественная','',''),
('sport',199,25,'Горнолыжные дисциплины','',''),
('sport',200,25,'Прыжки в воду','',''),
('sport',201,25,'Синхронное плавание','',''),
('sport',26,0,'Циклические','',''),
('sport',28,26,'Биатлон','',''),
('sport',29,26,'Велоспорт (длинные дистанции)','',''),
('sport',30,26,'Гребля академическая','',''),
('sport',31,26,'Гребля на байдарках и каноэ','',''),
('sport',32,26,'Конькобежный спорт','',''),
('sport',36,26,'Бег (длинные дистанции)','',''),
('sport',70,26,'Бег (средние дистанции)','',''),
('sport',77,26,'Лыжные гонки','',''),
('sport',38,26,'Марафонский бег','',''),
('sport',80,26,'Плавание (длинные дистанции)','',''),
('sport',79,26,'Плавание (средние дистанции)','',''),
('sport',40,26,'Спортивная ходьба','','');

update spr set name_eng = name_ru;

CREATE EXTENSION IF NOT EXISTS tablefunc;

create table history_login(
      id_user uuid NOT NULL,
      date_event timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
      ip_address TEXT,
     Country TEXT,
     Region TEXT,
     City TEXT
);
CREATE INDEX idx_history_login_user ON history_login(id_user);

alter table injection ADD COLUMN calc boolean default false;

alter table users ADD COLUMN admin boolean default false;
alter table users ADD COLUMN blocked boolean default false;
alter table users ADD COLUMN blocked_at timestamp DEFAULT CURRENT_TIMESTAMP;

CREATE TABLE feedback (
   id uuid NOT NULL,
   owner text NOT NULL DEFAULT '-',
   dt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   email text NOT NULL,
   name text NOT NULL,
   feedback text NOT NULL
);
alter table feedback ADD COLUMN location text NOT NULL DEFAULT '-';

CREATE TABLE course (
    id UUID PRIMARY KEY,
    owner UUID NOT NULL,
    descr VARCHAR(128) NOT NULL,
    course_start DATE NOT NULL,
    course_end DATE NOT NULL,
    type VARCHAR(10) NOT NULL,
    target text[] NOT NULL,
    notes TEXT NOT NULL
);

CREATE INDEX idx_course ON course(owner);

alter table injection ADD COLUMN delete boolean default false;