-- migrate:up
ALTER TABLE ONLY public.installations
    ADD CONSTRAINT installations_installation_id_key UNIQUE (installation_id);

-- migrate:down
ALTER TABLE ONLY public.installations
    DROP CONSTRAINT installations_installation_id_key;
