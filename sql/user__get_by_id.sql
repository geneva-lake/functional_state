CREATE OR REPLACE FUNCTION users.user__get_by_id(IN "@id" uuid)
    RETURNS user jsonb
    LANGUAGE 'plpgsql'
    VOLATILE SECURITY DEFINER
    PARALLEL UNSAFE
AS $BODY$

--Collect information about user from database
--returns json struct

begin

  return query select jsonb_build_object(
    'id', t.id,
    'name', t.name,
    'email', t.email)
   from users.user t
   where t.id = "@id";
end;
$BODY$