CREATE TABLE IF NOT EXISTS migrations (
   MigrationID SERIAL PRIMARY KEY,
   VersionID INTEGER default 0
);

DO LANGUAGE plpgsql $$
BEGIN
   IF NOT EXISTS (SELECT FROM migrations WHERE MigrationID = 1) THEN
      INSERT INTO migrations (MigrationID, VersionID) VALUES (1, 0);
   END IF;
END $$;