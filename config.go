package main

var CacheSql = map[string]struct{}{
	"": {},
	"SELECT opc.oid, opc.opcname, nsp.nspname FROM pg_opclass opc, pg_namespace nsp WHERE opc.opcnamespace = nsp.oid":  {},
	"SELECT opr.oid, opr.oprname, nsp.nspname FROM pg_operator opr, pg_namespace nsp WHERE opr.oprnamespace = nsp.oid": {},
	"SHOW server_version_num": {},
	"SELECT c.oid, nsp.nspname, c.collname FROM pg_collation c left join pg_namespace nsp on nsp.oid = c.collnamespace": {},
	"SELECT pc.collname, pn.nspname from pg_collation pc join pg_namespace pn on pc.collnamespace = pn.oid;":            {},
	"SELECT rolname AS name FROM pg_roles ORDER BY rolname ASC":                                                         {},
	"SELECT setting FROM pg_settings WHERE name = 'block_size'":                                                         {},
	"SELECT spcname AS name FROM pg_tablespace ORDER BY spcname ASC":                                                    {},
}
