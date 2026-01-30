-- migrate:up
CREATE TABLE installations (
	id SERIAL PRIMARY KEY,
	installation_id VARCHAR(255) NOT NULL,
	config_data JSONB NOT NULL,
	updated_by VARCHAR(255) NOT NULL,
	installed_by VARCHAR(255) NOT NULL
);


-- migrate:down
DROP TABLE installation_configurations;
