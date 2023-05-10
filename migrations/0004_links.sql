-- +goose Up
CREATE TYPE chat_record AS (
    id INT,
    title TEXT,
    description TEXT
);

CREATE TEMPORARY TABLE chat_results (
    all_ids    chat_record,
    city_ids   chat_record,
    dev_ids    chat_record,
    rec_ids    chat_record
);
-- +goose StatementBegin
-- select * from get_chat_ids(523019020);
CREATE FUNCTION get_chat_ids(user_id BIGINT)
    RETURNS TABLE (
        all_ids    RECORD,
        city_ids   RECORD,
        dev_ids    RECORD,
        rec_ids    RECORD
    ) AS $$
DECLARE
    _user RECORD;
BEGIN
    SELECT * INTO _user FROM bx_users
    JOIN tg_users ON bx_users.id = tg_users.bx_user
    WHERE tg_users.id = user_id;

    FOR all_ids IN SELECT id, title, description FROM tg_chats WHERE type = 'all' LOOP
        city_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        rec_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        dev_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);

        RETURN NEXT;
    END LOOP;

    FOR city_ids IN SELECT tg_chats.id, tg_chats.title FROM tg_chats JOIN bx_cities_chats ON tg_chats.id = bx_cities_chats.tg_chat WHERE bx_cities_chats.bx_city = _user.city AND tg_chats.type = 'city' order by sort asc LOOP
        all_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        rec_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        dev_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);

        RETURN NEXT;
    END LOOP;

    FOR rec_ids IN SELECT id, title, description FROM tg_chats WHERE type = 'recommended' order by sort asc LOOP
        all_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        city_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        dev_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);

        RETURN NEXT;
    END LOOP;

    FOR dev_ids IN SELECT tg_chats.id, tg_chats.title FROM tg_chats JOIN bx_departments_chats ON tg_chats.id = bx_departments_chats.tg_chat WHERE bx_departments_chats.bx_department = _user.department  AND tg_chats.type = 'department' ORDER BY tg_chats.sort ASC LOOP
        all_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        city_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);
        rec_ids := (SELECT ROW(null, null, null)::RECORD FROM tg_chats WHERE false);

        RETURN NEXT;
    END LOOP;

    RETURN;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
DROP FUNCTION get_chat_ids;