INSERT INTO album (id, name, created_at, updated_at)
VALUES ('967d5bb5-3a7a-4d5e-8a6c-febc8c5b3f13', 'Hollywood''s Bleeding', '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp),
       ('c809bf15-bc2c-4621-bb96-70af96fd5d67', 'AI YoungBoy 2', '2019-10-02 11:16:12'::timestamp, '2019-10-02 11:16:12'::timestamp),
       ('2367710a-d4fb-49f5-8860-557b337386dd', 'KIRK', '2019-10-05 05:21:11'::timestamp, '2019-10-05 05:21:11'::timestamp),
       ('b0a24f12-428f-4ff5-84d5-bc1fdcff6f03', 'Lover', '2019-10-11 19:43:18'::timestamp, '2019-10-11 19:43:18'::timestamp),
       ('e0bb80ec-75a6-4348-bfc3-6ac1e89b195e', 'So Much Fun', '2019-10-12 12:16:02'::timestamp, '2019-10-12 12:16:02'::timestamp);
INSERT INTO users(
    id, login, passwd, email, sex, birthday, height, weight)
VALUES ('3a56fd3a-e7ee-11eb-ba80-0242ac130004', 'tets', 'test123', 'test@email.test', 'M', current_date - interval '1 year' * 20, 190, 101);

INSERT INTO antro(
    id, owner,
    dt, general_age,
    general_hip, general_height, general_leglen, general_weight, general_handlen, general_shoulders,
    notes, basic, result_fat, result_nofat, result_energy)
VALUES ('350e2cd6-e8c0-11eb-9a03-0242ac130003', '3a56fd3a-e7ee-11eb-ba80-0242ac130004',
        current_date - interval '1 days' * 1, 31,
        0, 190, 0, 100, 1, 2,
        'test1', true, 0.7973375453893254, 99.20266245461067, 2512.7775090195905);

INSERT INTO antro(
    id, owner,
    dt, general_age,
    general_hip, general_height, general_leglen, general_weight, general_handlen, general_shoulders,
    notes, basic, result_fat, result_nofat, result_energy)
VALUES ('f0f789ce-e8c0-11eb-9a03-0242ac130003', '3a56fd3a-e7ee-11eb-ba80-0242ac130004',
        current_date - interval '1 days' * 3, 31,
        0, 190, 0, 100, 1, 2,
        'test1', true, 0.7973375453893254, 99.20266245461067, 2512.7775090195905);


INSERT INTO public.injection(
    id, owner, what)
VALUES ('3a56fd3a-e7ee-22eb-ba80-0242ac130004', '3a56fd3a-e7ee-11eb-ba80-0242ac130004', 1);
INSERT INTO public.injection_dose(
    id, id_injection,
    dose, drug, volume, solvent)
VALUES ('3a56fd3a-e7ee-33eb-ba80-0242ac130004', '3a56fd3a-e7ee-22eb-ba80-0242ac130004',
        1, '00000002-0003-0000-0000-ff00ff00ff00', 1, 'O');
INSERT INTO public.injection_dose(
    id, id_injection,
    dose, drug, volume, solvent)
VALUES ('3a56fd3a-e7ee-44eb-ba80-0242ac130004', '3a56fd3a-e7ee-22eb-ba80-0242ac130004',
        1, '00000006-0003-0000-0000-ff00ff00ff00', 5, 'L');