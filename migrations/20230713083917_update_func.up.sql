CREATE FUNCTION update_saved_timetable(new_value text, group_id integer) RETURNS boolean AS $$
DECLARE
    old_value text;
BEGIN
    SELECT shedule INTO old_value
    FROM saved_timetable
    WHERE groupp = group_id;

    IF old_value <> new_value THEN
        UPDATE saved_timetable SET shedule = new_value, date_update=Now()
        WHERE groupp = group_id;

        RAISE NOTICE 'New value: %, Is updated: true', new_value;
        RETURN true;
    ELSE
        RAISE NOTICE 'No rows updated';
        RETURN false;
    END IF;
END $$ LANGUAGE plpgsql;