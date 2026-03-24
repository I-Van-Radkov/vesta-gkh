-- 1. Заполняем категории заново
INSERT INTO gkh_categories 
    (id, slug, bonus_per_hit, min_hits_for_bonus, big_bonus)
VALUES
    ('ac4604bd-9961-476c-b1a5-fea530bf3fa4', 'heating', 18, 1, 65),  -- отопление
    ('e94bcd69-9060-4caa-ba6e-4d310af6fcb1', 'water',   18, 1, 65),  -- вода
    ('935f09f9-de9a-4e54-b5f1-47b4e466a633', 'repair',  18, 1, 65),  -- ремонт
    ('6feefdea-d966-4509-add8-f46fa5fd300d', 'electric',18, 1, 65),  -- электричество
    ('aab38fce-6d5a-4b30-a07f-daa1094b5750', 'waste',   18, 1, 65),  -- вывоз мусора (ТКО)
    ('290f54ed-1a79-479b-9699-084b154dd1e1', 'house',   18, 1, 65)   -- управляющая компания / дом в целом
ON CONFLICT (id) DO UPDATE 
    SET bonus_per_hit = EXCLUDED.bonus_per_hit,
        min_hits_for_bonus = EXCLUDED.min_hits_for_bonus,
        big_bonus = EXCLUDED.big_bonus;

-- 2. Ключевые слова (остаётся без изменений)
INSERT INTO gkh_keywords (keyword, weight, category, is_active) VALUES
('отопление',               55, 'heating', true),
('без отопления',           85, 'heating', true),
('отключили тепло',         85, 'heating', true),
('батареи холодные',        80, 'heating', true),
('теплосеть',               55, 'heating', true),
('порыв теплотрассы',       75, 'heating', true),
('авария на теплосети',     80, 'heating', true),
('теплоснабжение',          50, 'heating', true),

('без воды',                85, 'water',   true),
('отключили воду',          85, 'water',   true),
('горячей воды нет',        80, 'water',   true),
('порыв водовода',          75, 'water',   true),
('водоснабжение',           50, 'water',   true),

('обесточил',               80, 'electric',true),
('без света',               85, 'electric',true),
('отключение электричества',80, 'electric',true),
('электроснабжение',        50, 'electric',true),

('капремонт',               50, 'repair',  true),
('крыша',                   55, 'repair',  true),
('кровля',                  55, 'repair',  true),
('ремонт кровли',           60, 'repair',  true),

('тко',                     45, 'waste',   true),
('мусор',                   50, 'waste',   true),
('вывоз мусора',            55, 'waste',   true),
('контейнеры',              45, 'waste',   true),
('свалка',                  50, 'waste',   true),

('управляющая компания',    45, 'house',   true),
('ук',                      40, 'house',   true),
('тсж',                     40, 'house',   true),
('тариф',                   40, 'house',   true),
('перерасчет',              60, 'house',   true),
('благоустройство',         35, 'house',   true)
ON CONFLICT DO NOTHING;   -- если вдруг уже есть

-- 3. Привязка ключевых слов к категориям (остаётся как было)
INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT 'ac4604bd-9961-476c-b1a5-fea530bf3fa4', id 
FROM gkh_keywords WHERE keyword IN ('отопление','без отопления','отключили тепло','батареи холодные','теплосеть','порыв теплотрассы','авария на теплосети','теплоснабжение')
ON CONFLICT DO NOTHING;

INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT 'e94bcd69-9060-4caa-ba6e-4d310af6fcb1', id 
FROM gkh_keywords WHERE keyword IN ('без воды','отключили воду','горячей воды нет','порыв водовода','водоснабжение')
ON CONFLICT DO NOTHING;

INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT '6feefdea-d966-4509-add8-f46fa5fd300d', id 
FROM gkh_keywords WHERE keyword IN ('обесточил','без света','отключение электричества','электроснабжение')
ON CONFLICT DO NOTHING;

INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT '935f09f9-de9a-4e54-b5f1-47b4e466a633', id 
FROM gkh_keywords WHERE keyword IN ('капремонт','крыша','кровля','ремонт кровли')
ON CONFLICT DO NOTHING;

INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT 'aab38fce-6d5a-4b30-a07f-daa1094b5750', id 
FROM gkh_keywords WHERE keyword IN ('тко','мусор','вывоз мусора','контейнеры','свалка')
ON CONFLICT DO NOTHING;

INSERT INTO gkh_category_keywords (category_id, keyword_id)
SELECT '290f54ed-1a79-479b-9699-084b154dd1e1', id 
FROM gkh_keywords WHERE keyword IN ('управляющая компания','ук','тсж','тариф','перерасчет','благоустройство')
ON CONFLICT DO NOTHING;

-- 4. Регулярки (позитивные и негативные) — без изменений
INSERT INTO gkh_regex_rules (pattern, bonus_score, is_active) VALUES
('(?i)(отключили|без|отсутствует|холодно).{0,50}(отоплен|тепла|батаре|трубы|теплосет)', 85, true),
('(?i)(порыв|прорыв|авария|дефект).{0,40}(теплосет|трубопровод|водовод|канализац)', 80, true),
('(?i)(без|отключили|авария).{0,40}(вод|горяч|холодн)', 80, true),
('(?i)(обесточ|без свет|отключени).{0,40}(электр|свет)', 75, true),
('(?i)(мусор|тко|контейнер).{0,50}(переполн|не вывоз|запах|свалка)', 65, true),
('(?i)(капремонт|крыша|кровля).{0,30}(обрушилась|ремонт|течёт)', 70, true)
ON CONFLICT DO NOTHING;

INSERT INTO gkh_negative_regex (pattern, penalty, description, is_active) VALUES
('(?i)(дтп|столкновение|пьяный водитель|насмерть|погиб)', -65, 'ДТП и аварии', true),
('(?i)(убийств|ограблен|наркот|марихуан|убил)', -75, 'Тяжёлый криминал', true),
('(?i)(спорт|теннис|сноуборд|чемпионат|футбол)', -45, 'Спорт', true),
('(?i)(погода|снег|мороз|дождь|жара)', -35, 'Погода', true),
('(?i)(пожар.*(лес|трава|дом без связи с жкх))', -50, 'Пожар не ЖКХ', true)
ON CONFLICT DO NOTHING;